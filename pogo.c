#include <php.h>
#include <ext/json/php_json.h>
#include <zend_exceptions.h>
#include "pogo.h"
#include "pogo_arginfo.h"
#include "pogo_consts.h" 
#include "_cgo_export.h"

// SHM Includes
#include <sys/mman.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdlib.h>

// Class entries
zend_class_entry *go_future_ce;
zend_class_entry *go_channel_ce;
zend_class_entry *go_waitgroup_ce;
zend_class_entry *go_pool_ce;
zend_class_entry *go_worker_exception_ce;
zend_class_entry *go_timeout_exception_ce;

// Handlers
static zend_object_handlers pogo_handlers;

typedef struct {
    char *base;
    size_t size;
} shm_region_t;

static HashTable shm_registry;

static void proxy_log(char *msg, int level) {
    int php_level = E_WARNING;
    if (level >= 3) php_level = E_ERROR;
    if (level == 0) php_level = E_NOTICE;
    fprintf(stderr, "[GoHost] %s\n", msg);
}

static void pogo_free_object(zend_object *object)
{
    pogo_object *intern = pogo_object_from_obj(object);
    if (intern->owns_handle && intern->go_handle != 0) {
        removeGoObject(intern->go_handle);
        intern->go_handle = 0;
    }
    zend_object_std_dtor(&intern->std);
}

zend_object *pogo_create_object(zend_class_entry *class_type)
{
    pogo_object *intern = zend_object_alloc(sizeof(pogo_object), class_type);
    zend_object_std_init(&intern->std, class_type);
    object_properties_init(&intern->std, class_type);
    intern->std.handlers = &pogo_handlers;
    intern->go_handle = 0;
    intern->owns_handle = false;
    return &intern->std;
}

static void shm_registry_dtor(zval *pDest) {
    shm_region_t *region = (shm_region_t*)Z_PTR_P(pDest);
    if (region) {
        if (region->base) munmap(region->base, region->size);
        pefree(region, 1);
    }
}

PHP_FUNCTION(Pogo_Internal__gopogo_init) {
    _gopogo_init((uintptr_t)proxy_log);
}

PHP_FUNCTION(Pogo_Internal__shm_check) {
    zend_long fd;
    ZEND_PARSE_PARAMETERS_START(1, 1)
        Z_PARAM_LONG(fd)
    ZEND_PARSE_PARAMETERS_END();
    RETURN_BOOL(zend_hash_index_find(&shm_registry, (zend_ulong)fd) != NULL);
}

PHP_FUNCTION(Pogo_Internal__shm_read) {
    zend_long fd, offset, length;
    ZEND_PARSE_PARAMETERS_START(3, 3)
        Z_PARAM_LONG(fd)
        Z_PARAM_LONG(offset)
        Z_PARAM_LONG(length)
    ZEND_PARSE_PARAMETERS_END();

    shm_region_t *region = (shm_region_t*)zend_hash_index_find_ptr(&shm_registry, (zend_ulong)fd);
    if (!region) {
        zend_throw_exception(go_worker_exception_ce, "SHM FD not mapped", 0);
        RETURN_THROWS();
    }

    if (offset < 0 || length < 0 || (size_t)(offset + 1 + length) > region->size) {
        zend_throw_exception(go_worker_exception_ce, "SHM read out of bounds", 0);
        RETURN_THROWS();
    }

    unsigned char guard = (unsigned char)region->base[offset];
    if (guard != 0x02) {
        zend_throw_exception(go_worker_exception_ce, "SHM Data Corruption: Guard byte not READY", 0);
        RETURN_THROWS();
    }

    RETURN_STRINGL(region->base + offset + 1, length);
}

PHP_FUNCTION(Pogo_Internal__shm_decode) {
    zend_long fd, offset, length;
    ZEND_PARSE_PARAMETERS_START(3, 3)
        Z_PARAM_LONG(fd)
        Z_PARAM_LONG(offset)
        Z_PARAM_LONG(length)
    ZEND_PARSE_PARAMETERS_END();

    shm_region_t *region = (shm_region_t*)zend_hash_index_find_ptr(&shm_registry, (zend_ulong)fd);
    if (!region) RETURN_NULL();

    if (offset < 0 || length < 0 || (size_t)(offset + 1 + length) > region->size) {
        zend_throw_exception(go_worker_exception_ce, "SHM read out of bounds", 0);
        RETURN_THROWS();
    }

    unsigned char guard = (unsigned char)region->base[offset];
    if (guard != 0x02) {
        zend_throw_exception(go_worker_exception_ce, "SHM Data Corruption: Guard byte not READY", 0);
        RETURN_THROWS();
    }

    php_json_decode(return_value, region->base + offset + 1, (int)length, 1, PHP_JSON_PARSER_DEFAULT_DEPTH);
}

PHP_FUNCTION(Pogo_Internal_start_worker_pool)
{
    char *entrypoint = "job_runner.php";
    size_t entry_len = sizeof("job_runner.php") - 1;
    zend_long min_workers = 4;
    zend_long max_workers = 8;
    zend_long max_jobs = 0;
    zval *options = NULL;

    ZEND_PARSE_PARAMETERS_START(0, 5)
        Z_PARAM_OPTIONAL
        Z_PARAM_STRING(entrypoint, entry_len)
        Z_PARAM_LONG(min_workers)
        Z_PARAM_LONG(max_workers)
        Z_PARAM_LONG(max_jobs)
        Z_PARAM_ARRAY(options)
    ZEND_PARSE_PARAMETERS_END();

    zend_long shm_size = 64 * 1024 * 1024;
    zend_long ipc_timeout_ms = 500;
    zend_long scale_latency_ms = 50;
    zend_long job_timeout_ms = 0;

    if (options) {
        HashTable *ht = Z_ARRVAL_P(options);
        zval *val;
        if ((val = zend_hash_str_find(ht, "shm_size", sizeof("shm_size") - 1)) != NULL && Z_TYPE_P(val) == IS_LONG) {
            shm_size = Z_LVAL_P(val);
        }
        if ((val = zend_hash_str_find(ht, "ipc_timeout_ms", sizeof("ipc_timeout_ms") - 1)) != NULL && Z_TYPE_P(val) == IS_LONG) {
            ipc_timeout_ms = Z_LVAL_P(val);
        }
        if ((val = zend_hash_str_find(ht, "scale_latency_ms", sizeof("scale_latency_ms") - 1)) != NULL && Z_TYPE_P(val) == IS_LONG) {
            scale_latency_ms = Z_LVAL_P(val);
        }
        if ((val = zend_hash_str_find(ht, "job_timeout_ms", sizeof("job_timeout_ms") - 1)) != NULL && Z_TYPE_P(val) == IS_LONG) {
            job_timeout_ms = Z_LVAL_P(val);
        }
    }

    start_workers_wrapper(entrypoint, (int)entry_len, min_workers, max_workers, max_jobs, shm_size, ipc_timeout_ms, scale_latency_ms, job_timeout_ms);
}

PHP_FUNCTION(Pogo_Internal_dispatch)
{
    char *name;
    size_t name_len;
    zval *payload;

    ZEND_PARSE_PARAMETERS_START(2, 2)
        Z_PARAM_STRING(name, name_len)
        Z_PARAM_ARRAY(payload)
    ZEND_PARSE_PARAMETERS_END();

    dispatch_wrapper(name, (int)name_len, payload);
}

PHP_FUNCTION(Pogo_Internal_dispatch_task)
{
    char *task_name;
    size_t task_name_len;
    zval *payload = NULL;

    ZEND_PARSE_PARAMETERS_START(1, 2)
        Z_PARAM_STRING(task_name, task_name_len)
        Z_PARAM_OPTIONAL
        Z_PARAM_ARRAY(payload)
    ZEND_PARSE_PARAMETERS_END();

    uintptr_t ch_handle = dispatch_task_wrapper(task_name, (int)task_name_len, payload);
    if (ch_handle == 0) {
        zend_throw_exception(go_worker_exception_ce, "Pool is shutting down or unavailable", 0);
        RETURN_THROWS();
    }

    object_init_ex(return_value, go_future_ce);

    pogo_object *intern_fut = pogo_object_from_obj(Z_OBJ_P(return_value));
    intern_fut->go_handle = ch_handle;
    intern_fut->owns_handle = true; // FIX: Future owns the handle
}

PHP_FUNCTION(Pogo_Internal_async)
{
    char *class_name;
    size_t class_name_len;
    zval *args = NULL;

    ZEND_PARSE_PARAMETERS_START(1, 2)
        Z_PARAM_STRING(class_name, class_name_len)
        Z_PARAM_OPTIONAL
        Z_PARAM_ARRAY(args)
    ZEND_PARSE_PARAMETERS_END();

    uintptr_t ch_handle = async_wrapper(class_name, (int)class_name_len, args);
    if (ch_handle == 0) {
        zend_throw_exception(go_worker_exception_ce, "Pool is shutting down or unavailable", 0);
        RETURN_THROWS();
    }

    object_init_ex(return_value, go_future_ce);

    pogo_object *intern_fut = pogo_object_from_obj(Z_OBJ_P(return_value));
    intern_fut->go_handle = ch_handle;
    intern_fut->owns_handle = true; // FIX: Future owns the handle
}

typedef struct {
    zend_string *str_key;
    zend_ulong num_key;
    int key_type;
} select_key_t;

PHP_FUNCTION(Pogo_Internal_select)
{
    zval *cases;
    double timeout = -1.0;

    ZEND_PARSE_PARAMETERS_START(1, 2)
        Z_PARAM_ARRAY(cases)
        Z_PARAM_OPTIONAL
        Z_PARAM_DOUBLE(timeout)
    ZEND_PARSE_PARAMETERS_END();

    HashTable *ht = Z_ARRVAL_P(cases);
    int count = zend_hash_num_elements(ht);

    if (count == 0) {
        RETURN_NULL();
    }

    uintptr_t *handles = safe_emalloc(count, sizeof(uintptr_t), 0);
    select_key_t *keys = safe_emalloc(count, sizeof(select_key_t), 0);

    int i = 0;
    zval *val;
    zend_string *str_key;
    zend_ulong num_key;

    ZEND_HASH_FOREACH_KEY_VAL(ht, num_key, str_key, val) {
        if (str_key) {
            keys[i].key_type = HASH_KEY_IS_STRING;
            keys[i].str_key = str_key;
            zend_string_addref(str_key);
        } else {
            keys[i].key_type = HASH_KEY_IS_LONG;
            keys[i].num_key = num_key;
        }

        uintptr_t handle = 0;
        if (Z_TYPE_P(val) == IS_OBJECT) {
            if (Z_OBJ_P(val)->ce == go_channel_ce || Z_OBJ_P(val)->ce == go_future_ce) {
                pogo_object *intern = pogo_object_from_obj(Z_OBJ_P(val));
                handle = intern->go_handle;
            }
        }
        handles[i] = handle;
        i++;
    } ZEND_HASH_FOREACH_END();

    select_result res = select_wrapper(handles, count, timeout);

    for(int j=0; j<count; j++) {
        if (keys[j].key_type == HASH_KEY_IS_STRING) {
            zend_string_release(keys[j].str_key);
        }
    }
    efree(handles);

    if (res.status == 1) {
        efree(keys);
        RETURN_NULL();
    }

    array_init(return_value);
    int idx = (int)res.index;

    if (idx >= 0 && idx < count) {
        if (keys[idx].key_type == HASH_KEY_IS_STRING) {
            add_assoc_string(return_value, "key", ZSTR_VAL(keys[idx].str_key));
        } else {
            add_assoc_long(return_value, "key", keys[idx].num_key);
        }
    }
    efree(keys);

    if (res.value) {
        add_assoc_string(return_value, "value", res.value);
        free(res.value);
    } else {
        add_assoc_string(return_value, "value", "");
    }
}

PHP_FUNCTION(Pogo_Internal_get_pool_stats)
{
    zend_long poolID = 0;
    ZEND_PARSE_PARAMETERS_START(0, 1)
        Z_PARAM_OPTIONAL
        Z_PARAM_LONG(poolID)
    ZEND_PARSE_PARAMETERS_END();

    char *json_res = get_pool_stats_wrapper((long)poolID);
    if (json_res == NULL) RETURN_EMPTY_STRING();

    RETVAL_STRING(json_res);
    free(json_res);
}

PHP_FUNCTION(Pogo_Internal_version)
{
    char *v = Pogo_version();
    RETVAL_STRING(v);
    free(v);
}

PHP_METHOD(Pogo_Internal_Pool, __construct) {
    char *entrypoint;
    size_t entry_len;
    zend_long min = 1;
    zend_long max = 8;
    zend_long max_jobs = 0;
    zval *options = NULL;

    ZEND_PARSE_PARAMETERS_START(1, 5)
        Z_PARAM_STRING(entrypoint, entry_len)
        Z_PARAM_OPTIONAL
        Z_PARAM_LONG(min)
        Z_PARAM_LONG(max)
        Z_PARAM_LONG(max_jobs)
        Z_PARAM_ARRAY(options)
    ZEND_PARSE_PARAMETERS_END();

    pogo_object *intern = pogo_object_from_obj(Z_OBJ_P(ZEND_THIS));
    intern->go_handle = (uintptr_t)create_pool_wrapper();
    intern->owns_handle = false;

    zend_update_property_stringl(go_pool_ce, Z_OBJ_P(ZEND_THIS), "entrypoint", sizeof("entrypoint")-1, entrypoint, entry_len);
    zend_update_property_long(go_pool_ce, Z_OBJ_P(ZEND_THIS), "min", sizeof("min")-1, min);
    zend_update_property_long(go_pool_ce, Z_OBJ_P(ZEND_THIS), "max", sizeof("max")-1, max);
    zend_update_property_long(go_pool_ce, Z_OBJ_P(ZEND_THIS), "max_jobs", sizeof("max_jobs")-1, max_jobs);

    if (options) {
        zend_update_property(go_pool_ce, Z_OBJ_P(ZEND_THIS), "options", sizeof("options")-1, options);
    } else {
        zval empty_arr;
        array_init(&empty_arr);
        zend_update_property(go_pool_ce, Z_OBJ_P(ZEND_THIS), "options", sizeof("options")-1, &empty_arr);
        zval_ptr_dtor(&empty_arr);
    }
}

PHP_METHOD(Pogo_Internal_Pool, start) {
    pogo_object *intern = pogo_object_from_obj(Z_OBJ_P(ZEND_THIS));
    long poolID = (long)intern->go_handle;

    zval *entrypoint = zend_read_property(go_pool_ce, Z_OBJ_P(ZEND_THIS), "entrypoint", sizeof("entrypoint")-1, 1, NULL);
    zval *min = zend_read_property(go_pool_ce, Z_OBJ_P(ZEND_THIS), "min", sizeof("min")-1, 1, NULL);
    zval *max = zend_read_property(go_pool_ce, Z_OBJ_P(ZEND_THIS), "max", sizeof("max")-1, 1, NULL);
    zval *max_jobs = zend_read_property(go_pool_ce, Z_OBJ_P(ZEND_THIS), "max_jobs", sizeof("max_jobs")-1, 1, NULL);
    zval *options = zend_read_property(go_pool_ce, Z_OBJ_P(ZEND_THIS), "options", sizeof("options")-1, 1, NULL);

    zend_long shm_size = 64 * 1024 * 1024;
    zend_long ipc_timeout_ms = 500;
    zend_long scale_latency_ms = 50;
    zend_long job_timeout_ms = 0;

    if (options && Z_TYPE_P(options) == IS_ARRAY) {
        HashTable *ht = Z_ARRVAL_P(options);
        zval *val;
        if ((val = zend_hash_str_find(ht, "shm_size", sizeof("shm_size") - 1)) != NULL && Z_TYPE_P(val) == IS_LONG) {
            shm_size = Z_LVAL_P(val);
        }
        if ((val = zend_hash_str_find(ht, "ipc_timeout_ms", sizeof("ipc_timeout_ms") - 1)) != NULL && Z_TYPE_P(val) == IS_LONG) {
            ipc_timeout_ms = Z_LVAL_P(val);
        }
        if ((val = zend_hash_str_find(ht, "scale_latency_ms", sizeof("scale_latency_ms") - 1)) != NULL && Z_TYPE_P(val) == IS_LONG) {
            scale_latency_ms = Z_LVAL_P(val);
        }
        if ((val = zend_hash_str_find(ht, "job_timeout_ms", sizeof("job_timeout_ms") - 1)) != NULL && Z_TYPE_P(val) == IS_LONG) {
            job_timeout_ms = Z_LVAL_P(val);
        }
    }

    start_pool_wrapper(poolID, Z_STRVAL_P(entrypoint), (int)Z_STRLEN_P(entrypoint), Z_LVAL_P(min), Z_LVAL_P(max), Z_LVAL_P(max_jobs), shm_size, ipc_timeout_ms, scale_latency_ms, job_timeout_ms);
}

PHP_METHOD(Pogo_Internal_Pool, shutdown) {
    pogo_object *intern = pogo_object_from_obj(Z_OBJ_P(ZEND_THIS));
    long poolID = (long)intern->go_handle;
    shutdown_pool_wrapper(poolID);
}

PHP_METHOD(Pogo_Internal_Pool, submit) {
    char *class_name;
    size_t class_name_len;
    zval *args = NULL;

    ZEND_PARSE_PARAMETERS_START(1, 2)
        Z_PARAM_STRING(class_name, class_name_len)
        Z_PARAM_OPTIONAL
        Z_PARAM_ARRAY(args)
    ZEND_PARSE_PARAMETERS_END();

    pogo_object *intern = pogo_object_from_obj(Z_OBJ_P(ZEND_THIS));
    long poolID = (long)intern->go_handle;

    uintptr_t ch_handle = async_on_pool_wrapper(poolID, class_name, (int)class_name_len, args);
    if (ch_handle == 0) {
        zend_throw_exception(go_worker_exception_ce, "Pool is shutting down or unavailable", 0);
        RETURN_THROWS();
    }

    object_init_ex(return_value, go_future_ce);
    pogo_object *intern_fut = pogo_object_from_obj(Z_OBJ_P(return_value));
    intern_fut->go_handle = ch_handle;
    intern_fut->owns_handle = true; // FIX: Future owns the handle
}

PHP_METHOD(Pogo_Internal_Pool, id) {
    pogo_object *intern = pogo_object_from_obj(Z_OBJ_P(ZEND_THIS));
    RETURN_LONG((long)intern->go_handle);
}

PHP_METHOD(Pogo_Internal_WaitGroup, __construct) {
    pogo_object *intern = pogo_object_from_obj(Z_OBJ_P(ZEND_THIS));
    intern->go_handle = create_WaitGroup_object();
    intern->owns_handle = true;
}
PHP_METHOD(Pogo_Internal_WaitGroup, add) {
    zend_long delta = 1;
    ZEND_PARSE_PARAMETERS_START(0, 1) Z_PARAM_OPTIONAL Z_PARAM_LONG(delta) ZEND_PARSE_PARAMETERS_END();
    add_wrapper(pogo_object_from_obj(Z_OBJ_P(ZEND_THIS))->go_handle, delta);
}
PHP_METHOD(Pogo_Internal_WaitGroup, done) {
    done_wrapper(pogo_object_from_obj(Z_OBJ_P(ZEND_THIS))->go_handle);
}
PHP_METHOD(Pogo_Internal_WaitGroup, wait) {
    wait_wrapper(pogo_object_from_obj(Z_OBJ_P(ZEND_THIS))->go_handle);
}

PHP_METHOD(Pogo_Internal_Channel, __construct) {
    pogo_object *intern = pogo_object_from_obj(Z_OBJ_P(ZEND_THIS));
    intern->go_handle = create_Channel_object();
    intern->owns_handle = true;
}
PHP_METHOD(Pogo_Internal_Channel, init) {
    zend_long capacity = 0;
    ZEND_PARSE_PARAMETERS_START(0, 1) Z_PARAM_OPTIONAL Z_PARAM_LONG(capacity) ZEND_PARSE_PARAMETERS_END();
    init_wrapper(pogo_object_from_obj(Z_OBJ_P(ZEND_THIS))->go_handle, capacity);
}
PHP_METHOD(Pogo_Internal_Channel, push) {
    char *val; size_t val_len;
    ZEND_PARSE_PARAMETERS_START(1, 1) Z_PARAM_STRING(val, val_len) ZEND_PARSE_PARAMETERS_END();
    push_wrapper(pogo_object_from_obj(Z_OBJ_P(ZEND_THIS))->go_handle, val, (int)val_len);
}
PHP_METHOD(Pogo_Internal_Channel, pop) {
    char *res = pop_wrapper(pogo_object_from_obj(Z_OBJ_P(ZEND_THIS))->go_handle);
    if (res) { RETVAL_STRING(res); free(res); } else { RETURN_EMPTY_STRING(); }
}
PHP_METHOD(Pogo_Internal_Channel, close) {
    close_wrapper(pogo_object_from_obj(Z_OBJ_P(ZEND_THIS))->go_handle);
}

PHP_METHOD(Pogo_Internal_Future, __construct) {}
PHP_METHOD(Pogo_Internal_Future, await) {
    double timeout = -1.0;
    ZEND_PARSE_PARAMETERS_START(0, 1) Z_PARAM_OPTIONAL Z_PARAM_DOUBLE(timeout) ZEND_PARSE_PARAMETERS_END();

    pogo_object *intern = pogo_object_from_obj(Z_OBJ_P(ZEND_THIS));
    if (intern->go_handle == 0) {
        RETURN_NULL();
    }

    char *res = await_wrapper(intern->go_handle, timeout);
    if (res == NULL) {
        RETURN_NULL();
    }

    RETVAL_STRING(res);
    free(res);
}
PHP_METHOD(Pogo_Internal_Future, done) {
    pogo_object *intern = pogo_object_from_obj(Z_OBJ_P(ZEND_THIS));
    if (intern->go_handle == 0) RETURN_FALSE;

    char *res = poll_wrapper(intern->go_handle);
    if (res == NULL) RETURN_FALSE;
    
    RETURN_TRUE;
}
PHP_METHOD(Pogo_Internal_Future, cancel) {
    pogo_object *intern = pogo_object_from_obj(Z_OBJ_P(ZEND_THIS));
    if (intern->go_handle == 0) RETURN_FALSE;
    RETURN_BOOL(cancel_wrapper(intern->go_handle));
}

PHP_MSHUTDOWN_FUNCTION(pogo)
{
    zend_hash_destroy(&shm_registry);
    Pogo_shutdown_module();
    return SUCCESS;
}

PHP_MINIT_FUNCTION(pogo)
{
    memcpy(&pogo_handlers, &std_object_handlers, sizeof(zend_object_handlers));
    pogo_handlers.free_obj = pogo_free_object;
    pogo_handlers.offset = offsetof(pogo_object, std);

    go_future_ce = register_class_Pogo_Internal_Future();
    go_future_ce->create_object = pogo_create_object;

    go_channel_ce = register_class_Pogo_Internal_Channel();
    go_channel_ce->create_object = pogo_create_object;

    go_waitgroup_ce = register_class_Pogo_Internal_WaitGroup();
    go_waitgroup_ce->create_object = pogo_create_object;

    zend_class_entry ce_pool;
    INIT_CLASS_ENTRY(ce_pool, "Pogo\\Internal\\Pool", class_Pogo_Internal_Pool_methods);
    go_pool_ce = zend_register_internal_class(&ce_pool);
    go_pool_ce->create_object = pogo_create_object;
    zend_declare_property_string(go_pool_ce, "entrypoint", sizeof("entrypoint")-1, "", ZEND_ACC_PRIVATE);
    zend_declare_property_long(go_pool_ce, "min", sizeof("min")-1, 0, ZEND_ACC_PRIVATE);
    zend_declare_property_long(go_pool_ce, "max", sizeof("max")-1, 0, ZEND_ACC_PRIVATE);
    zend_declare_property_long(go_pool_ce, "max_jobs", sizeof("max_jobs")-1, 0, ZEND_ACC_PRIVATE);
    zend_declare_property_null(go_pool_ce, "options", sizeof("options")-1, ZEND_ACC_PRIVATE);

    zend_class_entry ce;
    INIT_CLASS_ENTRY(ce, "Pogo\\WorkerException", NULL);
    go_worker_exception_ce = zend_register_internal_class_ex(&ce, zend_ce_exception);
    INIT_CLASS_ENTRY(ce, "Pogo\\TimeoutException", NULL);
    go_timeout_exception_ce = zend_register_internal_class_ex(&ce, zend_ce_exception);

    _gopogo_init((uintptr_t)proxy_log);

    zend_hash_init(&shm_registry, 0, NULL, shm_registry_dtor, 1);

    char *env_fd = getenv("FRANKENPHP_WORKER_SHM_FD");
    if (env_fd) {
        long shm_fd = strtol(env_fd, NULL, 10);

        struct stat sb;
        if (fstat((int)shm_fd, &sb) != -1 && sb.st_size > 0) {
            size_t size = sb.st_size;
            char *base = mmap(NULL, size, PROT_READ|PROT_WRITE, MAP_SHARED, (int)shm_fd, 0);
            if (base == MAP_FAILED) {
                fprintf(stderr, "[GoWorker] SHM mmap failed\n");
            } else {
                shm_region_t *region = pemalloc(sizeof(shm_region_t), 1);
                region->base = base;
                region->size = size;
                zend_hash_index_update_ptr(&shm_registry, (zend_ulong)shm_fd, region);
            }
        }
    }

    return SUCCESS;
}

zend_module_entry pogo_module_entry = {
    STANDARD_MODULE_HEADER,
    "pogo",
    ext_functions,
    PHP_MINIT(pogo),
    PHP_MSHUTDOWN(pogo),
    NULL,
    NULL,
    NULL,
    "0.1",
    STANDARD_MODULE_PROPERTIES
};