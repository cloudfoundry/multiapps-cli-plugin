package models

type CloudFoundryApplication struct {
	Name      string `json:"name"`
	Guid      string `json:"guid"`
	State     string `json:"state"`
	SpaceGuid string `jsonry:"relationships.space.data.guid"`

	MtaNamespace string `jsonry:"metadata.annotations.mta_namespace"`
}

type AppProcessStatisticsResponse struct {
	Resources []ApplicationProcessStatistics `json:"resources"`
}

type ApplicationProcessStatistics struct {
	State  string `json:"state"`
	Memory int64  `jsonry:"usage.mem"`
	Disk   int64  `jsonry:"usage.disk"`
}

type ApplicationRoute struct {
	Host string `json:"host,omitempty"`
	Path string `json:"path,omitempty"`
	Url  string `json:"url"`
}
