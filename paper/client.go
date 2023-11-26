package paper

import (
	"fmt"
	"io"
	"net/http"

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

func GetProject(project string) (*PaperProject, error) {

	res, err := http.Get(ApiURL + fmt.Sprintf("projects/%s", project))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	paperProject := PaperProject{}
	err = decoder.Decode(&paperProject)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &paperProject, nil

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

	checksums, err := mcsrvapi.DownloadWithChecksums(ApiURL+fmt.Sprintf("projects/%s/versions/%s/builds/%d/downloads/%s", project, version, build, buildInfo.Downloads.Application.Name), path, mcsrvapi.ChecksumTypeSha256, buildInfo.Downloads.Application.Sha256)
	if err != nil {
		return nil, err
	}

	return checksums, nil

}
