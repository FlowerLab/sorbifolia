package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
)

const (
	defaultOutputFile   = "/_cf_file"
	defaultCompressFile = "/_cf_c_file"
	defaultOutputMeta   = "/_cf_meta"
)

// cryptutils-file e /data/server

func main() {
	if len(os.Args) < 3 {
		return
	}
	method := os.Args[1]
	filename := os.Args[2]

	switch method {
	case "d":
		bts, err := os.ReadFile(defaultOutputMeta)
		if err != nil {
			log.Fatalln(err)
		}

		var key string
		if key, err = GetKey(string(bts)); err != nil {
			log.Fatalln(err)
		}

		if err = decrypt(key, defaultOutputFile, defaultCompressFile); err != nil {
			log.Fatalln(err)
		}

		if err = UnCompress(defaultCompressFile, filename); err != nil {
			log.Fatalln(err)
		}
	case "e":
		if err := Compress(filename, defaultCompressFile); err != nil {
			log.Fatalln(err)
		}

		hash, err := getFileHash(defaultCompressFile)
		if err != nil {
			log.Fatalln(err)
		}

		var m []byte
		if m, err = RegisterKey(hash); err != nil {
			log.Fatalln(err)
		}

		if err = os.WriteFile(defaultOutputMeta, m, 0666); err != nil {
			log.Fatalln(err)
		}

		var key string
		if key, err = GetKey(string(m)); err != nil {
			log.Fatalln(err)
		}

		if err = encrypt(key, defaultCompressFile, defaultOutputFile); err != nil {
			log.Fatalln(err)
		}
		if err = os.RemoveAll(defaultCompressFile); err != nil {
			log.Fatalln(err)
		}
	}
}

func getFileHash(filename string) (string, error) {
	var (
		h         = sha256.New()
		file, err = os.Open(filename)
	)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(h, file)
	if _ = file.Close(); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
