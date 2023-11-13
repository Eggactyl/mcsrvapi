package paper

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
var ApiURL = fmt.Sprintf("https://api.papermc.io/v%d/", ApiVersion)

func GetProjects() (*PaperProjects, error) {

	res, err := http.Get(ApiURL + "projects")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	paperProjects := PaperProjects{}
	err = decoder.Decode(&paperProjects)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &paperProjects, nil

}

func GetProjectVersion(project string, version string) (*PaperProjectVersion, error) {

	res, err := http.Get(ApiURL + fmt.Sprintf("projects/%s/versions/%s", project, version))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	paperProjectVersion := PaperProjectVersion{}
	err = decoder.Decode(&paperProjectVersion)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &paperProjectVersion, nil

}

func GetProjectVersionBuilds(project string, version string) (*PaperProjectVersionBuilds, error) {

	res, err := http.Get(ApiURL + fmt.Sprintf("projects/%s/versions/%s/builds", project, version))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	paperProjectVersionBuilds := PaperProjectVersionBuilds{}
	err = decoder.Decode(&paperProjectVersionBuilds)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &paperProjectVersionBuilds, nil

}

func GetProjectVersionBuild(project string, version string, build int) (*PaperProjectVersionBuild, error) {

	res, err := http.Get(ApiURL + fmt.Sprintf("projects/%s/versions/%s/builds/%d", project, version, build))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	paperProjectVersionBuild := PaperProjectVersionBuild{}
	err = decoder.Decode(&paperProjectVersionBuild)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &paperProjectVersionBuild, nil

}

func GetProjectVersionGroup(project string, group string) (*PaperProjectVersionGroup, error) {

	res, err := http.Get(ApiURL + fmt.Sprintf("projects/%s/version_group/%s", project, group))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	paperProjectVersionGroup := PaperProjectVersionGroup{}
	err = decoder.Decode(&paperProjectVersionGroup)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &paperProjectVersionGroup, nil

}

func GetProjectVersionGroupBuilds(project string, group string) (*PaperProjectVersionGroupBuilds, error) {

	res, err := http.Get(ApiURL + fmt.Sprintf("projects/%s/version_group/%s/builds", project, group))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	paperProjectVersionGroupBuilds := PaperProjectVersionGroupBuilds{}
	err = decoder.Decode(&paperProjectVersionGroupBuilds)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &paperProjectVersionGroupBuilds, nil

}

func GetProjectVersionBuildDownload(project string, version string, build int, path string) (*mcsrvapi.DownloadChecksums, error) {

	buildInfo, err := GetProjectVersionBuild(project, version, build)
	if err != nil {
		return nil, err
	}

	res, err := http.Get(ApiURL + fmt.Sprintf("projects/%s/versions/%s/builds/%d/downloads/%s", project, version, build, buildInfo.Downloads.Application.Name))
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

	if hex.EncodeToString(sha256Sum) != buildInfo.Downloads.Application.Sha256 {

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
