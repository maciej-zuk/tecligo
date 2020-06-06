package temap

// #cgo CFLAGS: -Wall
// #include <stdlib.h>
// #include "lzx.h"
import "C"

import "unsafe"

func LZXDecompress(data []byte) []byte {
	indata := C.CBytes(data)
	inlen := C.int(len(data))
	outlen := C.int(0)
	outdata := C.LZXDecompressFull((*C.uchar)(indata), inlen, (*C.int)(unsafe.Pointer(&outlen)))
	output := C.GoBytes(unsafe.Pointer(outdata), outlen)
	C.free(indata)
	C.free(unsafe.Pointer(outdata))
	return output
}
