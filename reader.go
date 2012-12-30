package clz4

import (
	"bytes"
	"io"
	"io/ioutil"
)

type Reader struct {
	bytes.Buffer
	underlying io.Reader
}

// Note: this buffers the whole of the underlying reader, since this go code is
//       "lz4" unaware and does not know how to partially read an lz4 stream.
func (r *Reader) Read(data []byte) (int, error) {
	if len(data) < r.Buffer.Len() {
		// Not enough data in the buffer
		
		input, err := ioutil.ReadAll(r.underlying)
		if err != nil { return 0, err }
		
		output := []byte{}
		err = Uncompress(input, &output)
		if err != nil { return 0, err }
		
		r.Buffer.Write(output)		
	}
	return r.Buffer.Read(data)
}

func NewReader(r io.Reader) *Reader {
	return &Reader{underlying: r}
}

// Typecheck
var _ io.Reader = NewReader(&Reader{})

