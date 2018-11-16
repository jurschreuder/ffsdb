package ffsdb

import (
	"bufio"
	"encoding/binary"
	"math"
	"os"
)

type Ffsdb struct {
	Path string

	sliceLen  int
	fd        *os.File
	reader    *bufio.Reader
	writer    *bufio.Writer
	buffer    []byte
	isFlushed bool
}

func NewFfsdb(path string, sliceLen int, removeOld bool) (*Ffsdb, error) {
	if removeOld {
		os.Remove(path)
	}
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return &Ffsdb{}, nil
	}
	fdb := Ffsdb{
		Path:      path,
		sliceLen:  sliceLen,
		fd:        fd,
		reader:    bufio.NewReaderSize(fd, sliceLen*8),
		writer:    bufio.NewWriter(fd),
		buffer:    make([]byte, sliceLen*8),
		isFlushed: true,
	}
	return &fdb, nil
}

func BytesToFloat64(bytes []byte) float64 {
	bits := binary.BigEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, bits)
	return bytes
}

func Float64SliceToBytes(fs []float64, bs []byte) {
	bi := 0
	for _, f := range fs {
		fbs := Float64ToBytes(f)
		for i := 0; i < 8; i++ {
			bs[bi+i] = fbs[i]
		}
		bi += 8
	}
}

func BytesToFloat64Slice(bs []byte) []float64 {
	fs := make([]float64, len(bs)/8)
	fbs := make([]byte, 8)
	n := 0
	for i := 0; i < len(bs); i += 8 {
		for j := 0; j < 8; j++ {
			fbs[j] = bs[i+j]
		}
		fs[n] = BytesToFloat64(fbs)
		n++
	}
	return fs
}

func (fdb *Ffsdb) Rewind() {
	if !fdb.isFlushed {
		fdb.Flush()
	}
	fdb.fd.Seek(0, 0)
}

func (fdb *Ffsdb) ReadNext() ([]float64, bool) {
	if !fdb.isFlushed {
		fdb.Flush()
	}
	_, err := fdb.reader.Read(fdb.buffer)
	//_, err := fdb.fd.Read(fdb.buffer)
	if err != nil {
		return []float64{}, false
	}
	//	if n != fdb.sliceLen*8 {
	//		return []float64{}, false
	//	}
	return BytesToFloat64Slice(fdb.buffer), true
}

func (fdb *Ffsdb) Add(vals []float64) error {
	fdb.isFlushed = false
	Float64SliceToBytes(vals, fdb.buffer)
	_, err := fdb.writer.Write(fdb.buffer)
	return err
}

func (fdb *Ffsdb) Flush() {
	fdb.writer.Flush()
	fdb.isFlushed = true
}
