package stream

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"log"

	"github.com/eaburns/bit"
	// "github.com/kr/pretty"
	"github.com/pwaller/go-clz4"
)

var _ = log.Prefix // TODO(pwaller): remove me

var (
	ErrBadMagic                   = errors.New("Bad magic number, expected 0x04224D18")
	ErrUnsupportedVersion         = errors.New("Unsupported stream version")
	ErrUnsupportedBlockDependence = errors.New("Unsupported stream format: (Block Independence == false)")
)

type Reader struct {
	underlying   io.ReadSeeker
	currentBlock io.Reader
	bitReader    *bit.Reader
	header       Header
}

type Header struct {
	magic   uint32
	version uint8
	flags   struct {
		blockIndependence, blockChecksum, streamSize,
		streamChecksum, reserved, presetDictionary bool
	}
	blockMaximum     uint64
	streamSize       int64
	presetDictionary uint32
	headerChecksum   uint8
}

func NewReader(u io.ReadSeeker) (*Reader, error) {
	r := &Reader{
		underlying: u,
		bitReader:  bit.NewReader(u),
	}
	r.header.streamSize = -1

	err := r.readHeader()
	if err != nil {
		return nil, err
	}
	err = r.checkHeader()
	if err != nil {
		return nil, err
	}
	err = r.NextBlock()
	return r, err
}

func (r *Reader) NextBlock() error {
	var size uint32
	err := binary.Read(r.underlying, binary.LittleEndian, &size)
	if err != nil {
		return err
	}
	if size == 0 {
		return io.EOF
	}
	buf, err := ioutil.ReadAll(io.LimitReader(r.underlying, int64(size)))
	if err != nil {
		log.Println("Not able to read")
		return err
	}
	out := make([]byte, 0, size*10)
	clz4.UncompressUnknownOutputSize(buf, &out)
	r.currentBlock = bytes.NewBuffer(out)
	return nil
}

// func (r *Reader) FindBlocks() error {

// 	var size uint32 = 1
// 	for size != 0 {
// 		err := binary.Read(r.underlying, binary.LittleEndian, &size)
// 		if err != nil {
// 			return err
// 		}
// 		log.Println("Size:", size)
// 		_, err = r.underlying.Seek(int64(size), os.SEEK_CUR)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	return nil
// }

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.currentBlock.Read(p)
	if err == io.EOF {
		return n, r.NextBlock()
	}
	return
}

func (r *Reader) readHeader() error {

	h := &r.header

	sizes := []uint{32, 2, 1, 1, 1, 1, 1, 1, 1, 3, 4}
	header, err := r.bitReader.ReadFields(sizes...)
	if err != nil {
		return err
	}

	h.magic = uint32(header[0])
	h.version = uint8(header[1])

	h.flags.blockIndependence = header[2] == 1
	h.flags.blockChecksum = header[3] == 1
	h.flags.streamSize = header[4] == 1
	h.flags.streamChecksum = header[5] == 1
	h.flags.reserved = header[6] == 1
	h.flags.presetDictionary = header[7] == 1

	_ = header[8] // reserved field
	h.blockMaximum = header[9]
	_ = header[10] // reserved field

	if h.flags.streamSize {
		s, err := r.bitReader.Read(64)
		// Hmm, bitreader is Big Endian :(
		// h.streamSize = int64(s)
		h.streamSize = 0 // TODO
		_ = s
		if err != nil {
			return err
		}
	}

	if h.flags.presetDictionary {
		p, err := r.bitReader.Read(32)
		h.presetDictionary = uint32(p)
		if err != nil {
			return err
		}
	}

	s, err := r.bitReader.Read(8)
	h.headerChecksum = uint8(s)
	if err != nil {
		return err
	}

	return nil
}

func (r *Reader) checkHeader() error {
	h := r.header
	if h.magic != 0x04224D18 {
		return ErrBadMagic
	}
	if h.version != 1 {
		return ErrUnsupportedVersion
	}
	fl := h.flags

	if fl.blockIndependence == false {
		// If support for streaming inter-dependent blocks is added to go-clz4,
		// this needs to be removed.
		return ErrUnsupportedBlockDependence
	}

	return nil
}
