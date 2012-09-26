package clz4

import (
	"strings"
	"testing"
)

func TestCompression(t *testing.T) {
	input := []byte(strings.Repeat("Hello world, this is quite something", 10))
	output := []byte{}
	err := Compress(input, &output)
	if err != nil {
		t.Fatalf("Compression failed: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("Output buffer is empty..")
	}
	t.Logf("Sizes: input=%d, output=%d, ratio=%.2f", len(input), len(output),
		float64(len(output))/float64(len(input)))
	decompressed := make([]byte, len(input))
	err = Uncompress(output, &decompressed)
	if err != nil {
		t.Fatalf("Decompression failed: %v", err)
	}
	if string(decompressed) != string(input) {
		t.Fatalf("Decompressed output != input: %q != %q", decompressed, input)
	}

	decompressed = make([]byte, 0, len(input))
	err = UncompressUnknownOutputSize(output, &decompressed)
	if err != nil {
		t.Fatalf("Decompression failed: %v", err)
	}
	if string(decompressed) != string(input) {
		t.Fatalf("Decompressed output != input: %q != %q", decompressed, input)
	}

	decompressed = make([]byte, 0, len(input)-1)
	err = UncompressUnknownOutputSize(output, &decompressed)
	if err == nil {
		t.Fatalf("UncompressUnknownOutputSize didn't fail when it should have")
	}
}
