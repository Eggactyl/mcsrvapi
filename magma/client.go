package magma

import (
	"fmt"
	"io"
	"net/http"

	"github.com/eggactyl/mcsrvapi"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

var ApiVersion = 2
var ApiURL = fmt.Sprintf("https://api.magmafoundation.org/api/v%d/", ApiVersion)

func GetVersions() ([]string, error) {

	res, err := http.Get(ApiURL + "allVersions")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	magmaVersions := []string{}
	err = decoder.Decode(&magmaVersions)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return magmaVersions, nil

}

func GetVersionReleases(version string) (*[]MagmaRelease, error) {

	res, err := http.Get(ApiURL + version)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	versionReleases := []MagmaRelease{}
	err = decoder.Decode(&versionReleases)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &versionReleases, nil

}

func GetVersionDownload(version string, tag string, path string) (*mcsrvapi.DownloadChecksums, error) {

	checksums, err := mcsrvapi.DownloadWithChecksums(ApiURL+fmt.Sprintf("%s/latest/%s/download", version, tag), path, mcsrvapi.ChecksumTypeNone, "")
	if err != nil {
		return nil, err
	}

	return checksums, nil

}
