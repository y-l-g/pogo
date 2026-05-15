package pogo

/*
#include <stdlib.h>
#include "pogo.h"
*/
import "C"
import (
	"unsafe"

	"github.com/dunglas/frankenphp"
)

func init() {
	frankenphp.RegisterExtension(unsafe.Pointer(&C.pogo_module_entry))
}
