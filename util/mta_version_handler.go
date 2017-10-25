package util

import "github.com/SAP/cf-mta-plugin/clients/models"

const unknownMtaVersion string = "0.0.0-unknown"

// GetMtaVersionAsString returns an MTA's version as a string or "?" if the version is unknown.
func GetMtaVersionAsString(mta *models.Mta) string {
	return getDefaultIfUnknown(mta.Metadata.Version)
}

func getDefaultIfUnknown(version string) string {
	if version != unknownMtaVersion {
		return version
	}
	return "?"
}
