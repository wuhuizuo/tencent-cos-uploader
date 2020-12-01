package main

import (
	"net/http"
	"os"

	cos "github.com/tencentyun/cos-go-sdk-v5"
)

// bucketCfg bucket arguments
type bucketCfg struct {
	BucketName   string
	BucketRegion string

	Auth cos.AuthorizationTransport
}

func newCosClient(cfg *bucketCfg) *cos.Client {
	bucketURL := cos.NewBucketURL(cfg.BucketName, cfg.BucketRegion, true)
	httpClient := &http.Client{Transport: &cfg.Auth}

	return cos.NewClient(&cos.BaseURL{BucketURL: bucketURL}, httpClient)
}

func guessFileContentType(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
