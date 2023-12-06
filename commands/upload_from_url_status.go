package commands

type UploadFromUrlStatus struct {
	FileId          string
	MtaId           string
	ClientActions   []string
	ExecutionStatus ExecutionStatus
}
