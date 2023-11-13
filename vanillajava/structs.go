package vanillajava

import "time"

type VanillaJavaManifest struct {
	Latest   LatestVanillaJavaManifestVersion `json:"latest"`
	Versions []VanillaJavaManifestVersion     `json:"versions"`
}

type LatestVanillaJavaManifestVersion struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

type VanillaJavaManifestVersion struct {
	Id              string    `json:"id"`
	Type            string    `json:"type"`
	Url             string    `json:"url"`
	Time            time.Time `json:"time"`
	ReleaseTime     time.Time `json:"releaseTime"`
	Sha1            string    `json:"sha1"`
	ComplianceLevel int       `json:"complianceLevel"`
}

type VanillaJavaVersion struct {
	Downloads VanillaJavaVersionDownloads
}

type VanillaJavaVersionDownloads struct {
	Client         VanillaJavaVersionDownload `json:"client"`
	ClientMappings VanillaJavaVersionDownload `json:"client_mappings"`
	Server         VanillaJavaVersionDownload `json:"server"`
	ServerMappings VanillaJavaVersionDownload `json:"server_mappings"`
}

type VanillaJavaVersionDownload struct {
	Sha1 string `json:"sha1"`
	Size int64  `json:"size"`
	URL  string `json:"url"`
}
