package models

type Domain struct {
	Name string
	Guid string
}

func NewDomain(name, guid string) Domain {
	return Domain{Name: name, Guid: guid}
}
