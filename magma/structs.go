package magma

import "time"

type MagmaRelease struct {
	Name          string    `json:"name"`
	TagName       string    `json:"tag_name"`
	CreatedAt     time.Time `json:"created_at"`
	Link          string    `json:"link"`
	InstallerLink string    `json:"installer_link"`
	GitCommitURL  string    `json:"git_commit_url"`
	Archived      bool      `json:"archived"`
}
