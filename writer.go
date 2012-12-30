package clz4

import (
	"fmt"
	"io"
)

type Writer struct {
	io.Writer
}

func (w *Writer) Write(data []byte) (int, error) {
	compressed_bytes := []byte{}
	Compress(data, &compressed_bytes)
	n, err := w.Writer.Write(compressed_bytes)
	if n != len(compressed_bytes) {
		return n, fmt.Errorf("short write: %d - expected %d", n, len(compressed_bytes))
	}
	return len(data), err
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w}
}

// Typecheck
var _ io.Writer = NewWriter(&Writer{})

