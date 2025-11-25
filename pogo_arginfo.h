/* This is a generated file, edit the .stub.php file instead.
 * Stub hash: baa32204a9b4511c5326a3ef38ec46be45739039 */

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Go__gopogo_init, 0, 0, IS_VOID, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Go__shm_read, 0, 3, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO(0, fd, IS_LONG, 0)
	ZEND_ARG_TYPE_INFO(0, offset, IS_LONG, 0)
	ZEND_ARG_TYPE_INFO(0, length, IS_LONG, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Go__shm_decode, 0, 3, IS_MIXED, 0)
	ZEND_ARG_TYPE_INFO(0, fd, IS_LONG, 0)
	ZEND_ARG_TYPE_INFO(0, offset, IS_LONG, 0)
	ZEND_ARG_TYPE_INFO(0, length, IS_LONG, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Go__shm_check, 0, 1, _IS_BOOL, 0)
	ZEND_ARG_TYPE_INFO(0, fd, IS_LONG, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Go_start_worker_pool, 0, 0, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, entrypoint, IS_STRING, 0, "\"job_runner.php\"")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, minWorkers, IS_LONG, 0, "4")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, maxWorkers, IS_LONG, 0, "8")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, maxJobs, IS_LONG, 0, "0")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, options, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Go_dispatch, 0, 2, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO(0, workerName, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO(0, payload, IS_ARRAY, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_OBJ_INFO_EX(arginfo_Go_dispatch_task, 0, 1, Go\\Future, 0)
	ZEND_ARG_TYPE_INFO(0, taskName, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, payload, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Go_select, 0, 1, IS_ARRAY, 1)
	ZEND_ARG_TYPE_INFO(0, cases, IS_ARRAY, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, timeout, IS_DOUBLE, 1, "null")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_OBJ_INFO_EX(arginfo_Go_async, 0, 1, Go\\Future, 0)
	ZEND_ARG_TYPE_INFO(0, class, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, args, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Go_get_pool_stats, 0, 0, IS_ARRAY, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, poolID, IS_LONG, 0, "0")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_INFO_EX(arginfo_class_Go_Future___construct, 0, 0, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Go_Future_await, 0, 0, IS_MIXED, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, timeout, IS_DOUBLE, 1, "null")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Go_Future_done, 0, 0, _IS_BOOL, 0)
ZEND_END_ARG_INFO()

#define arginfo_class_Go_Future_cancel arginfo_class_Go_Future_done

#define arginfo_class_Go_WaitGroup___construct arginfo_class_Go_Future___construct

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Go_WaitGroup_add, 0, 0, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, delta, IS_LONG, 0, "1")
ZEND_END_ARG_INFO()

#define arginfo_class_Go_WaitGroup_done arginfo_Go__gopogo_init

#define arginfo_class_Go_WaitGroup_wait arginfo_Go__gopogo_init

#define arginfo_class_Go_Channel___construct arginfo_class_Go_Future___construct

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Go_Channel_init, 0, 0, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, capacity, IS_LONG, 0, "0")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Go_Channel_push, 0, 1, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO(0, value, IS_STRING, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Go_Channel_pop, 0, 0, IS_STRING, 0)
ZEND_END_ARG_INFO()

#define arginfo_class_Go_Channel_close arginfo_Go__gopogo_init

#define arginfo_class_Go_Contract_Resettable_reset arginfo_Go__gopogo_init

ZEND_BEGIN_ARG_INFO_EX(arginfo_class_Go_Runtime_Pool___construct, 0, 0, 1)
	ZEND_ARG_TYPE_INFO(0, entrypoint, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, minWorkers, IS_LONG, 0, "1")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, maxWorkers, IS_LONG, 0, "8")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, maxJobs, IS_LONG, 0, "0")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, options, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

#define arginfo_class_Go_Runtime_Pool_start arginfo_Go__gopogo_init

#define arginfo_class_Go_Runtime_Pool_shutdown arginfo_Go__gopogo_init

ZEND_BEGIN_ARG_WITH_RETURN_OBJ_INFO_EX(arginfo_class_Go_Runtime_Pool_submit, 0, 1, Go\\Future, 0)
	ZEND_ARG_TYPE_INFO(0, jobClass, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, args, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

ZEND_FUNCTION(Go__gopogo_init);
ZEND_FUNCTION(Go__shm_read);
ZEND_FUNCTION(Go__shm_decode);
ZEND_FUNCTION(Go__shm_check);
ZEND_FUNCTION(Go_start_worker_pool);
ZEND_FUNCTION(Go_dispatch);
ZEND_FUNCTION(Go_dispatch_task);
ZEND_FUNCTION(Go_select);
ZEND_FUNCTION(Go_async);
ZEND_FUNCTION(Go_get_pool_stats);
ZEND_METHOD(Go_Future, __construct);
ZEND_METHOD(Go_Future, await);
ZEND_METHOD(Go_Future, done);
ZEND_METHOD(Go_Future, cancel);
ZEND_METHOD(Go_WaitGroup, __construct);
ZEND_METHOD(Go_WaitGroup, add);
ZEND_METHOD(Go_WaitGroup, done);
ZEND_METHOD(Go_WaitGroup, wait);
ZEND_METHOD(Go_Channel, __construct);
ZEND_METHOD(Go_Channel, init);
ZEND_METHOD(Go_Channel, push);
ZEND_METHOD(Go_Channel, pop);
ZEND_METHOD(Go_Channel, close);
ZEND_METHOD(Go_Runtime_Pool, __construct);
ZEND_METHOD(Go_Runtime_Pool, start);
ZEND_METHOD(Go_Runtime_Pool, shutdown);
ZEND_METHOD(Go_Runtime_Pool, submit);

static const zend_function_entry ext_functions[] = {
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Go", "_gopogo_init"), zif_Go__gopogo_init, arginfo_Go__gopogo_init, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Go", "_shm_read"), zif_Go__shm_read, arginfo_Go__shm_read, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Go", "_shm_decode"), zif_Go__shm_decode, arginfo_Go__shm_decode, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Go", "_shm_check"), zif_Go__shm_check, arginfo_Go__shm_check, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Go", "start_worker_pool"), zif_Go_start_worker_pool, arginfo_Go_start_worker_pool, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Go", "dispatch"), zif_Go_dispatch, arginfo_Go_dispatch, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Go", "dispatch_task"), zif_Go_dispatch_task, arginfo_Go_dispatch_task, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Go", "select"), zif_Go_select, arginfo_Go_select, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Go", "async"), zif_Go_async, arginfo_Go_async, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Go", "get_pool_stats"), zif_Go_get_pool_stats, arginfo_Go_get_pool_stats, 0, NULL, NULL)
	ZEND_FE_END
};

static const zend_function_entry class_Go_Future_methods[] = {
	ZEND_ME(Go_Future, __construct, arginfo_class_Go_Future___construct, ZEND_ACC_PUBLIC)
	ZEND_ME(Go_Future, await, arginfo_class_Go_Future_await, ZEND_ACC_PUBLIC)
	ZEND_ME(Go_Future, done, arginfo_class_Go_Future_done, ZEND_ACC_PUBLIC)
	ZEND_ME(Go_Future, cancel, arginfo_class_Go_Future_cancel, ZEND_ACC_PUBLIC)
	ZEND_FE_END
};

static const zend_function_entry class_Go_WaitGroup_methods[] = {
	ZEND_ME(Go_WaitGroup, __construct, arginfo_class_Go_WaitGroup___construct, ZEND_ACC_PUBLIC)
	ZEND_ME(Go_WaitGroup, add, arginfo_class_Go_WaitGroup_add, ZEND_ACC_PUBLIC)
	ZEND_ME(Go_WaitGroup, done, arginfo_class_Go_WaitGroup_done, ZEND_ACC_PUBLIC)
	ZEND_ME(Go_WaitGroup, wait, arginfo_class_Go_WaitGroup_wait, ZEND_ACC_PUBLIC)
	ZEND_FE_END
};

static const zend_function_entry class_Go_Channel_methods[] = {
	ZEND_ME(Go_Channel, __construct, arginfo_class_Go_Channel___construct, ZEND_ACC_PUBLIC)
	ZEND_ME(Go_Channel, init, arginfo_class_Go_Channel_init, ZEND_ACC_PUBLIC)
	ZEND_ME(Go_Channel, push, arginfo_class_Go_Channel_push, ZEND_ACC_PUBLIC)
	ZEND_ME(Go_Channel, pop, arginfo_class_Go_Channel_pop, ZEND_ACC_PUBLIC)
	ZEND_ME(Go_Channel, close, arginfo_class_Go_Channel_close, ZEND_ACC_PUBLIC)
	ZEND_FE_END
};

static const zend_function_entry class_Go_Contract_Resettable_methods[] = {
	ZEND_RAW_FENTRY("reset", NULL, arginfo_class_Go_Contract_Resettable_reset, ZEND_ACC_PUBLIC|ZEND_ACC_ABSTRACT, NULL, NULL)
	ZEND_FE_END
};

static const zend_function_entry class_Go_Runtime_Pool_methods[] = {
	ZEND_ME(Go_Runtime_Pool, __construct, arginfo_class_Go_Runtime_Pool___construct, ZEND_ACC_PUBLIC)
	ZEND_ME(Go_Runtime_Pool, start, arginfo_class_Go_Runtime_Pool_start, ZEND_ACC_PUBLIC)
	ZEND_ME(Go_Runtime_Pool, shutdown, arginfo_class_Go_Runtime_Pool_shutdown, ZEND_ACC_PUBLIC)
	ZEND_ME(Go_Runtime_Pool, submit, arginfo_class_Go_Runtime_Pool_submit, ZEND_ACC_PUBLIC)
	ZEND_FE_END
};

static zend_class_entry *register_class_Go_Future(void)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Go", "Future", class_Go_Future_methods);
	class_entry = zend_register_internal_class_with_flags(&ce, NULL, 0);

	zval property_result_default_value;
	ZVAL_NULL(&property_result_default_value);
	zend_string *property_result_name = zend_string_init("result", sizeof("result") - 1, 1);
	zend_declare_typed_property(class_entry, property_result_name, &property_result_default_value, ZEND_ACC_PRIVATE, NULL, (zend_type) ZEND_TYPE_INIT_MASK(MAY_BE_ANY));
	zend_string_release(property_result_name);

	zval property_resolved_default_value;
	ZVAL_FALSE(&property_resolved_default_value);
	zend_string *property_resolved_name = zend_string_init("resolved", sizeof("resolved") - 1, 1);
	zend_declare_typed_property(class_entry, property_resolved_name, &property_resolved_default_value, ZEND_ACC_PRIVATE, NULL, (zend_type) ZEND_TYPE_INIT_MASK(MAY_BE_BOOL));
	zend_string_release(property_resolved_name);

	zval property_error_default_value;
	ZVAL_NULL(&property_error_default_value);
	zend_string *property_error_name = zend_string_init("error", sizeof("error") - 1, 1);
	zend_declare_typed_property(class_entry, property_error_name, &property_error_default_value, ZEND_ACC_PRIVATE, NULL, (zend_type) ZEND_TYPE_INIT_MASK(MAY_BE_STRING|MAY_BE_NULL));
	zend_string_release(property_error_name);

	zval property_channel_default_value;
	ZVAL_NULL(&property_channel_default_value);
	zend_string *property_channel_name = zend_string_init("channel", sizeof("channel") - 1, 1);
	zend_string *property_channel_class_Go_Channel = zend_string_init("Go\\Channel", sizeof("Go\\Channel")-1, 1);
	zend_declare_typed_property(class_entry, property_channel_name, &property_channel_default_value, ZEND_ACC_PRIVATE, NULL, (zend_type) ZEND_TYPE_INIT_CLASS(property_channel_class_Go_Channel, 0, MAY_BE_NULL));
	zend_string_release(property_channel_name);

	return class_entry;
}

static zend_class_entry *register_class_Go_WaitGroup(void)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Go", "WaitGroup", class_Go_WaitGroup_methods);
	class_entry = zend_register_internal_class_with_flags(&ce, NULL, 0);

	return class_entry;
}

static zend_class_entry *register_class_Go_Channel(void)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Go", "Channel", class_Go_Channel_methods);
	class_entry = zend_register_internal_class_with_flags(&ce, NULL, 0);

	return class_entry;
}

static zend_class_entry *register_class_Go_WorkerException(zend_class_entry *class_entry_Exception)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Go", "WorkerException", NULL);
	class_entry = zend_register_internal_class_with_flags(&ce, class_entry_Exception, 0);

	return class_entry;
}

static zend_class_entry *register_class_Go_TimeoutException(zend_class_entry *class_entry_Exception)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Go", "TimeoutException", NULL);
	class_entry = zend_register_internal_class_with_flags(&ce, class_entry_Exception, 0);

	return class_entry;
}

static zend_class_entry *register_class_Go_Contract_Resettable(void)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Go\\Contract", "Resettable", class_Go_Contract_Resettable_methods);
	class_entry = zend_register_internal_interface(&ce);

	return class_entry;
}

static zend_class_entry *register_class_Go_Runtime_Pool(void)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Go\\Runtime", "Pool", class_Go_Runtime_Pool_methods);
	class_entry = zend_register_internal_class_with_flags(&ce, NULL, 0);

	return class_entry;
}
