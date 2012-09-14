package clz4

// #cgo CFLAGS: -O3
// #include "src/lz4.c"
import "C"

import (
	"log"
	"unsafe"
)

func LZ4_uncompress_unknownOutputSize(input []byte, output *[]byte) {
	ip := (*C.char)(unsafe.Pointer(&(input)[0]))
	op := (*C.char)(unsafe.Pointer(&(*output)[0]))
	resultlen := C.LZ4_uncompress_unknownOutputSize(
		ip, op, C.int(len(input)), C.int(cap(*output)))
	if int(resultlen) >= cap(*output) {
		log.Panic("Looks like decompression buffer wasn't big enough!")
	}
	(*output) = (*output)[0:resultlen]
}
