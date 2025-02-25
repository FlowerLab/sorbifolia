package password

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"

	"golang.org/x/crypto/argon2"
)

func New() Generator {
	return Argon2{
		Time:    1,
		Memory:  64 * 1024,
		Threads: 1,
		KeyLen:  16,
	}
}

type Argon2 struct {
	Time, Memory uint32
	Threads      uint8
	KeyLen       uint32
}

func (a Argon2) MustGenerate(password string) string {
	data, _ := a.Generate(password)
	return data
}

func (a Argon2) Generate(password string) (string, error) {
	data := make([]byte, 4+4+1+4+a.KeyLen)
	binary.BigEndian.PutUint32(data[0:4], a.Time)
	binary.BigEndian.PutUint32(data[4:8], a.Memory)
	data[8] = a.Threads
	binary.BigEndian.PutUint32(data[9:13], a.KeyLen)
	_, _ = rand.Read(data[13 : 13+a.KeyLen])
	data = append(data, argon2.IDKey([]byte(password), data[13:13+a.KeyLen],
		a.Time, a.Memory, a.Threads, a.KeyLen)...)

	return base64.RawStdEncoding.EncodeToString(data), nil
}

func (a Argon2) Compare(hashedPassword, password string) bool {
	data, err := base64.RawStdEncoding.DecodeString(hashedPassword)
	if err != nil || len(data) < 13 {
		return false
	}
	time := binary.BigEndian.Uint32(data[0:4])
	memory := binary.BigEndian.Uint32(data[4:8])
	threads := data[8]
	keyLen := binary.BigEndian.Uint32(data[9:13])
	if len(data) < 13+int(keyLen) {
		return false
	}

	return bytes.Equal(
		argon2.IDKey([]byte(password), data[13:13+keyLen], time, memory, threads, keyLen),
		data[13+keyLen:])
}
