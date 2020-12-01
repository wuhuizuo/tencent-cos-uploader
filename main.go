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

// uploadCfg local vs remote link
type uploadCfg struct {
	Local  string
	Remote string
}

type runCfg struct {
	Bucket *bucketCfg
	File   uploadCfg
}

func main() {
	cfg := parseCliArgs()

	contentType, err := guessFileContentType(cfg.File.Local)
	if err != nil {
		log.Fatalln("guess file content type failed:", err)
	}
	log.Println("guess file content type: ", contentType)

	putOptions := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
	}

	client := newCosClient(cfg.Bucket)
	resp, err := client.Object.PutFromFile(
		context.TODO(),
		cfg.File.Remote,
		cfg.File.Local,
		putOptions)
	if err != nil {
		log.Fatalln(err)
	}

	etag := resp.Header.Get("ETag")

	log.Println("file md5：", etag)
}

func parseCliArgs() *runCfg {
	var bucketArgs bucketCfg
	var upload uploadCfg
	flag.StringVar(&bucketArgs.BucketName, "bucketName", "", "[required] cos bucket name")
	flag.StringVar(&bucketArgs.BucketRegion, "bucketRegion", "", "[required] cos bucket region, like ap-shanghai,ap-beijing")
	flag.StringVar(&bucketArgs.Auth.SecretID, "secretID", "", "[required] cos client secret id")
	flag.StringVar(&bucketArgs.Auth.SecretKey, "secretKey", "", "[required] cos client secret key")
	flag.StringVar(&bucketArgs.Auth.SessionToken, "sessionToken", "", "[optional] cos client temporary session token")
	flag.StringVar(&upload.Local, "upload-file", "", "[require] file path to upload")
	flag.StringVar(&upload.Remote, "remote-path", "", "[require] remote bucket object path to store(contained filename part)")

	showVersion := flag.Bool("version", false, "prints current version")
	flag.Usage = usage

	flag.Parse()

	if showVersion != nil && *showVersion {
		printVersion()
		os.Exit(0)
	}

	if bucketArgs.BucketName == "" ||
		bucketArgs.BucketRegion == "" ||
		bucketArgs.Auth.SecretID == "" ||
		bucketArgs.Auth.SecretKey == "" ||
		upload.Local == "" ||
		upload.Remote == "" {
		usage()
		os.Exit(1)
	}

	return &runCfg{Bucket: &bucketArgs, File: upload}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options], options list:\n", os.Args[0])
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Fprintln(os.Stdout, "Version:\t", version)
	fmt.Fprintln(os.Stdout, "Build date:\t", buildDate)
}

func init() {
	log.SetOutput(os.Stdout)
}
