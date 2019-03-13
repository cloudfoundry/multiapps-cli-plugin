package commands

import (
	"fmt"
	"github.com/cloudfoundry/cli/cf/terminal"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"golang.org/x/sync/errgroup"
)

//FileUploader uploads files for the service with the specified service ID
type FileUploader struct {
	files     []string
	mtaClient mtaclient.MtaClientOperations
}

//NewFileUploader creates a new file uploader for the specified service ID, files, and SLMP client
func NewFileUploader(files []string, mtaClient mtaclient.MtaClientOperations) *FileUploader {
	return &FileUploader{
		files:     files,
		mtaClient: mtaClient,
	}
}

//UploadFiles uploads the files
func (f *FileUploader) UploadFiles() ([]*models.FileMetadata, ExecutionStatus) {
	log.Tracef("Uploading files '%v'\n", f.files)

	// Get all files that are already uploaded
	uploadedMtaFiles, err := f.mtaClient.GetMtaFiles()
	if err != nil {
		ui.Failed("Could not get mta files: %s", baseclient.NewClientError(err))
		return nil, Failure
	}

	// Determine which files to uplaod
	filesToUpload := []os.File{}
	alreadyUploadedFiles := []*models.FileMetadata{}
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
		if !isFileAlreadyUploaded(file, fileInfo, uploadedMtaFiles, &alreadyUploadedFiles) {
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
	uploadedFiles := []*models.FileMetadata{}
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
			uploaded, err := uploadInChunks(fullPath, fileToUpload, f.mtaClient)
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

func uploadInChunks(fullPath string, fileToUpload os.File, mtaClient mtaclient.MtaClientOperations) ([]*models.FileMetadata, error) {
	// Upload the file
	fileToUploadParts, err := util.SplitFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("Could not process file %q: %v", fullPath, baseclient.NewClientError(err))
	}
	defer attemptToRemoveFileParts(fileToUploadParts)

	var uploaderGroup errgroup.Group
	uploadedFilesChannel := make(chan *models.FileMetadata)
	defer close(uploadedFilesChannel)
	for _, fileToUploadPart := range fileToUploadParts {
		filePart, err := os.Open(fileToUploadPart)
		if err != nil {
			return nil, fmt.Errorf("Could not open file part %s of file %s", filePart.Name(), fullPath)
		}
		uploaderGroup.Go(func() error {
			file, err := uploadFilePart(filePart, fileToUpload.Name(), mtaClient)
			if err != nil {
				return err
			}
			uploadedFilesChannel <- file
			return nil
		})
	}
	uploadedFileParts := []*models.FileMetadata{}
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

func uploadFilePart(filePart *os.File, baseFileName string, client mtaclient.MtaClientOperations) (*models.FileMetadata, error) {
	uploadedFile, err := client.UploadMtaFile(*filePart)
	defer filePart.Close()
	if err != nil {
		return nil, fmt.Errorf("Could not create file %s: %s", terminal.EntityNameColor(baseFileName), baseclient.NewClientError(err))
	}
	return uploadedFile, nil
}

func isFileAlreadyUploaded(newFilePath string, fileInfo os.FileInfo, oldFiles []*models.FileMetadata, alreadyUploadedFiles *[]*models.FileMetadata) bool {
	newFileDigests := make(map[string]string)
	for _, oldFile := range oldFiles {
		if oldFile.Name != fileInfo.Name() {
			continue
		}
		if newFileDigests[oldFile.DigestAlgorithm] == "" {
			digest, err := util.ComputeFileChecksum(newFilePath, oldFile.DigestAlgorithm)
			if err != nil {
				ui.Failed("Could not compute digest of file %s: %s", terminal.EntityNameColor(newFilePath), baseclient.NewClientError(err))
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
