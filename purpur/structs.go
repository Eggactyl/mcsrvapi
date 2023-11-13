package purpur

type PurpurProjects struct {
	Projects []string `json:"projects"`
}

type PurpurProject struct {
	Project  string   `json:"project"`
	Versions []string `json:"versions"`
}

type PurpurProjectVersion struct {
	Project string              `json:"project"`
	Version string              `json:"version"`
	Builds  PurpurProjectBuilds `json:"builds"`
}

type PurpurProjectBuilds struct {
	Latest string   `json:"latest"`
	All    []string `json:"all"`
}

type PurpurProjectBuild struct {
	Project   string                      `json:"project"`
	Version   string                      `json:"version"`
	Build     string                      `json:"build"`
	Result    string                      `json:"result"`
	Timestamp int64                       `json:"timestamp"`
	Duration  int64                       `json:"duration"`
	Commits   []PurpurProjectBuildCommits `json:"commits"`
	Md5       string                      `json:"md5"`
}

type PurpurProjectBuildCommits struct {
	Author      string `json:"author"`
	Email       string `json:"email"`
	Description string `json:"description"`
	Hash        string `json:"hash"`
	Timestamp   int64  `json:"timestamp"`
}
