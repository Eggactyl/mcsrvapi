package purpur

import (
	"fmt"
	"io"
	"net/http"

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

	req, err := http.NewRequest("GET", ApiURL+fmt.Sprintf("%s/%s", project, version), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cache-Control", "no-cache")

	res, err := http.DefaultClient.Do(req)
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

func (p *PurpurProjectVersion) GetLatestBuild() (*PurpurProjectBuild, error) {
	return GetProjectVersionBuild(p.Project, p.Version, p.Builds.Latest)
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

func (p *PurpurProjectBuild) DownloadURL() string {
	return fmt.Sprintf(ApiURL+"%s/%s/%s/download", p.Project, p.Version, p.Build)
}

func GetProjectVersionBuildDownload(project string, version string, build string, path string) (*mcsrvapi.DownloadChecksums, error) {

	buildInfo, err := GetProjectVersionBuild(project, version, build)
	if err != nil {
		return nil, err
	}

	checksums, err := mcsrvapi.DownloadWithChecksums(ApiURL+fmt.Sprintf("%s/%s/%s/download", project, version, build), path, mcsrvapi.ChecksumTypeMd5, buildInfo.Md5)
	if err != nil {
		return nil, err
	}

	return checksums, nil

}
