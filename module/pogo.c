#include <php.h>
#include <Zend/zend_exceptions.h>
#include <Zend/zend_smart_str.h>
#include <ext/json/php_json.h>
#include <ext/spl/spl_exceptions.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>

#include "_cgo_export.h"
#include "pogo.h"
#include "pogo_arginfo.h"

ZEND_BEGIN_MODULE_GLOBALS(pogo)
	HashTable *active_tasks;
ZEND_END_MODULE_GLOBALS(pogo)

ZEND_DECLARE_MODULE_GLOBALS(pogo)

#define POGO_G(v) ZEND_MODULE_GLOBALS_ACCESSOR(pogo, v)

static PHP_GINIT_FUNCTION(pogo)
{
	pogo_globals->active_tasks = NULL;
}

static zend_class_entry *runtime_exception_ce(void)
{
	return spl_ce_RuntimeException;
}

static void throw_runtime_from_cstr(char *message)
{
	if (message == NULL) {
		zend_throw_exception(runtime_exception_ce(), "Pogo error", 0);
		return;
	}

	zend_throw_exception(runtime_exception_ce(), message, 0);
	free(message);
}

PHP_MINIT_FUNCTION(pogo)
{
	return SUCCESS;
}

PHP_RINIT_FUNCTION(pogo)
{
	ALLOC_HASHTABLE(POGO_G(active_tasks));
	zend_hash_init(POGO_G(active_tasks), 8, NULL, NULL, 0);
	return SUCCESS;
}

PHP_RSHUTDOWN_FUNCTION(pogo)
{
	zend_ulong task;

	if (POGO_G(active_tasks) == NULL) {
		return SUCCESS;
	}

	ZEND_HASH_FOREACH_NUM_KEY(POGO_G(active_tasks), task) {
		go_pogo_cancel((uint64_t) task);
	} ZEND_HASH_FOREACH_END();

	zend_hash_destroy(POGO_G(active_tasks));
	FREE_HASHTABLE(POGO_G(active_tasks));
	POGO_G(active_tasks) = NULL;

	return SUCCESS;
}

static void track_task(uint64_t task)
{
	zval dummy;

	if (POGO_G(active_tasks) == NULL) {
		return;
	}

	ZVAL_NULL(&dummy);
	zend_hash_index_update(POGO_G(active_tasks), (zend_ulong) task, &dummy);
}

static void untrack_task(uint64_t task)
{
	if (POGO_G(active_tasks) == NULL) {
		return;
	}

	zend_hash_index_del(POGO_G(active_tasks), (zend_ulong) task);
}

zend_module_entry pogo_module_entry = {
	STANDARD_MODULE_HEADER,
	"pogo",
	ext_functions,
	PHP_MINIT(pogo),
	NULL,
	PHP_RINIT(pogo),
	PHP_RSHUTDOWN(pogo),
	NULL,
	"0.1.0",
	ZEND_MODULE_GLOBALS(pogo),
	PHP_GINIT(pogo),
	NULL,
	NULL,
	STANDARD_MODULE_PROPERTIES_EX
};

PHP_FUNCTION(pogo_spawn)
{
	char *class_name;
	size_t class_name_len;
	char *pool_name = "default";
	size_t pool_name_len = sizeof("default") - 1;
	zval *args = NULL;
	zval empty_args;
	smart_str json = {0};
	char *err = NULL;

	ZEND_PARSE_PARAMETERS_START(1, 3)
		Z_PARAM_STRING(class_name, class_name_len)
		Z_PARAM_OPTIONAL
		Z_PARAM_ARRAY(args)
		Z_PARAM_STRING(pool_name, pool_name_len)
	ZEND_PARSE_PARAMETERS_END();

	if (args == NULL) {
		array_init(&empty_args);
		args = &empty_args;
	}

	if (php_json_encode(&json, args, PHP_JSON_UNESCAPED_SLASHES) == FAILURE) {
		if (args == &empty_args) {
			zval_ptr_dtor(&empty_args);
		}
		smart_str_free(&json);
		zend_throw_exception(runtime_exception_ce(), "Failed to encode Pogo job args as JSON", 0);
		RETURN_THROWS();
	}

	if (args == &empty_args) {
		zval_ptr_dtor(&empty_args);
	}

	if (EG(exception)) {
		smart_str_free(&json);
		RETURN_THROWS();
	}

	if (json.s == NULL) {
		zend_throw_exception(runtime_exception_ce(), "Failed to encode Pogo job args as JSON", 0);
		RETURN_THROWS();
	}

	smart_str_0(&json);

	uint64_t task = go_pogo_spawn(
		pool_name,
		pool_name_len,
		class_name,
		class_name_len,
		ZSTR_VAL(json.s),
		ZSTR_LEN(json.s),
		&err
	);

	smart_str_free(&json);

	if (err != NULL) {
		throw_runtime_from_cstr(err);
		RETURN_THROWS();
	}

	track_task(task);

	RETURN_LONG((zend_long) task);
}

PHP_FUNCTION(pogo_await)
{
	zend_long task;
	double timeout = 5.0;
	char *err = NULL;
	char *json = NULL;

	ZEND_PARSE_PARAMETERS_START(1, 2)
		Z_PARAM_LONG(task)
		Z_PARAM_OPTIONAL
		Z_PARAM_DOUBLE(timeout)
	ZEND_PARSE_PARAMETERS_END();

	untrack_task((uint64_t) task);

	json = go_pogo_await((uint64_t) task, timeout, &err);

	if (err != NULL) {
		if (json != NULL) {
			free(json);
		}
		throw_runtime_from_cstr(err);
		RETURN_THROWS();
	}

	if (json == NULL) {
		RETURN_NULL();
	}

	php_json_decode(return_value, json, (int) strlen(json), 1, PHP_JSON_PARSER_DEFAULT_DEPTH);
	free(json);

	if (EG(exception)) {
		RETURN_THROWS();
	}
}

PHP_FUNCTION(pogo_pool_size)
{
	char *pool_name = "default";
	size_t pool_name_len = sizeof("default") - 1;

	ZEND_PARSE_PARAMETERS_START(0, 1)
		Z_PARAM_OPTIONAL
		Z_PARAM_STRING(pool_name, pool_name_len)
	ZEND_PARSE_PARAMETERS_END();

	RETURN_LONG((zend_long) go_pogo_pool_size(pool_name, pool_name_len));
}
