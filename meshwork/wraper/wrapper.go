package wrapper

import (
	"encoding/binary"
	"io"
	"math"
)

// Wrapper is a binary packer helps you pack data into an io.Writer.
type Wrapper struct {
	writer io.Writer
	endian binary.ByteOrder
	err    error
}

// NewWrapper returns a *Packer hold an io.Writer. User must provide the byte order explicitly.
func NewWrapper(endian binary.ByteOrder, writer io.Writer) *Wrapper {
	return &Wrapper{
		writer: writer,
		endian: endian,
	}
}

// Error returns an error if any errors exists
func (w *Wrapper) Error() error {
	return w.err
}

// PushByte write a single byte into writer.
func (w *Wrapper) PushByte(b byte) *Wrapper {
	return w.errFilter(func() {
		_, w.err = w.writer.Write([]byte{b})
	})
}

// PushBytes write a bytes array into writer.
func (w *Wrapper) PushBytes(bytes []byte) *Wrapper {
	return w.errFilter(func() {
		_, w.err = w.writer.Write(bytes)
	})
}

// PushUint8 write a uint8 into writer.
func (w *Wrapper) PushUint8(i uint8) *Wrapper {
	return w.errFilter(func() {
		_, w.err = w.writer.Write([]byte{byte(i)})
	})
}

// PushUint16 write a uint16 into writer.
func (w *Wrapper) PushUint16(i uint16) *Wrapper {
	return w.errFilter(func() {
		buffer := make([]byte, 2)
		w.endian.PutUint16(buffer, i)
		_, w.err = w.writer.Write(buffer)
	})
}

//PushInt16 write a int16 into writer.
func (w *Wrapper) PushInt16(i int16) *Wrapper {
	return w.PushUint16(uint16(i))
}

// PushUint32 write a uint32 into writer.
func (w *Wrapper) PushUint32(i uint32) *Wrapper {
	return w.errFilter(func() {
		buffer := make([]byte, 4)
		w.endian.PutUint32(buffer, i)
		_, w.err = w.writer.Write(buffer)
	})
}

// PushInt32 write a int32 into writer.
func (w *Wrapper) PushInt32(i int32) *Wrapper {
	return w.PushUint32(uint32(i))
}

// PushUint64 write a uint64 into writer.
func (w *Wrapper) PushUint64(i uint64) *Wrapper {
	return w.errFilter(func() {
		buffer := make([]byte, 8)
		w.endian.PutUint64(buffer, i)
		_, w.err = w.writer.Write(buffer)
	})
}

// PushInt64 write a int64 into writer.
func (w *Wrapper) PushInt64(i int64) *Wrapper {
	return w.PushUint64(uint64(i))
}

// PushFloat32 write a float32 into writer.
func (w *Wrapper) PushFloat32(i float32) *Wrapper {
	return w.PushUint32(math.Float32bits(i))
}

// PushFloat64 write a float64 into writer.
func (w *Wrapper) PushFloat64(i float64) *Wrapper {
	return w.PushUint64(math.Float64bits(i))
}

// PushString write a string into writer.
func (w *Wrapper) PushString(s string) *Wrapper {
	return w.errFilter(func() {
		_, w.err = w.writer.Write([]byte(s))
	})
}

func (w *Wrapper) errFilter(f func()) *Wrapper {
	if w.err == nil {
		f()
	}
	return w
}
