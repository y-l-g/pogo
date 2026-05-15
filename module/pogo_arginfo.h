/* This is a generated-compatible file, edit pogo.stub.php and regenerate when tooling is available. */

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_pogo_dispatch, 0, 1, IS_LONG, 0)
	ZEND_ARG_TYPE_INFO(0, class, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, args, IS_ARRAY, 0, "[]")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, pool, IS_STRING, 0, "\"default\"")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_pogo_await, 0, 1, IS_MIXED, 0)
	ZEND_ARG_TYPE_INFO(0, handle, IS_LONG, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, timeout, IS_DOUBLE, 0, "5.0")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_pogo_pool_size, 0, 0, IS_LONG, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, pool, IS_STRING, 0, "\"default\"")
ZEND_END_ARG_INFO()

ZEND_FUNCTION(pogo_dispatch);
ZEND_FUNCTION(pogo_await);
ZEND_FUNCTION(pogo_pool_size);

static const zend_function_entry ext_functions[] = {
	ZEND_FE(pogo_dispatch, arginfo_pogo_dispatch)
	ZEND_FE(pogo_await, arginfo_pogo_await)
	ZEND_FE(pogo_pool_size, arginfo_pogo_pool_size)
	ZEND_FE_END
};
