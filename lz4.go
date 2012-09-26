// Fast compression using C lz4 source code.
// The functions in this package should be robust against malicious inputs.
package clz4

// #cgo CFLAGS: -O3
// int LZ4_compressBound(int isize) { return (isize + (isize/255) + 16); }
// #include "src/lz4.h"
// #include "src/lz4.c"
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Get a char pointer to the first byte of a slice
func charp(input *[]byte) *C.char {
	result_header := (*reflect.SliceHeader)(unsafe.Pointer(input))
	return (*C.char)(unsafe.Pointer(result_header.Data))
}

// Uncompress with an known output size. `len(*output)` should be equal to the
// length of the uncompressed output.
func Uncompress(input []byte, output *[]byte) error {
	ip, op := charp(&input), charp(output)

	resultlen := int(C.LZ4_uncompress(ip, op, C.int(len(*output))))

	if resultlen < len(input) {
		return fmt.Errorf("Decompression didn't read all the input "+
			" Expected %d bytes, read %d", len(input), resultlen)
	}
	return nil
}

// Uncompress with an unknown output size. The destination buffer should have
// a cap() of sufficient size to hold the output, otherwise an error is returned.
func UncompressUnknownOutputSize(input []byte, output *[]byte) error {
	ip, op := charp(&input), charp(output)

	resultlen := int(C.LZ4_uncompress_unknownOutputSize(
		ip, op, C.int(len(input)), C.int(cap(*output))))

	if resultlen < 0 {
		return fmt.Errorf("decompression failed: resultlen=%d", resultlen)
	}

	*output = (*output)[0:resultlen]
	if len(*output) != resultlen {
		return fmt.Errorf("Failed to resize destination buffer")
	}
	return nil
}

func CompressBound(input []byte) int {
	return int(C.LZ4_compressBound(C.int(len(input))))
}

// Compresses `input` and puts the content in `output`. `output` will be 
// re-allocated if there is not sufficient space to store the maximum length
// compressed output.
func Compress(input []byte, output *[]byte) error {
	cb := CompressBound(input)
	if cap(*output) < cb {
		*output = make([]byte, 0, cb)
	}

	ip, op := charp(&input), charp(output)
	resultlen := int(C.LZ4_compress_limitedOutput(ip, op,
		C.int(len(input)), C.int(cap(*output))))
	if resultlen > cap(*output) {
		return fmt.Errorf("LZ4 overran compression buffer. This shouldn't happen. "+
			"Expected: %d, got %d", cap(*output), resultlen)
	}
	*output = (*output)[0:resultlen]
	return nil
}
