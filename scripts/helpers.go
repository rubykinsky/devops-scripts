package devops_scripts

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3FileUpload uploads a file to an S3 bucket
func S3FileUpload(bucket, key, filePath string) error {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-2")}, nil)
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return err
	}
	fileHash := hex.EncodeToString(hash.Sum(nil))

	data := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key + "-" + fileHash),
		Body:   file,
	}

	_, err = s3.New(sess).PutObject(data)
	if err != nil {
		return err
	}

	return nil
}

// GetFileSize returns the size of a file in bytes
func GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

// GetFileExtension returns the extension of a file
func GetFileExtension(filePath string) string {
	ext := filepath.Ext(filePath)
	return strings.ToLower(ext)
}

// GetFileSHA256 returns the SHA256 of a file
func GetFileSHA256(filePath string) (string, error) {
	hash := sha256.New()
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetPrettyFilePath returns a pretty file path
func GetPrettyFilePath(filePath string) string {
	return strings.Replace(filePath, "\\", "/", -1)
}

// GetRelativeFilePath returns the relative file path
func GetRelativeFilePath(filePath, baseDir string) string {
	return filepath.Rel(baseDir, filePath)
}