/* This is a generated file, edit the .stub.php file instead.
 * Stub hash: 197ca9fed088def349dcd5b4c6fcfe855e39c19e */

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo__gopogo_init, 0, 0, IS_VOID, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo__shm_read, 0, 3, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO(0, fd, IS_LONG, 0)
	ZEND_ARG_TYPE_INFO(0, offset, IS_LONG, 0)
	ZEND_ARG_TYPE_INFO(0, length, IS_LONG, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo__shm_decode, 0, 3, IS_MIXED, 0)
	ZEND_ARG_TYPE_INFO(0, fd, IS_LONG, 0)
	ZEND_ARG_TYPE_INFO(0, offset, IS_LONG, 0)
	ZEND_ARG_TYPE_INFO(0, length, IS_LONG, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo__shm_check, 0, 1, _IS_BOOL, 0)
	ZEND_ARG_TYPE_INFO(0, fd, IS_LONG, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo_start_worker_pool, 0, 0, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, entrypoint, IS_STRING, 0, "\"job_runner.php\"")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, minWorkers, IS_LONG, 0, "4")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, maxWorkers, IS_LONG, 0, "8")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, maxJobs, IS_LONG, 0, "0")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, options, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo_dispatch, 0, 2, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO(0, workerName, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO(0, payload, IS_ARRAY, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_OBJ_INFO_EX(arginfo_Pogo_dispatch_task, 0, 1, Pogo\\Future, 0)
	ZEND_ARG_TYPE_INFO(0, taskName, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, payload, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo_select, 0, 1, IS_ARRAY, 1)
	ZEND_ARG_TYPE_INFO(0, cases, IS_ARRAY, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, timeout, IS_DOUBLE, 1, "null")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_OBJ_INFO_EX(arginfo_Pogo_async, 0, 1, Pogo\\Future, 0)
	ZEND_ARG_TYPE_INFO(0, class, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, args, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo_get_pool_stats, 0, 0, IS_ARRAY, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, poolID, IS_LONG, 0, "0")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_INFO_EX(arginfo_class_Pogo_Future___construct, 0, 0, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Pogo_Future_await, 0, 0, IS_MIXED, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, timeout, IS_DOUBLE, 1, "null")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Pogo_Future_done, 0, 0, _IS_BOOL, 0)
ZEND_END_ARG_INFO()

#define arginfo_class_Pogo_Future_cancel arginfo_class_Pogo_Future_done

#define arginfo_class_Pogo_WaitGroup___construct arginfo_class_Pogo_Future___construct

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Pogo_WaitGroup_add, 0, 0, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, delta, IS_LONG, 0, "1")
ZEND_END_ARG_INFO()

#define arginfo_class_Pogo_WaitGroup_done arginfo_Pogo__gopogo_init

#define arginfo_class_Pogo_WaitGroup_wait arginfo_Pogo__gopogo_init

#define arginfo_class_Pogo_Channel___construct arginfo_class_Pogo_Future___construct

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Pogo_Channel_init, 0, 0, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, capacity, IS_LONG, 0, "0")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Pogo_Channel_push, 0, 1, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO(0, value, IS_STRING, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Pogo_Channel_pop, 0, 0, IS_STRING, 0)
ZEND_END_ARG_INFO()

#define arginfo_class_Pogo_Channel_close arginfo_Pogo__gopogo_init

#define arginfo_class_Pogo_Contract_Resettable_reset arginfo_Pogo__gopogo_init

ZEND_BEGIN_ARG_INFO_EX(arginfo_class_Pogo_Runtime_Pool___construct, 0, 0, 1)
	ZEND_ARG_TYPE_INFO(0, entrypoint, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, minWorkers, IS_LONG, 0, "1")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, maxWorkers, IS_LONG, 0, "8")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, maxJobs, IS_LONG, 0, "0")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, options, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

#define arginfo_class_Pogo_Runtime_Pool_start arginfo_Pogo__gopogo_init

#define arginfo_class_Pogo_Runtime_Pool_shutdown arginfo_Pogo__gopogo_init

ZEND_BEGIN_ARG_WITH_RETURN_OBJ_INFO_EX(arginfo_class_Pogo_Runtime_Pool_submit, 0, 1, Pogo\\Future, 0)
	ZEND_ARG_TYPE_INFO(0, jobClass, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, args, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

ZEND_FUNCTION(Pogo__gopogo_init);
ZEND_FUNCTION(Pogo__shm_read);
ZEND_FUNCTION(Pogo__shm_decode);
ZEND_FUNCTION(Pogo__shm_check);
ZEND_FUNCTION(Pogo_start_worker_pool);
ZEND_FUNCTION(Pogo_dispatch);
ZEND_FUNCTION(Pogo_dispatch_task);
ZEND_FUNCTION(Pogo_select);
ZEND_FUNCTION(Pogo_async);
ZEND_FUNCTION(Pogo_get_pool_stats);
ZEND_METHOD(Pogo_Future, __construct);
ZEND_METHOD(Pogo_Future, await);
ZEND_METHOD(Pogo_Future, done);
ZEND_METHOD(Pogo_Future, cancel);
ZEND_METHOD(Pogo_WaitGroup, __construct);
ZEND_METHOD(Pogo_WaitGroup, add);
ZEND_METHOD(Pogo_WaitGroup, done);
ZEND_METHOD(Pogo_WaitGroup, wait);
ZEND_METHOD(Pogo_Channel, __construct);
ZEND_METHOD(Pogo_Channel, init);
ZEND_METHOD(Pogo_Channel, push);
ZEND_METHOD(Pogo_Channel, pop);
ZEND_METHOD(Pogo_Channel, close);
ZEND_METHOD(Pogo_Runtime_Pool, __construct);
ZEND_METHOD(Pogo_Runtime_Pool, start);
ZEND_METHOD(Pogo_Runtime_Pool, shutdown);
ZEND_METHOD(Pogo_Runtime_Pool, submit);

static const zend_function_entry ext_functions[] = {
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo", "_gopogo_init"), zif_Pogo__gopogo_init, arginfo_Pogo__gopogo_init, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo", "_shm_read"), zif_Pogo__shm_read, arginfo_Pogo__shm_read, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo", "_shm_decode"), zif_Pogo__shm_decode, arginfo_Pogo__shm_decode, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo", "_shm_check"), zif_Pogo__shm_check, arginfo_Pogo__shm_check, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo", "start_worker_pool"), zif_Pogo_start_worker_pool, arginfo_Pogo_start_worker_pool, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo", "dispatch"), zif_Pogo_dispatch, arginfo_Pogo_dispatch, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo", "dispatch_task"), zif_Pogo_dispatch_task, arginfo_Pogo_dispatch_task, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo", "select"), zif_Pogo_select, arginfo_Pogo_select, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo", "async"), zif_Pogo_async, arginfo_Pogo_async, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo", "get_pool_stats"), zif_Pogo_get_pool_stats, arginfo_Pogo_get_pool_stats, 0, NULL, NULL)
	ZEND_FE_END
};

static const zend_function_entry class_Pogo_Future_methods[] = {
	ZEND_ME(Pogo_Future, __construct, arginfo_class_Pogo_Future___construct, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Future, await, arginfo_class_Pogo_Future_await, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Future, done, arginfo_class_Pogo_Future_done, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Future, cancel, arginfo_class_Pogo_Future_cancel, ZEND_ACC_PUBLIC)
	ZEND_FE_END
};

static const zend_function_entry class_Pogo_WaitGroup_methods[] = {
	ZEND_ME(Pogo_WaitGroup, __construct, arginfo_class_Pogo_WaitGroup___construct, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_WaitGroup, add, arginfo_class_Pogo_WaitGroup_add, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_WaitGroup, done, arginfo_class_Pogo_WaitGroup_done, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_WaitGroup, wait, arginfo_class_Pogo_WaitGroup_wait, ZEND_ACC_PUBLIC)
	ZEND_FE_END
};

static const zend_function_entry class_Pogo_Channel_methods[] = {
	ZEND_ME(Pogo_Channel, __construct, arginfo_class_Pogo_Channel___construct, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Channel, init, arginfo_class_Pogo_Channel_init, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Channel, push, arginfo_class_Pogo_Channel_push, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Channel, pop, arginfo_class_Pogo_Channel_pop, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Channel, close, arginfo_class_Pogo_Channel_close, ZEND_ACC_PUBLIC)
	ZEND_FE_END
};

static const zend_function_entry class_Pogo_Contract_Resettable_methods[] = {
	ZEND_RAW_FENTRY("reset", NULL, arginfo_class_Pogo_Contract_Resettable_reset, ZEND_ACC_PUBLIC|ZEND_ACC_ABSTRACT, NULL, NULL)
	ZEND_FE_END
};

static const zend_function_entry class_Pogo_Runtime_Pool_methods[] = {
	ZEND_ME(Pogo_Runtime_Pool, __construct, arginfo_class_Pogo_Runtime_Pool___construct, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Runtime_Pool, start, arginfo_class_Pogo_Runtime_Pool_start, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Runtime_Pool, shutdown, arginfo_class_Pogo_Runtime_Pool_shutdown, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Runtime_Pool, submit, arginfo_class_Pogo_Runtime_Pool_submit, ZEND_ACC_PUBLIC)
	ZEND_FE_END
};

static zend_class_entry *register_class_Pogo_Future(void)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Pogo", "Future", class_Pogo_Future_methods);
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
	zend_string *property_channel_class_Pogo_Channel = zend_string_init("Pogo\\Channel", sizeof("Pogo\\Channel")-1, 1);
	zend_declare_typed_property(class_entry, property_channel_name, &property_channel_default_value, ZEND_ACC_PRIVATE, NULL, (zend_type) ZEND_TYPE_INIT_CLASS(property_channel_class_Pogo_Channel, 0, MAY_BE_NULL));
	zend_string_release(property_channel_name);

	return class_entry;
}

static zend_class_entry *register_class_Pogo_WaitGroup(void)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Pogo", "WaitGroup", class_Pogo_WaitGroup_methods);
	class_entry = zend_register_internal_class_with_flags(&ce, NULL, 0);

	return class_entry;
}

static zend_class_entry *register_class_Pogo_Channel(void)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Pogo", "Channel", class_Pogo_Channel_methods);
	class_entry = zend_register_internal_class_with_flags(&ce, NULL, 0);

	return class_entry;
}

static zend_class_entry *register_class_Pogo_WorkerException(zend_class_entry *class_entry_Exception)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Pogo", "WorkerException", NULL);
	class_entry = zend_register_internal_class_with_flags(&ce, class_entry_Exception, 0);

	return class_entry;
}

static zend_class_entry *register_class_Pogo_TimeoutException(zend_class_entry *class_entry_Exception)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Pogo", "TimeoutException", NULL);
	class_entry = zend_register_internal_class_with_flags(&ce, class_entry_Exception, 0);

	return class_entry;
}

static zend_class_entry *register_class_Pogo_Contract_Resettable(void)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Pogo\\Contract", "Resettable", class_Pogo_Contract_Resettable_methods);
	class_entry = zend_register_internal_interface(&ce);

	return class_entry;
}

static zend_class_entry *register_class_Pogo_Runtime_Pool(void)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Pogo\\Runtime", "Pool", class_Pogo_Runtime_Pool_methods);
	class_entry = zend_register_internal_class_with_flags(&ce, NULL, 0);

	return class_entry;
}
