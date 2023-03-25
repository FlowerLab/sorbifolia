package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"io"
	"os"

	"go.x2ox.com/sorbifolia/cryptutils"
)

func encrypt(key, src, dst string) error {
	cs, err := parseCry(key)
	if err != nil {
		return err
	}

	var (
		sf *os.File
		df *os.File
	)
	if sf, err = os.Open(src); err != nil {
		return err
	}
	if df, err = os.Create(dst); err != nil {
		_ = sf.Close()
		return err
	}

	var (
		bf  = bufio.NewWriter(df)
		buf = make([]byte, 1024*128)
		i   = 0
	)

	for {
		if i, err = sf.Read(buf); err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		cs.Encrypt(buf[:i], buf[:i])
		if _, err = bf.Write(buf[:i]); err != nil {
			break
		}
	}

	_ = sf.Close()

	if err != nil {
		_ = sf.Close()
		return err
	}

	if err = bf.Flush(); err != nil {
		_ = df.Close()
		return err
	}
	return df.Close()
}

func decrypt(key, src, dst string) error {
	cs, err := parseCry(key)
	if err != nil {
		return err
	}

	var (
		df *os.File
		sf *os.File
	)
	if df, err = os.Create(dst); err != nil {
		return err
	}
	if sf, err = os.Open(src); err != nil {
		_ = df.Close()
		return err
	}

	var (
		bf  = bufio.NewWriter(df)
		buf = make([]byte, 1024*128)
		i   = 0
	)

	for {
		if i, err = sf.Read(buf); err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		cs.Decrypt(buf[:i], buf[:i])
		if _, err = bf.Write(buf[:i]); err != nil {
			break
		}
	}

	_ = sf.Close()

	if err != nil {
		_ = df.Close()
		return err
	}

	if err = bf.Flush(); err != nil {
		_ = df.Close()
		return err
	}
	return df.Close()
}

func parseCry(key string) (cryptutils.CryptStream, error) {
	data, err := hex.DecodeString(key[:16])
	if err != nil {
		return nil, err
	}

	var cb cipher.Block
	if cb, err = aes.NewCipher(data[:16]); err != nil {
		return nil, err
	}

	return cryptutils.CTR(cb, data[16:]), nil
}
