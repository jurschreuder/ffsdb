package ffsdb

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
)

type Ffsdb struct {
	Path string

	save32Bit  bool
	sliceLen   int
	sliceLenBs int64
	fd         *os.File
	fdAt       *os.File
	reader     *bufio.Reader
	writer     *bufio.Writer
	buffer     []byte
	isFlushed  bool
}

func NewFfsdb(path string, sliceLen int, removeOld, save32bit bool) (*Ffsdb, error) {
	if removeOld {
		os.Remove(path)
	}
	floatBit := 8
	if save32bit {
		floatBit = 4
	}
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return &Ffsdb{}, nil
	}
	fdAt, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return &Ffsdb{}, nil
	}
	fdb := Ffsdb{
		Path:       path,
		save32Bit:  save32bit,
		sliceLen:   sliceLen,
		sliceLenBs: int64(sliceLen * floatBit),
		fd:         fd,
		fdAt:       fdAt,
		reader:     bufio.NewReaderSize(fd, sliceLen*floatBit),
		writer:     bufio.NewWriter(fd),
		buffer:     make([]byte, sliceLen*floatBit),
		isFlushed:  true,
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

func Bytes32ToFloat64(bytes []byte) float64 {
	bits := binary.BigEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float64(float)
}

func Float64ToBytes32(float float64) []byte {
	bits := math.Float32bits(float32(float))
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, bits)
	return bytes
}

func Float64SliceToBytes32(fs []float64, bs []byte) {
	bi := 0
	for _, f := range fs {
		fbs := Float64ToBytes32(f)
		for i := 0; i < 4; i++ {
			bs[bi+i] = fbs[i]
		}
		bi += 4
	}
}

func Bytes32ToFloat64Slice(bs []byte) []float64 {
	fs := make([]float64, len(bs)/4)
	fbs := make([]byte, 4)
	n := 0
	for i := 0; i < len(bs); i += 4 {
		for j := 0; j < 4; j++ {
			fbs[j] = bs[i+j]
		}
		fs[n] = Bytes32ToFloat64(fbs)
		n++
	}
	return fs
}

func (fdb *Ffsdb) Close() {
	fdb.fd.Close()
	fdb.fdAt.Close()
}

func (fdb *Ffsdb) Rewind() {
	if !fdb.isFlushed {
		fdb.Flush()
	}
	fdb.fd.Seek(0, 0)
}

func (fdb *Ffsdb) Seek(id int64) error {
	_, err := fdb.fd.Seek(id*fdb.sliceLenBs, 0)
	return err
}

func (fdb *Ffsdb) ReadId(id int64) ([]float64, error) {
	if !fdb.isFlushed {
		fdb.Flush()
	}
	err := fdb.Seek(id)
	if err != nil {
		return []float64{}, err
	}
	_, err = fdb.reader.Read(fdb.buffer)
	if err != nil {
		return []float64{}, err
	}
	if fdb.save32Bit {
		return Bytes32ToFloat64Slice(fdb.buffer), nil
	}
	return BytesToFloat64Slice(fdb.buffer), nil
}

func (fdb *Ffsdb) ReadNext() ([]float64, bool) {
	if !fdb.isFlushed {
		fdb.Flush()
	}
	_, err := fdb.reader.Read(fdb.buffer)
	if err != nil {
		return []float64{}, false
	}
	if fdb.save32Bit {
		return Bytes32ToFloat64Slice(fdb.buffer), true
	}
	return BytesToFloat64Slice(fdb.buffer), true
}

func (fdb *Ffsdb) Add(vals []float64) error {
	if len(vals) != fdb.sliceLen {
		return errors.New(fmt.Sprint("vals length was ", len(vals), " expected ", fdb.sliceLen))
	}
	return fdb.AddUnsafe(vals)
}

func (fdb *Ffsdb) AddUnsafe(vals []float64) error {
	fdb.isFlushed = false
	if fdb.save32Bit {
		Float64SliceToBytes32(vals, fdb.buffer)
	} else {
		Float64SliceToBytes(vals, fdb.buffer)
	}
	_, err := fdb.writer.Write(fdb.buffer)
	return err
}

func (fdb *Ffsdb) Update(id int64, vals []float64) error {
	if len(vals) != fdb.sliceLen {
		return errors.New(fmt.Sprint("vals length was ", len(vals), " expected ", fdb.sliceLen))
	}
	return fdb.UpdateUnsafe(id, vals)
}

func (fdb *Ffsdb) UpdateUnsafe(id int64, vals []float64) error {
	if fdb.save32Bit {
		Float64SliceToBytes32(vals, fdb.buffer)
	} else {
		Float64SliceToBytes(vals, fdb.buffer)
	}
	_, err := fdb.fdAt.WriteAt(fdb.buffer, id*fdb.sliceLenBs)
	return err
}

func (fdb *Ffsdb) Flush() {
	fdb.writer.Flush()
	fdb.isFlushed = true
}
