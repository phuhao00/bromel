package wrapper

import (
	"encoding/binary"
	"io"
	"math"
	"unsafe"
)

// UnWrapper helps you unpack binary data from an io.Reader.
type UnWrapper struct {
	reader io.Reader
	endian binary.ByteOrder
	err    error
}

// NewUnWrapper returns a *Unpacker which hold an io.Reader. User must provide the byte order explicitly.
func NewUnWrapper(endian binary.ByteOrder, reader io.Reader) *UnWrapper {
	return &UnWrapper{
		reader: reader,
		endian: endian,
	}
}

// Error returns an error if any errors exists
func (u *UnWrapper) Error() error {
	return u.err
}

// ShiftByte fetch the first byte in io.Reader. Returns a byte and an error if
// exists.
func (u *UnWrapper) ShiftByte() (byte, error) {
	buffer := make([]byte, 1)
	_, err := u.reader.Read(buffer)
	return buffer[0], err
}

// FetchByte fetch the first byte in io.Reader and set to b.
func (u *UnWrapper) FetchByte(b *byte) *UnWrapper {
	return u.errFilter(func() {
		*b, u.err = u.ShiftByte()
	})
}

// ShiftBytes fetch n bytes in io.Reader. Returns a byte array and an error if
// exists.
func (u *UnWrapper) ShiftBytes(_n uint64) ([]byte, error) {
	buf := make([]byte, _n)
	_, err := io.ReadFull(u.reader, buf)
	return buf, err
}

// FetchBytes read n bytes and set to bytes.
func (u *UnWrapper) FetchBytes(n uint64, bytes *[]byte) *UnWrapper {
	return u.errFilter(func() {
		*bytes, u.err = u.ShiftBytes(n)
	})
}

// ShiftUint8 fetch 1 byte in io.Reader and covert it to uint8
func (u *UnWrapper) ShiftUint8() (uint8, error) {
	buffer := make([]byte, 1)
	_, err := u.reader.Read(buffer)
	return uint8(buffer[0]), err
}

// FetchUint8 read 1 byte, convert it to uint8 and set it to i.
func (u *UnWrapper) FetchUint8(i *uint8) *UnWrapper {
	return u.errFilter(func() {
		*i, u.err = u.ShiftUint8()
	})
}

// ShiftUint16 fetch 2 bytes in io.Reader and convert it to uint16.
func (u *UnWrapper) ShiftUint16() (uint16, error) {
	buffer := make([]byte, 2)
	if _, err := u.reader.Read(buffer); err != nil {
		return 0, err
	}
	return u.endian.Uint16(buffer), nil
}

// FetchUint16 read 2 bytes, convert it to uint16 and set it to i.
func (u *UnWrapper) FetchUint16(i *uint16) *UnWrapper {
	return u.errFilter(func() {
		*i, u.err = u.ShiftUint16()
	})
}

// ShiftInt16 fetch 2 bytes in io.Reader and convert it to int16.
func (u *UnWrapper) ShiftInt16() (int16, error) {
	i, err := u.ShiftUint16()
	return int16(i), err
}

// FetchInt16 read 2 bytes, convert it to int16 and set it to i.
func (u *UnWrapper) FetchInt16(i *int16) *UnWrapper {
	return u.FetchUint16((*uint16)(unsafe.Pointer(i)))
}

// ShiftUint32 fetch 4 bytes in io.Reader and convert it to uint32.
func (u *UnWrapper) ShiftUint32() (uint32, error) {
	buffer := make([]byte, 4)
	if _, err := u.reader.Read(buffer); err != nil {
		return 0, err
	}
	return u.endian.Uint32(buffer), nil
}

// ShiftInt32 fetch 4 bytes in io.Reader and convert it to int32.
func (u *UnWrapper) ShiftInt32() (int32, error) {
	i, err := u.ShiftUint32()
	return int32(i), err
}

// FetchUint32 read 4 bytes, convert it to uint32 and set it to i.
func (u *UnWrapper) FetchUint32(i *uint32) *UnWrapper {
	return u.errFilter(func() {
		*i, u.err = u.ShiftUint32()
	})
}

// FetchInt32 read 4 bytes, convert it to int32 and set it to i.
func (u *UnWrapper) FetchInt32(i *int32) *UnWrapper {
	return u.FetchUint32((*uint32)(unsafe.Pointer(i)))
}

// ShiftUint64 fetch 8 bytes in io.Reader and convert it to uint64.
func (u *UnWrapper) ShiftUint64() (uint64, error) {
	buffer := make([]byte, 8)
	if _, err := u.reader.Read(buffer); err != nil {
		return 0, err
	}
	return u.endian.Uint64(buffer), nil
}

// ShiftInt64 fetch 8 bytes in io.Reader and convert it to int64.
func (u *UnWrapper) ShiftInt64() (int64, error) {
	i, err := u.ShiftUint64()
	return int64(i), err
}

// FetchUint64 read 8 bytes, convert it to uint64 and set it to i.
func (u *UnWrapper) FetchUint64(i *uint64) *UnWrapper {
	return u.errFilter(func() {
		*i, u.err = u.ShiftUint64()
	})
}

// FetchInt64 read 8 bytes, convert it to int64 and set it to i.
func (u *UnWrapper) FetchInt64(i *int64) *UnWrapper {
	return u.FetchUint64((*uint64)(unsafe.Pointer(i)))
}

// ShiftFloat32 fetch 4 bytes in io.Reader and convert it to float32.
func (u *UnWrapper) ShiftFloat32() (float32, error) {
	buffer := make([]byte, 4)
	if _, err := u.reader.Read(buffer); err != nil {
		return 0, err
	}
	return math.Float32frombits(u.endian.Uint32(buffer)), nil
}

// ShiftFloat64 fetch 8 bytes in io.Reader and convert it to float64.
func (u *UnWrapper) ShiftFloat64() (float64, error) {
	buffer := make([]byte, 8)
	if _, err := u.reader.Read(buffer); err != nil {
		return 0, err
	}
	return math.Float64frombits(u.endian.Uint64(buffer)), nil
}

// FetchFloat32 read 4 bytes, convert it to float32 and set it to i.
func (u *UnWrapper) FetchFloat32(i *float32) *UnWrapper {
	return u.errFilter(func() {
		*i, u.err = u.ShiftFloat32()
	})
}

// FetchFloat64 read 8 bytes, convert it to float64 and set it to i.
func (u *UnWrapper) FetchFloat64(i *float64) *UnWrapper {
	return u.errFilter(func() {
		*i, u.err = u.ShiftFloat64()
	})
}

// ShiftString fetch n bytes, convert it to string. Returns string and an error.
func (u *UnWrapper) ShiftString(n uint64) (string, error) {
	buffer := make([]byte, n)
	if _, err := u.reader.Read(buffer); err != nil {
		return "", err
	}
	return string(buffer), nil
}

// FetchString read n bytes, convert it to string and set t to s.
func (u *UnWrapper) FetchString(n uint64, s *string) *UnWrapper {
	return u.errFilter(func() {
		*s, u.err = u.ShiftString(n)
	})
}

// StringWithUint16Prefix read 2 bytes as string length, then read N bytes,
// convert it to string and set it to s.
func (u *UnWrapper) StringWithUint16Prefix(s *string) *UnWrapper {
	return u.errFilter(func() {
		var n uint16
		n, u.err = u.ShiftUint16()
		u.FetchString(uint64(n), s)
	})
}

// StringWithUint32Prefix read 4 bytes as string length, then read N bytes,
// convert it to string and set it to s.
func (u *UnWrapper) StringWithUint32Prefix(s *string) *UnWrapper {
	return u.errFilter(func() {
		var n uint32
		n, u.err = u.ShiftUint32()
		u.FetchString(uint64(n), s)
	})
}

// StringWithUint64Prefix read 8 bytes as string length, then read N bytes,
// convert it to string and set it to s.
func (u *UnWrapper) StringWithUint64Prefix(s *string) *UnWrapper {
	return u.errFilter(func() {
		var n uint64
		n, u.err = u.ShiftUint64()
		u.FetchString(n, s)
	})
}

// BytesWithUint16Prefix read 2 bytes as bytes length, then read N bytes and set
// it to bytes.
func (u *UnWrapper) BytesWithUint16Prefix(bytes *[]byte) *UnWrapper {
	return u.errFilter(func() {
		var n uint16
		n, u.err = u.ShiftUint16()
		u.FetchBytes(uint64(n), bytes)
	})
}

// BytesWithUint32Prefix read 4 bytes as bytes length, then read N bytes and set
// it to bytes.
func (u *UnWrapper) BytesWithUint32Prefix(bytes *[]byte) *UnWrapper {
	return u.errFilter(func() {
		var n uint32
		n, u.err = u.ShiftUint32()
		u.FetchBytes(uint64(n), bytes)
	})
}

// BytesWithUint64Prefix read 8 bytes as bytes length, then read N bytes and set
// it to bytes.
func (u *UnWrapper) BytesWithUint64Prefix(bytes *[]byte) *UnWrapper {
	return u.errFilter(func() {
		var n uint64
		n, u.err = u.ShiftUint64()
		u.FetchBytes(n, bytes)
	})
}

func (u *UnWrapper) errFilter(f func()) *UnWrapper {
	if u.err == nil {
		f()
	}
	return u
}
