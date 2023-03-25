package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getKeyServerURL(s string) string {
	return fmt.Sprintf("https://%s/api/%s", os.Getenv("KEY_SERVER"), s)
}

func RegisterKey(hash string) ([]byte, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(meta); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, getKeyServerURL(hash), &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Sbfa %s", os.Getenv("KEY_SERVER_KEY")))

	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return nil, err
	}

	var bts []byte
	bts, err = io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	return bts, err
}

func GetKey(id string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, getKeyServerURL(id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Sbfa %s", os.Getenv("KEY_SERVER_KEY")))
	req.Header.Set("X-SB-FA-Name", os.Getenv("_SB_FA_NAME"))
	req.Header.Set("X-SB-FA-Tag", os.Getenv("_SB_FA_TAG"))

	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return nil, err
	}

	var bts []byte
	bts, err = io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return bts, err
}
