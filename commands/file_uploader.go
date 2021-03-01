package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry/cli/cf/terminal"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
)

//FileUploader uploads files in chunks for the specified namespace
type FileUploader struct {
	mtaClient           mtaclient.MtaClientOperations
	namespace           string
	uploadChunkSizeInMB uint64
}

//NewFileUploader creates a new file uploader for the specified namespace
func NewFileUploader(mtaClient mtaclient.MtaClientOperations, namespace string, uploadChunkSizeInMB uint64) *FileUploader {
	return &FileUploader{
		mtaClient:           mtaClient,
		namespace:           namespace,
		uploadChunkSizeInMB: uploadChunkSizeInMB,
	}
}

//UploadFiles uploads the files
func (f *FileUploader) UploadFiles(files []string) ([]*models.FileMetadata, ExecutionStatus) {
	log.Tracef("Uploading files '%v'\n", files)

	// Get all files that are already uploaded
	uploadedMtaFiles, err := f.mtaClient.GetMtaFiles(&f.namespace)
	if err != nil {
		ui.Failed("Could not get mta files: %s", baseclient.NewClientError(err))
		return nil, Failure
	}

	// Determine which files to upload
	var filesToUpload []os.File
	var alreadyUploadedFiles []*models.FileMetadata
	for _, file := range files {
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
		if !f.isFileAlreadyUploaded(file, fileInfo, uploadedMtaFiles, &alreadyUploadedFiles) {
			// If not, add it to the list of uploaded files
			fileToUpload, err := os.Open(file)
			if err != nil {
				ui.Failed("Could not open file %s", terminal.EntityNameColor(file))
				return nil, Failure
			}
			defer fileToUpload.Close()
			filesToUpload = append(filesToUpload, *fileToUpload)
		}
	}

	// If there are new files to upload, upload them
	var uploadedFiles []*models.FileMetadata
	uploadedFiles = append(uploadedFiles, alreadyUploadedFiles...)
	if len(filesToUpload) != 0 {
		ui.Say("Uploading %d files...", len(filesToUpload))

		// Iterate over all files to be uploaded
		for _, fileToUpload := range filesToUpload {
			// Print the full path of the file
			fullPath, err := filepath.Abs(fileToUpload.Name())
			if err != nil {
				ui.Failed("Could not get absolute path of file %s: %s", terminal.EntityNameColor(fileToUpload.Name()), err.Error())
				return nil, Failure
			}
			ui.Say("  " + fullPath)

			// Upload the file
			uploaded, err := f.uploadInChunks(fullPath, fileToUpload)
			if err != nil {
				ui.Failed("Could not upload file %s: %s", terminal.EntityNameColor(fileToUpload.Name()), err.Error())
				return nil, Failure
			}
			uploadedFiles = append(uploadedFiles, uploaded...)
		}
		ui.Ok()
	}
	return uploadedFiles, Success
}

func (f *FileUploader) uploadInChunks(fullPath string, fileToUpload os.File) ([]*models.FileMetadata, error) {
	// Upload the file
	err := util.ValidateChunkSize(fullPath, f.uploadChunkSizeInMB)
	if err != nil {
		return nil, fmt.Errorf("Could not valide file %q: %v", fullPath, baseclient.NewClientError(err))
	}
	fileToUploadParts, err := util.SplitFile(fullPath, f.uploadChunkSizeInMB)
	if err != nil {
		return nil, fmt.Errorf("Could not process file %q: %v", fullPath, baseclient.NewClientError(err))
	}
	defer attemptToRemoveFileParts(fileToUploadParts)

	uploadedFilesChannel := make(chan *models.FileMetadata)
	defer close(uploadedFilesChannel)
	errorChannel := make(chan error)
	defer close(errorChannel)

	for _, fileToUploadPart := range fileToUploadParts {
		filePart, err := os.Open(fileToUploadPart)
		if err != nil {
			return nil, fmt.Errorf("Could not open file part %s of file %s", fileToUploadPart, fullPath)
		}
		go func() {
			fileInfo, err := filePart.Stat()
			if err != nil {
				errorChannel <- err
				return
			}
			file, err := f.uploadFilePart(filePart, fileToUpload.Name(), fileInfo.Size())
			if err != nil {
				errorChannel <- err
				return
			}
			uploadedFilesChannel <- file
		}()
	}

	var uploadedFileParts []*models.FileMetadata
	for len(uploadedFileParts) < len(fileToUploadParts) {
		select {
		case uploadedFile := <-uploadedFilesChannel:
			uploadedFileParts = append(uploadedFileParts, uploadedFile)
		case err := <-errorChannel:
			return nil, err
		}
	}
	return uploadedFileParts, nil
}

func attemptToRemoveFileParts(fileParts []string) {
	// If more than one file parts exists, then remove them.
	// If there is only one, then this is the archive itself
	if len(fileParts) <= 1 {
		return
	}
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

func (f *FileUploader) uploadFilePart(filePart *os.File, baseFileName string, fileSize int64) (*models.FileMetadata, error) {
	uploadedFile, err := f.mtaClient.UploadMtaFile(*filePart, fileSize, &f.namespace)
	defer filePart.Close()
	if err != nil {
		return nil, fmt.Errorf("Could not create file %s: %s", terminal.EntityNameColor(baseFileName), baseclient.NewClientError(err))
	}
	return uploadedFile, nil
}

func (f *FileUploader) isFileAlreadyUploaded(newFilePath string, fileInfo os.FileInfo, oldFiles []*models.FileMetadata, alreadyUploadedFiles *[]*models.FileMetadata) bool {
	for _, oldFile := range oldFiles {
		if oldFile.Name != fileInfo.Name() || oldFile.Namespace != f.namespace {
			continue
		}
		digest, err := util.ComputeFileChecksum(newFilePath, oldFile.DigestAlgorithm)
		if err != nil {
			ui.Failed("Could not compute digest of file %s: %s", terminal.EntityNameColor(newFilePath), baseclient.NewClientError(err))
			return false
		}

		if strings.ToUpper(digest) == oldFile.Digest {
			*alreadyUploadedFiles = append(*alreadyUploadedFiles, oldFile)
			ui.Say("Previously uploaded file %s with same digest detected, new upload will be skipped.",
				terminal.EntityNameColor(fileInfo.Name()))
			return true
		}
	}
	return false
}
