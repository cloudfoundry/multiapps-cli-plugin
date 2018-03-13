package models

type SharedDomain struct {
	Name string
	Guid string
	Url  string
}

func NewSharedDomain(name, guid, url string) SharedDomain {
	return SharedDomain{Name: name, Guid: guid, Url: url}
}
