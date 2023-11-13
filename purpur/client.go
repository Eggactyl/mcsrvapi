package purpur

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/eggactyl/mcsrvapi"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

var ApiVersion = 2
var ApiURL = fmt.Sprintf("https://api.purpurmc.org/v%d/", ApiVersion)

func GetProjects() (*PurpurProjects, error) {

	res, err := http.Get(ApiURL)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	purpurProjects := PurpurProjects{}
	err = decoder.Decode(&purpurProjects)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &purpurProjects, nil

}

func GetProject(project string) (*PurpurProject, error) {

	res, err := http.Get(ApiURL + project)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	purpurProject := PurpurProject{}
	err = decoder.Decode(&purpurProject)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &purpurProject, nil

}

func GetProjectVersion(project string, version string) (*PurpurProjectVersion, error) {

	res, err := http.Get(ApiURL + fmt.Sprintf("%s/%s", project, version))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	purpurProjectVersion := PurpurProjectVersion{}
	err = decoder.Decode(&purpurProjectVersion)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &purpurProjectVersion, nil

}

func GetProjectVersionBuild(project string, version string, build string) (*PurpurProjectBuild, error) {

	res, err := http.Get(ApiURL + fmt.Sprintf("%s/%s/%s", project, version, build))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	purpurProjectVersionBuild := PurpurProjectBuild{}
	err = decoder.Decode(&purpurProjectVersionBuild)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &purpurProjectVersionBuild, nil

}

func GetProjectVersionBuildDownload(project string, version string, build string, path string) (*mcsrvapi.DownloadChecksums, error) {

	buildInfo, err := GetProjectVersionBuild(project, version, build)
	if err != nil {
		return nil, err
	}

	res, err := http.Get(ApiURL + fmt.Sprintf("%s/%s/%s/download", project, version, build))
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

	file.Sync()

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

	if hex.EncodeToString(md5Sum) != buildInfo.Md5 {

		file.Close()

		if err := os.Remove(path); err != nil {
			return nil, err
		}

		return nil, errors.New("mismatched checksums")

	}

	checksums := mcsrvapi.DownloadChecksums{
		Sha256: hex.EncodeToString(sha256Sum),
		Sha1:   hex.EncodeToString(sha1Sum),
		Md5:    hex.EncodeToString(md5Sum),
	}

	return &checksums, nil

}
