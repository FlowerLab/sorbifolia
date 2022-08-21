package padding

import (
	"bytes"
	"crypto/rand"
	"errors"
)

var (
	ErrDataLength      = errors.New("data length is wrong")
	ErrInvalidPadding  = errors.New("invalid padding")
	ErrSizeOutOfBounds = errors.New("size is out of bounds")
)

type Padding interface {
	Pad(buf []byte, size int) ([]byte, error)
	UnPad(buf []byte, size int) ([]byte, error)
}

type PKCS7 struct{}

func (PKCS7) Pad(data []byte, size int) ([]byte, error) {
	if err := checkSize(size); err != nil {
		return nil, err
	}
	buf := padSize(len(data), size)
	return append(data, bytes.Repeat([]byte{byte(buf)}, buf)...), nil
}

func (PKCS7) UnPad(data []byte, size int) ([]byte, error) {
	dataLen := len(data)
	if err := checkUnPadSize(dataLen, size); err != nil {
		return nil, err
	}

	paddingBytes := int(data[dataLen-1])
	if paddingBytes > size || paddingBytes <= 0 {
		return nil, ErrInvalidPadding
	}
	for _, v := range data[dataLen-paddingBytes : dataLen-1] {
		if int(v) != paddingBytes {
			return nil, ErrInvalidPadding
		}
	}

	return data[:dataLen-int(data[dataLen-1])], nil
}

type NoPadding struct{}

func (NoPadding) Pad(buf []byte, size int) ([]byte, error)   { return buf, nil }
func (NoPadding) UnPad(buf []byte, size int) ([]byte, error) { return buf, nil }

type ZeroPadding struct{}

func (ZeroPadding) Pad(data []byte, size int) ([]byte, error) {
	if err := checkSize(size); err != nil {
		return nil, err
	}
	return append(data, bytes.Repeat([]byte{byte(0)}, padSize(len(data), size))...), nil
}

func (ZeroPadding) UnPad(data []byte, size int) ([]byte, error) {
	dataLen := len(data)
	if err := checkUnPadSize(dataLen, size); err != nil {
		return nil, err
	}

	paddingBytes := 0
	for data[dataLen-1-paddingBytes] == 0 {
		paddingBytes++
	}
	if paddingBytes > size || paddingBytes <= 0 {
		return nil, ErrInvalidPadding
	}
	return data[0 : dataLen-paddingBytes], nil
}

type ISO10126 struct{}

func (ISO10126) Pad(data []byte, size int) ([]byte, error) {
	if err := checkSize(size); err != nil {
		return nil, err
	}

	var (
		paddingBytes = padSize(len(data), size)
		buf          = make([]byte, paddingBytes-1)
	)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}
	return append(data, append(buf, byte(paddingBytes))...), nil
}

func (ISO10126) UnPad(data []byte, size int) ([]byte, error) {
	dataLen := len(data)
	if err := checkUnPadSize(dataLen, size); err != nil {
		return nil, err
	}

	i := int(data[dataLen-1])
	if i > size || i <= 0 {
		return nil, ErrInvalidPadding
	}
	return data[0 : dataLen-i], nil
}

type ANSIx923 struct{}

func (ANSIx923) Pad(data []byte, size int) ([]byte, error) {
	if err := checkSize(size); err != nil {
		return nil, err
	}
	paddingBytes := padSize(len(data), size)

	return append(data, append(
		bytes.Repeat([]byte{byte(0)}, paddingBytes-1), byte(paddingBytes))...), nil
}

func (ANSIx923) UnPad(data []byte, size int) ([]byte, error) {
	dataLen := len(data)
	if err := checkUnPadSize(dataLen, size); err != nil {
		return nil, err
	}

	i := int(data[dataLen-1])
	if i > size || i <= 0 {
		return nil, ErrInvalidPadding
	}

	if dataLen-i < dataLen-2 {
		for _, v := range data[dataLen-i : dataLen-2] {
			if int(v) != 0 {
				return nil, errors.New("invalid padding found")
			}
		}
	}
	return data[0 : dataLen-i], nil
}

func padSize(dataSize, blockSize int) int {
	return blockSize - (dataSize % blockSize)
}

func checkSize(size int) error {
	if size < 1 || size >= 256 {
		return ErrSizeOutOfBounds
	}
	return nil
}

func checkUnPadSize(dataLen, size int) error {
	if dataLen%size != 0 {
		return ErrDataLength
	}
	return nil
}
