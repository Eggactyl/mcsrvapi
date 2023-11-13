package mcsrvapi

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
)

var ErrInvalidChecksumType = errors.New("invalid checksum type")
var ErrMismatchedChecksum = errors.New("mismatched checksums")

type DownloadChecksums struct {
	Sha1   string
	Sha256 string
	Md5    string
}

type ChecksumType = string

const (
	ChecksumTypeSha1   = "sha1"
	ChecksumTypeSha256 = "sha256"
	ChecksumTypeMd5    = "md5"
)

func DownloadWithChecksums(url string, path string, checksumType string, checksum string) (*DownloadChecksums, error) {

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil && err != io.EOF {
		return nil, err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	sha256Hash := sha256.New()
	sha1Hash := sha1.New()
	md5Hash := md5.New()

	buffer := make([]byte, 4096)

	for {

		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if n == 0 {
			break
		}

		if _, err := sha256Hash.Write(buffer[:n]); err != nil {
			return nil, err
		}

		if _, err := sha1Hash.Write(buffer[:n]); err != nil {
			return nil, err
		}

		if _, err := md5Hash.Write(buffer[:n]); err != nil {
			return nil, err
		}

	}

	sha256Sum := sha256Hash.Sum(nil)
	sha1Sum := sha1Hash.Sum(nil)
	md5Sum := md5Hash.Sum(nil)

	var actualChecksum string

	switch checksumType {
	case "sha1":
		actualChecksum = hex.EncodeToString(sha1Sum)
	case "sha256":
		actualChecksum = hex.EncodeToString(sha256Sum)
	case "md5":
		actualChecksum = hex.EncodeToString(md5Sum)
	default:
		return nil, ErrInvalidChecksumType
	}

	if actualChecksum != checksum {

		file.Close()

		if err := os.Remove(path); err != nil {
			return nil, err
		}

		return nil, ErrMismatchedChecksum

	}

	checksums := DownloadChecksums{
		Sha256: hex.EncodeToString(sha256Sum),
		Sha1:   hex.EncodeToString(sha1Sum),
		Md5:    hex.EncodeToString(md5Sum),
	}

	return &checksums, nil

}
