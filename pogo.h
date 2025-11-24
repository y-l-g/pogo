#ifndef _POGO_H
#define _POGO_H

#include <php.h>
#include <stdint.h>
#include <Zend/zend_types.h>
#include <stddef.h>
#include <stdbool.h>

extern zend_module_entry pogo_module_entry;

/**
 * Custom object structure for PHP classes wrapping Go resources.
 * 
 * @var go_handle   The uintptr_t handle returned by cgo.NewHandle().
 * @var owns_handle If true, the destructor will free the Go handle.
 *                  Channels own the handle. Futures borrow it.
 * @var std         Standard Zend object header.
 */
typedef struct
{
    uintptr_t go_handle;
    bool owns_handle;
    zend_object std;
} pogo_object;

/**
 * Helper to cast a standard zend_object pointer to our custom struct.
 */
static inline pogo_object *pogo_object_from_obj(zend_object *obj)
{
    return (pogo_object *)((char *)(obj)-offsetof(pogo_object, std));
}

#endif