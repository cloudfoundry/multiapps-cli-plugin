package restclient

// RestClientOperations is an interface having all RestClient operations
type RestClientOperations interface {
	PurgeConfiguration(org, space string) error
}
