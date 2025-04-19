package domain

type SystemInfo struct {
	ID          string `json:"id" example:"uid123"`
	Hostname    string `json:"hostname" example:"jarvis"`
	Platform    string `json:"platform" example:"Linux"`
	Arch        string `json:"arch" example:"x64"`
	CpuCount    int    `json:"cpu_count" example:"16"`
	NodeVersion string `json:"node_version" example:"22.1.2"`
}
