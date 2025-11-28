/* This is a generated file, edit the .stub.php file instead.
 * Stub hash: df6899199c1d49deed581f22d720f98d4ca99528 */

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo_Internal__gopogo_init, 0, 0, IS_VOID, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo_Internal__shm_read, 0, 3, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO(0, fd, IS_LONG, 0)
	ZEND_ARG_TYPE_INFO(0, offset, IS_LONG, 0)
	ZEND_ARG_TYPE_INFO(0, length, IS_LONG, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo_Internal__shm_decode, 0, 3, IS_MIXED, 0)
	ZEND_ARG_TYPE_INFO(0, fd, IS_LONG, 0)
	ZEND_ARG_TYPE_INFO(0, offset, IS_LONG, 0)
	ZEND_ARG_TYPE_INFO(0, length, IS_LONG, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo_Internal__shm_check, 0, 1, _IS_BOOL, 0)
	ZEND_ARG_TYPE_INFO(0, fd, IS_LONG, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo_Internal_start_worker_pool, 0, 0, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, entrypoint, IS_STRING, 0, "\"job_runner.php\"")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, minWorkers, IS_LONG, 0, "4")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, maxWorkers, IS_LONG, 0, "8")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, maxJobs, IS_LONG, 0, "0")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, options, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo_Internal_dispatch, 0, 2, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO(0, workerName, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO(0, payload, IS_ARRAY, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_OBJ_INFO_EX(arginfo_Pogo_Internal_dispatch_task, 0, 1, Pogo\\Internal\\Future, 0)
	ZEND_ARG_TYPE_INFO(0, taskName, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, payload, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo_Internal_select, 0, 1, IS_ARRAY, 1)
	ZEND_ARG_TYPE_INFO(0, cases, IS_ARRAY, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, timeout, IS_DOUBLE, 1, "null")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_OBJ_INFO_EX(arginfo_Pogo_Internal_async, 0, 1, Pogo\\Internal\\Future, 0)
	ZEND_ARG_TYPE_INFO(0, class, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, args, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo_Internal_get_pool_stats, 0, 0, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, poolID, IS_LONG, 0, "0")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Pogo_Internal_version, 0, 0, IS_STRING, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_INFO_EX(arginfo_class_Pogo_Internal_Future___construct, 0, 0, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Pogo_Internal_Future_await, 0, 0, IS_STRING, 1)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, timeout, IS_DOUBLE, 1, "null")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Pogo_Internal_Future_done, 0, 0, _IS_BOOL, 0)
ZEND_END_ARG_INFO()

#define arginfo_class_Pogo_Internal_Future_cancel arginfo_class_Pogo_Internal_Future_done

#define arginfo_class_Pogo_Internal_WaitGroup___construct arginfo_class_Pogo_Internal_Future___construct

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Pogo_Internal_WaitGroup_add, 0, 0, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, delta, IS_LONG, 0, "1")
ZEND_END_ARG_INFO()

#define arginfo_class_Pogo_Internal_WaitGroup_done arginfo_Pogo_Internal__gopogo_init

#define arginfo_class_Pogo_Internal_WaitGroup_wait arginfo_Pogo_Internal__gopogo_init

#define arginfo_class_Pogo_Internal_Channel___construct arginfo_class_Pogo_Internal_Future___construct

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Pogo_Internal_Channel_init, 0, 0, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, capacity, IS_LONG, 0, "0")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Pogo_Internal_Channel_push, 0, 1, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO(0, value, IS_STRING, 0)
ZEND_END_ARG_INFO()

#define arginfo_class_Pogo_Internal_Channel_pop arginfo_Pogo_Internal_version

#define arginfo_class_Pogo_Internal_Channel_close arginfo_Pogo_Internal__gopogo_init

ZEND_BEGIN_ARG_INFO_EX(arginfo_class_Pogo_Internal_Pool___construct, 0, 0, 1)
	ZEND_ARG_TYPE_INFO(0, entrypoint, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, minWorkers, IS_LONG, 0, "1")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, maxWorkers, IS_LONG, 0, "8")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, maxJobs, IS_LONG, 0, "0")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, options, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

#define arginfo_class_Pogo_Internal_Pool_start arginfo_Pogo_Internal__gopogo_init

#define arginfo_class_Pogo_Internal_Pool_shutdown arginfo_Pogo_Internal__gopogo_init

ZEND_BEGIN_ARG_WITH_RETURN_OBJ_INFO_EX(arginfo_class_Pogo_Internal_Pool_submit, 0, 1, Pogo\\Internal\\Future, 0)
	ZEND_ARG_TYPE_INFO(0, jobClass, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, args, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_Pogo_Internal_Pool_id, 0, 0, IS_LONG, 0)
ZEND_END_ARG_INFO()

ZEND_FUNCTION(Pogo_Internal__gopogo_init);
ZEND_FUNCTION(Pogo_Internal__shm_read);
ZEND_FUNCTION(Pogo_Internal__shm_decode);
ZEND_FUNCTION(Pogo_Internal__shm_check);
ZEND_FUNCTION(Pogo_Internal_start_worker_pool);
ZEND_FUNCTION(Pogo_Internal_dispatch);
ZEND_FUNCTION(Pogo_Internal_dispatch_task);
ZEND_FUNCTION(Pogo_Internal_select);
ZEND_FUNCTION(Pogo_Internal_async);
ZEND_FUNCTION(Pogo_Internal_get_pool_stats);
ZEND_FUNCTION(Pogo_Internal_version);
ZEND_METHOD(Pogo_Internal_Future, __construct);
ZEND_METHOD(Pogo_Internal_Future, await);
ZEND_METHOD(Pogo_Internal_Future, done);
ZEND_METHOD(Pogo_Internal_Future, cancel);
ZEND_METHOD(Pogo_Internal_WaitGroup, __construct);
ZEND_METHOD(Pogo_Internal_WaitGroup, add);
ZEND_METHOD(Pogo_Internal_WaitGroup, done);
ZEND_METHOD(Pogo_Internal_WaitGroup, wait);
ZEND_METHOD(Pogo_Internal_Channel, __construct);
ZEND_METHOD(Pogo_Internal_Channel, init);
ZEND_METHOD(Pogo_Internal_Channel, push);
ZEND_METHOD(Pogo_Internal_Channel, pop);
ZEND_METHOD(Pogo_Internal_Channel, close);
ZEND_METHOD(Pogo_Internal_Pool, __construct);
ZEND_METHOD(Pogo_Internal_Pool, start);
ZEND_METHOD(Pogo_Internal_Pool, shutdown);
ZEND_METHOD(Pogo_Internal_Pool, submit);
ZEND_METHOD(Pogo_Internal_Pool, id);

static const zend_function_entry ext_functions[] = {
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo\\Internal", "_gopogo_init"), zif_Pogo_Internal__gopogo_init, arginfo_Pogo_Internal__gopogo_init, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo\\Internal", "_shm_read"), zif_Pogo_Internal__shm_read, arginfo_Pogo_Internal__shm_read, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo\\Internal", "_shm_decode"), zif_Pogo_Internal__shm_decode, arginfo_Pogo_Internal__shm_decode, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo\\Internal", "_shm_check"), zif_Pogo_Internal__shm_check, arginfo_Pogo_Internal__shm_check, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo\\Internal", "start_worker_pool"), zif_Pogo_Internal_start_worker_pool, arginfo_Pogo_Internal_start_worker_pool, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo\\Internal", "dispatch"), zif_Pogo_Internal_dispatch, arginfo_Pogo_Internal_dispatch, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo\\Internal", "dispatch_task"), zif_Pogo_Internal_dispatch_task, arginfo_Pogo_Internal_dispatch_task, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo\\Internal", "select"), zif_Pogo_Internal_select, arginfo_Pogo_Internal_select, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo\\Internal", "async"), zif_Pogo_Internal_async, arginfo_Pogo_Internal_async, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo\\Internal", "get_pool_stats"), zif_Pogo_Internal_get_pool_stats, arginfo_Pogo_Internal_get_pool_stats, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Pogo\\Internal", "version"), zif_Pogo_Internal_version, arginfo_Pogo_Internal_version, 0, NULL, NULL)
	ZEND_FE_END
};

static const zend_function_entry class_Pogo_Internal_Future_methods[] = {
	ZEND_ME(Pogo_Internal_Future, __construct, arginfo_class_Pogo_Internal_Future___construct, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Internal_Future, await, arginfo_class_Pogo_Internal_Future_await, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Internal_Future, done, arginfo_class_Pogo_Internal_Future_done, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Internal_Future, cancel, arginfo_class_Pogo_Internal_Future_cancel, ZEND_ACC_PUBLIC)
	ZEND_FE_END
};

static const zend_function_entry class_Pogo_Internal_WaitGroup_methods[] = {
	ZEND_ME(Pogo_Internal_WaitGroup, __construct, arginfo_class_Pogo_Internal_WaitGroup___construct, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Internal_WaitGroup, add, arginfo_class_Pogo_Internal_WaitGroup_add, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Internal_WaitGroup, done, arginfo_class_Pogo_Internal_WaitGroup_done, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Internal_WaitGroup, wait, arginfo_class_Pogo_Internal_WaitGroup_wait, ZEND_ACC_PUBLIC)
	ZEND_FE_END
};

static const zend_function_entry class_Pogo_Internal_Channel_methods[] = {
	ZEND_ME(Pogo_Internal_Channel, __construct, arginfo_class_Pogo_Internal_Channel___construct, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Internal_Channel, init, arginfo_class_Pogo_Internal_Channel_init, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Internal_Channel, push, arginfo_class_Pogo_Internal_Channel_push, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Internal_Channel, pop, arginfo_class_Pogo_Internal_Channel_pop, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Internal_Channel, close, arginfo_class_Pogo_Internal_Channel_close, ZEND_ACC_PUBLIC)
	ZEND_FE_END
};

static const zend_function_entry class_Pogo_Internal_Pool_methods[] = {
	ZEND_ME(Pogo_Internal_Pool, __construct, arginfo_class_Pogo_Internal_Pool___construct, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Internal_Pool, start, arginfo_class_Pogo_Internal_Pool_start, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Internal_Pool, shutdown, arginfo_class_Pogo_Internal_Pool_shutdown, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Internal_Pool, submit, arginfo_class_Pogo_Internal_Pool_submit, ZEND_ACC_PUBLIC)
	ZEND_ME(Pogo_Internal_Pool, id, arginfo_class_Pogo_Internal_Pool_id, ZEND_ACC_PUBLIC)
	ZEND_FE_END
};

static zend_class_entry *register_class_Pogo_Internal_Future(void)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Pogo\\Internal", "Future", class_Pogo_Internal_Future_methods);
	class_entry = zend_register_internal_class_with_flags(&ce, NULL, 0);

	return class_entry;
}

static zend_class_entry *register_class_Pogo_Internal_WaitGroup(void)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Pogo\\Internal", "WaitGroup", class_Pogo_Internal_WaitGroup_methods);
	class_entry = zend_register_internal_class_with_flags(&ce, NULL, 0);

	return class_entry;
}

static zend_class_entry *register_class_Pogo_Internal_Channel(void)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Pogo\\Internal", "Channel", class_Pogo_Internal_Channel_methods);
	class_entry = zend_register_internal_class_with_flags(&ce, NULL, 0);

	return class_entry;
}

static zend_class_entry *register_class_Pogo_Internal_Pool(void)
{
	zend_class_entry ce, *class_entry;

	INIT_NS_CLASS_ENTRY(ce, "Pogo\\Internal", "Pool", class_Pogo_Internal_Pool_methods);
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
