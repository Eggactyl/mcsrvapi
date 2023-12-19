package vanillajava

import (
	"errors"
	"io"
	"net/http"

	"github.com/eggactyl/mcsrvapi"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

var ErrNotFound = errors.New("mc version not found")
var ApiURL = "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"

func GetManifest() (*VanillaJavaManifest, error) {

	res, err := http.Get(ApiURL)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	manifest := VanillaJavaManifest{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&manifest)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &manifest, nil

}

func (m *VanillaJavaManifest) GetManifestVersion(version string) (*VanillaJavaManifestVersion, error) {

	var chosenVersion *VanillaJavaManifestVersion

	for _, mVersion := range m.Versions {

		if mVersion.Id == version {

			chosenVersion = &mVersion
			break

		}

	}

	if chosenVersion == nil {
		return nil, ErrNotFound
	}

	return chosenVersion, nil

}

func (m *VanillaJavaManifest) GetVersion(version string) (*VanillaJavaVersion, error) {

	manifestVer, err := m.GetManifestVersion(version)
	if err != nil {
		return nil, err
	}

	res, err := http.Get(manifestVer.Url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	vanillaJavaVersion := VanillaJavaVersion{}
	err = decoder.Decode(&vanillaJavaVersion)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &vanillaJavaVersion, nil

}

func (m *VanillaJavaManifest) ClientDownload(version string, path string) (*mcsrvapi.DownloadChecksums, error) {

	buildInfo, err := m.GetVersion(version)
	if err != nil {
		return nil, err
	}

	checksums, err := mcsrvapi.DownloadWithChecksums(buildInfo.Downloads.Client.URL, path, mcsrvapi.ChecksumTypeSha1, buildInfo.Downloads.Client.Sha1)
	if err != nil {
		return nil, err
	}

	return checksums, nil

}

func (m *VanillaJavaManifest) ServerDownload(version string, path string) (*mcsrvapi.DownloadChecksums, error) {

	buildInfo, err := m.GetVersion(version)
	if err != nil {
		return nil, err
	}

	checksums, err := mcsrvapi.DownloadWithChecksums(buildInfo.Downloads.Server.URL, path, mcsrvapi.ChecksumTypeSha1, buildInfo.Downloads.Server.Sha1)
	if err != nil {
		return nil, err
	}

	return checksums, nil

}
