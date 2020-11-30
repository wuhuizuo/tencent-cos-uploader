package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// release version info
var (
	version   string
	buildDate string
)

// bucketCfg bucket arguments
type bucketCfg struct {
	BucketName   string
	BucketRegion string

	Auth cos.AuthorizationTransport
}

type fileList []string

func (l *fileList) String() string {
	return ""
}

func (l *fileList) Set(value string) error {
	*l = append(*l, value)

	return nil
}

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	var bucketArgs bucketCfg
	var localFile, remotePath string
	flag.StringVar(&bucketArgs.BucketName, "bucketName", "", "[required] cos bucket name")
	flag.StringVar(&bucketArgs.BucketRegion, "bucketRegion", "", "[required] cos bucket region")
	flag.StringVar(&bucketArgs.Auth.SecretID, "secretID", "", "[required] cos client secret id")
	flag.StringVar(&bucketArgs.Auth.SecretKey, "secretKey", "", "[required] cos client secret key")
	flag.StringVar(&bucketArgs.Auth.SessionToken, "sessionToken", "", "[optional] cos client temporary session token")
	flag.StringVar(&localFile, "upload-file", "", "[require] file path to upload")
	flag.StringVar(&remotePath, "remote-path", "", "[require] remote bucket object path to store(contained filename part)")

	showVersion := flag.Bool("version", false, "prints current version")
	flag.Usage = usage

	flag.Parse()

	if showVersion != nil {
		printVersion()
		os.Exit(0)
	}

	contentType, err := guessFileContentType(localFile)
	if err != nil {
		log.Fatalln("guess file content failed:", err)
	}
	putOptions := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
	}

	client := newCosClient(bucketArgs.BucketName, bucketArgs.BucketRegion, &bucketArgs.Auth)
	resp, err := client.Object.PutFromFile(context.TODO(), remotePath, localFile, putOptions)
	if err != nil {
		log.Fatalln(err)
	}
	etag := resp.Header.Get("ETag")

	log.Println("file md5：", etag)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options], options list:\n", os.Args[0])
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Fprintln(os.Stdout, "Version:\t", version)
	fmt.Fprintln(os.Stdout, "Build date:\t", buildDate)
}
