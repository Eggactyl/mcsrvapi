package paper

import (
	"time"
)

type PaperProjects struct {
	Projects []string `json:"projects"`
}

type PaperProject struct {
	ProjectID     string   `json:"project_id"`
	ProjectName   string   `json:"project_name"`
	VersionGroups []string `json:"version_groups"`
	Versions      []string `json:"versions"`
}

type PaperProjectVersion struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Version     string `json:"version"`
	Builds      []int  `json:"builds"`
}

type PaperProjectVersionBuilds struct {
	ProjectID   string                         `json:"project_id"`
	ProjectName string                         `json:"project_name"`
	Version     string                         `json:"string"`
	Builds      []PaperProjectVersionBuildBase `json:"builds"`
}

type PaperProjectVersionBuildBase struct {
	Build     int                   `json:"build"`
	Time      time.Time             `json:"time"`
	Channel   string                `json:"channel"`
	Promoted  bool                  `json:"promoted"`
	Changes   []PaperProjectChanges `json:"changes"`
	Downloads PaperProjectDownloads `json:"downloads"`
}

type PaperProjectVersionBuild struct {
	ProjectID   string                `json:"project_id"`
	ProjectName string                `json:"project_name"`
	Build       int                   `json:"build"`
	Time        time.Time             `json:"time"`
	Channel     string                `json:"string"`
	Promoted    bool                  `json:"promoted"`
	Changes     []PaperProjectChanges `json:"changes"`
	Downloads   PaperProjectDownloads `json:"downloads"`
}

type PaperProjectChanges struct {
	Commit  string `json:"commit"`
	Summary string `json:"summary"`
	Message string `json:"message"`
}

type PaperProjectDownloads struct {
	Application    PaperProjectDownload `json:"application"`
	MojangMappings PaperProjectDownload `json:"mojang-mappings"`
}

type PaperProjectDownload struct {
	Name   string `json:"name"`
	Sha256 string `json:"sha256"`
}

type PaperProjectVersionGroup struct {
	ProjectID    string   `json:"project_id"`
	ProjectName  string   `json:"project_name"`
	VersionGroup string   `json:"version_group"`
	Versions     []string `json:"versions"`
}

type PaperProjectVersionGroupBuilds struct {
	ProjectID    string                         `json:"project_id"`
	ProjectName  string                         `json:"project_name"`
	VersionGroup string                         `json:"version_group"`
	Versions     []string                       `json:"versions"`
	Builds       []PaperProjectVersionBuildBase `json:"builds"`
}
