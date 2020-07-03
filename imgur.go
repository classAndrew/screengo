package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func mainf() {

}

// EncodeB64 encodes strings in to b64
func EncodeB64(filename string) string {
	f, _ := os.Open(filename)
	r, _ := ioutil.ReadAll(f)
	buf := make([]byte, int(((4*len(r)/3)+3)&-4))
	base64.StdEncoding.Encode(buf, r)
	f.Close()
	return string(buf)
}

// UploadImgur uploads files to imgur
func UploadImgur(b64 string) string {
	key := "3077b71d954c094"
	url := "https://api.imgur.com/3/image"
	method := "POST"
	data := bytes.NewBuffer(make([]byte, 0))
	writer := multipart.NewWriter(data)
	writer.WriteField("image", b64)
	writer.Close()
	req, _ := http.NewRequest(method, url, data)
	req.Header.Add("Authorization", "Client-ID "+key)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return string(body)
}

// FindID finds ID
func FindID(res string) string {
	begin := strings.Index(res, "id")
	return res[begin+5 : begin+12]
}

// FindDelHash finds delete hash
func FindDelHash(res string) string {
	begin := strings.Index(res, "deletehash")
	return res[begin+13 : begin+28]
}
