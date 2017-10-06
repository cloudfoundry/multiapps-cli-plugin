package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry/cli/cf/terminal"

	slmpclient "github.com/SAP/cf-mta-plugin/clients/slmpclient"
	"github.com/SAP/cf-mta-plugin/log"
	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/SAP/cf-mta-plugin/util"

	"github.com/SAP/cf-mta-plugin/clients/models"
	"golang.org/x/sync/errgroup"
)

const (
	ServiceVersion1_0 = "1.0"
	ServiceVersion1_1 = "1.1"
)

//FileUploader uploads files for the service with the specified service ID
type FileUploader struct {
	serviceID  string
	files      []string
	slmpClient slmpclient.SlmpClientOperations
}

//NewFileUploader creates a new file uploader for the specified service ID, files, and SLMP client
func NewFileUploader(serviceID string, files []string, slmpClient slmpclient.SlmpClientOperations) *FileUploader {
	return &FileUploader{
		serviceID:  serviceID,
		files:      files,
		slmpClient: slmpClient,
	}
}

//UploadFiles uploads the files
func (f *FileUploader) UploadFiles() ([]*models.File, ExecutionStatus) {
	log.Tracef("Uploading files '%v'\n", f.files)

	// Get all files that are already uploaded
	serviceFiles, err := f.slmpClient.GetServiceFiles(f.serviceID)
	if err != nil {
		ui.Failed("Could not get files for service %s: %s", terminal.EntityNameColor(f.serviceID), err)
		return nil, Failure
	}

	// Determine which files to uplaod
	filesToUpload := []os.File{}
	alreadyUploadedFiles := []*models.File{}
	for _, file := range f.files {
		// Check if the file exists
		fileInfo, err := os.Stat(file)
		if os.IsNotExist(err) {
			ui.Failed("Could not find file %s", terminal.EntityNameColor(file))
			return nil, Failure
		} else if err != nil {
			ui.Failed("Could not get information for file %s", terminal.EntityNameColor(file))
			return nil, Failure
		}

		// Check if the files is already uploaded
		if !isFileAlreadyUploaded(file, fileInfo, serviceFiles, &alreadyUploadedFiles) {
			// If not, add it to the list of uploaded files
			fileToUpload, err := os.Open(file)
			defer fileToUpload.Close()
			if err != nil {
				ui.Failed("Could not open file %s", terminal.EntityNameColor(file))
				return nil, Failure
			}
			filesToUpload = append(filesToUpload, *fileToUpload)
		}
	}

	// If there are new files to upload, upload them
	uploadedFiles := []*models.File{}
	uploadedFiles = append(uploadedFiles, alreadyUploadedFiles...)
	if len(filesToUpload) != 0 {
		ui.Say("Uploading %d files...", len(filesToUpload))

		// Iterate over all files to be uploaded
		for _, fileToUpload := range filesToUpload {
			// Print the full path of the file
			fullPath, err := filepath.Abs(fileToUpload.Name())
			if err != nil {
				ui.Failed("Could not get absolute path of file %s", terminal.EntityNameColor(fileToUpload.Name()))
				return nil, Failure
			}
			ui.Say("  " + fullPath)

			// Recreate the session if it is expired
			EnsureSlmpSession(f.slmpClient)

			// Upload the file
			shouldUploadInChunks, err := f.shouldUploadInChunks()
			if err != nil {
				ui.Failed("Could not get versions for service %s", terminal.EntityNameColor(f.serviceID), err)
				return nil, Failure
			}

			uploaded, err := f.upload(shouldUploadInChunks, fileToUpload, fullPath)
			if err != nil {
				ui.Failed("Could not upload file %s", terminal.EntityNameColor(fileToUpload.Name()))
				return nil, Failure
			}
			uploadedFiles = append(uploadedFiles, uploaded...)
		}
		ui.Ok()
	}
	return uploadedFiles, Success
}

func (f *FileUploader) upload(shouldUploadInChunks bool, fileToUpload os.File, filePath string) ([]*models.File, error) {
	if shouldUploadInChunks {
		// upload files in chunks
		return uploadInChunks(filePath, fileToUpload, f.serviceID, f.slmpClient)
	}

	//upload normally
	file, err := uploadFile(&fileToUpload, fileToUpload.Name(), f.serviceID, f.slmpClient)
	return []*models.File{file}, err
}

func (f *FileUploader) shouldUploadInChunks() (bool, error) {
	serviceVersions, err := f.slmpClient.GetServiceVersions(f.serviceID)
	if err != nil {
		return false, err

	}
	baseServiceVersion := getBaseServiceVersion(serviceVersions.ComponentVersions)
	return baseServiceVersion == ServiceVersion1_1, nil
}

func getBaseServiceVersion(serviceVersions []*models.ComponentVersion) string {
	baseServiceVersion := ServiceVersion1_0
	for _, version := range serviceVersions {
		if *version.Version == ServiceVersion1_1 {
			return ServiceVersion1_1
		}
	}
	return baseServiceVersion
}

func uploadInChunks(fullPath string, fileToUpload os.File, serviceID string, slmpClient slmpclient.SlmpClientOperations) ([]*models.File, error) {
	// Upload the file
	fileToUploadParts, err := util.SplitFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("Could not process file %q: %v", fullPath, err)
	}
	defer attemptToRemoveFileParts(fileToUploadParts)

	var uploaderGroup errgroup.Group
	uploadedFilesChannel := make(chan *models.File)
	defer close(uploadedFilesChannel)
	for _, fileToUploadPart := range fileToUploadParts {
		filePart, err := os.Open(fileToUploadPart)
		if err != nil {
			return nil, fmt.Errorf("Could not open file part %s of file %s", filePart.Name(), fullPath)
		}
		uploaderGroup.Go(func() error {
			file, err := uploadFilePart(filePart, fileToUpload.Name(), serviceID, slmpClient)
			if err != nil {
				return err
			}
			uploadedFilesChannel <- file
			return nil
		})
	}
	uploadedFileParts := []*models.File{}
	var retrieverGroup errgroup.Group
	retrieverGroup.Go(func() error {
		for uploadedFile := range uploadedFilesChannel {
			uploadedFileParts = append(uploadedFileParts, uploadedFile)
			if len(uploadedFileParts) == len(fileToUploadParts) {
				break
			}
		}
		return nil
	})

	err = uploaderGroup.Wait()
	if err != nil {
		return nil, err
	}
	err = retrieverGroup.Wait()
	if err != nil {
		return nil, err
	}
	return uploadedFileParts, nil
}

func attemptToRemoveFileParts(fileParts []string) {
	// If more than one file parts exists, then remove them.
	// If there is only one, then this is the archive itself
	if len(fileParts) > 1 {
		for _, filePart := range fileParts {
			filePartAbsPath, err := filepath.Abs(filePart)
			if err != nil {
				ui.Warn("Error retrieving absolute file path of %q", filePart)
			}
			err = os.Remove(filePartAbsPath)
			if err != nil {
				ui.Warn("Error cleaning up temporary files")
			}
		}
	}
}

func uploadFilePart(filePart *os.File, baseFileName string, serviceID string, client slmpclient.SlmpClientOperations) (*models.File, error) {
	return uploadFile(filePart, baseFileName, serviceID, client)
}

func uploadFile(file *os.File, baseFileName string, serviceID string, client slmpclient.SlmpClientOperations) (*models.File, error) {
	createdFiles, err := client.CreateServiceFile(serviceID, *file)
	defer file.Close()
	if err != nil || len(createdFiles.Files) == 0 {
		return nil, fmt.Errorf("Could not create file %s for service %s: %s", terminal.EntityNameColor(baseFileName), terminal.EntityNameColor(serviceID), err)
	}
	return createdFiles.Files[0], nil
}

func isFileAlreadyUploaded(newFilePath string, fileInfo os.FileInfo, oldFiles models.Files, alreadyUploadedFiles *[]*models.File) bool {
	newFileDigests := make(map[string]string)
	for _, oldFile := range oldFiles.Files {
		if *oldFile.FileName != fileInfo.Name() {
			continue
		}
		if newFileDigests[oldFile.DigestAlgorithm] == "" {
			digest, err := util.ComputeFileChecksum(newFilePath, oldFile.DigestAlgorithm)
			if err != nil {
				ui.Failed("Could not compute digest of file %s: %s", terminal.EntityNameColor(newFilePath), err)
			}
			newFileDigests[oldFile.DigestAlgorithm] = strings.ToUpper(digest)
		}
		if newFileDigests[oldFile.DigestAlgorithm] == oldFile.Digest {
			*alreadyUploadedFiles = append(*alreadyUploadedFiles, oldFile)
			ui.Say("Previously uploaded file %s with same digest detected, new upload will be skipped.",
				terminal.EntityNameColor(fileInfo.Name()))
			return true
		}
	}
	return false
}
