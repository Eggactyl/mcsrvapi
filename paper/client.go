package paper

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/eggactyl/mcsrvapi"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

var ApiVersion = 2
var ApiURL = fmt.Sprintf("https://api.papermc.io/v%d/", ApiVersion)

var ErrNoStableVer = errors.New("no stable version found")

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

func (p *PaperProject) GetVersionBuilds(version string) (*PaperProjectVersionBuilds, error) {
	return GetProjectVersionBuilds(p.ProjectID, version)
}

func (p *PaperProject) GetVersion(version string) (*PaperProjectVersion, error) {
	return GetProjectVersion(p.ProjectID, version)
}

func (v PaperProjectVersionBuilds) LatestBuild() *PaperProjectVersionBuildBase {
	return &v.Builds[len(v.Builds)-1]
}

func (v PaperProjectVersionBuilds) LatestStableBuild() (*PaperProjectVersionBuildBase, error) {

	latestBuildNum := 0
	latestBuildIndex := -1

	for i, build := range v.Builds {

		if build.Channel == "default" {

			if latestBuildNum < build.Build {
				latestBuildIndex = i
				latestBuildNum = build.Build
			}

		}

	}

	if latestBuildIndex > len(v.Builds)-1 || latestBuildIndex < 0 {
		return nil, ErrNoStableVer
	}

	return &v.Builds[latestBuildIndex], nil

}

func (p *PaperProject) DownloadURL(version string, stable bool) (*string, error) {

	verInfo, err := p.GetVersionBuilds(version)
	if err != nil {
		return nil, err
	}

	var latestBuild *PaperProjectVersionBuildBase

	if !stable {
		latestBuild = verInfo.LatestBuild()
	} else {
		latestBuild, err = verInfo.LatestStableBuild()
		if err != nil {
			return nil, err
		}
	}

	downloadURL := fmt.Sprintf(ApiURL+"projects/%s/versions/%s/builds/%d/downloads/%s", p.ProjectID, version, latestBuild.Build, latestBuild.Downloads.Application.Name)

	return &downloadURL, nil

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
