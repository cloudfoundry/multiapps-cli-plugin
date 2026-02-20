package commands

type UploadFromUrlStatus struct {
	FileId          string
	MtaId           string
	SchemaVersion   string
	ClientActions   []string
	ExecutionStatus ExecutionStatus
}
