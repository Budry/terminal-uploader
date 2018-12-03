package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func uploadFile(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func main() {

	if len(os.Args) < 4 {
		panic("Invalid arguments")
	}

	path, _ := os.Getwd()
	path += "/" + os.Args[3]
	uri := os.Args[1]
	fieldName := os.Args[2]

	extraParams := map[string]string{}

	request, err := uploadFile(uri, extraParams, fieldName, path)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(body)
	}
}
