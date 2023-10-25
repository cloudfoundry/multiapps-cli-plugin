package commands

import (
	"fmt"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"code.cloudfoundry.org/cli/cf/terminal"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/baseclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
)

// FileUploader uploads files in chunks for the specified namespace
type FileUploader struct {
	mtaClient                mtaclient.MtaClientOperations
	namespace                string
	uploadChunkSizeInMB      uint64
	sequentialUpload         bool
	shouldDisableProgressBar bool
}

type progressBarReader struct {
	pb       *pb.ProgressBar
	written  atomic.Int64
	file     io.ReadSeeker
	fileName string
}

func (r *progressBarReader) Name() string {
	return r.fileName
}

func (r *progressBarReader) Read(p []byte) (int, error) {
	n, err := r.file.Read(p)
	if n > 0 {
		r.written.Add(int64(n))
		r.pb.Add(n)
	}
	return n, err
}

func (r *progressBarReader) Seek(offset int64, whence int) (int64, error) {
	newOffset, err := r.file.Seek(offset, whence)
	if whence == io.SeekStart && r.written.Load() != 0 {
		r.pb.Add64(-r.written.Load() - offset)
		r.written.Store(offset)
	}
	return newOffset, err
}

func (r *progressBarReader) Close() error {
	//no-op as we close the file part manually
	return nil
}

// NewFileUploader creates a new file uploader for the specified namespace
func NewFileUploader(mtaClient mtaclient.MtaClientOperations, namespace string, uploadChunkSizeInMB uint64,
	sequentialUpload, shouldDisableProgressBar bool) *FileUploader {
	return &FileUploader{
		mtaClient:                mtaClient,
		namespace:                namespace,
		uploadChunkSizeInMB:      uploadChunkSizeInMB,
		sequentialUpload:         sequentialUpload,
		shouldDisableProgressBar: shouldDisableProgressBar,
	}
}

// UploadFiles uploads the files
func (f *FileUploader) UploadFiles(files []string) ([]*models.FileMetadata, ExecutionStatus) {
	log.Tracef("Uploading files '%v'\n", files)

	// Get all files that are already uploaded
	uploadedMtaFiles, err := f.mtaClient.GetMtaFiles(&f.namespace)
	if err != nil {
		ui.Failed("Could not get mta files: %s", baseclient.NewClientError(err))
		return nil, Failure
	}

	// Determine which files to upload
	var filesToUpload []*os.File
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
			filesToUpload = append(filesToUpload, fileToUpload)
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
			// as per the Go docs, the file name is the one passed to os.Open
			// and we pass the absolute path
			ui.Say("  " + fileToUpload.Name())
			// Upload the file
			uploaded, err := f.uploadInChunks(fileToUpload)
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

func (f *FileUploader) uploadInChunks(fileToUpload *os.File) ([]*models.FileMetadata, error) {
	err := util.ValidateChunkSize(fileToUpload.Name(), f.uploadChunkSizeInMB)
	if err != nil {
		return nil, fmt.Errorf("Could not valide file %q: %v", fileToUpload.Name(), err)
	}
	fileToUploadParts, err := util.SplitFile(fileToUpload.Name(), f.uploadChunkSizeInMB)
	if err != nil {
		return nil, fmt.Errorf("Could not process file %q: %v", fileToUpload.Name(), err)
	}
	defer attemptToRemoveFileParts(fileToUploadParts)

	var uploadedFileParts []*models.FileMetadata
	uploadedFilesChannel := make(chan *models.FileMetadata)
	errorChannel := make(chan error)

	fileInfo, err := fileToUpload.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not get information on file %q: %v", fileToUpload.Name(), err)
	}

	progressBar := pb.New64(fileInfo.Size()).SetUnits(pb.U_BYTES)
	progressBar.ShowTimeLeft = false
	progressBar.ShowElapsedTime = true
	progressBar.NotPrint = f.shouldDisableProgressBar
	progressBar.Start()
	defer progressBar.Finish()

	for _, fileToUploadPart := range fileToUploadParts {
		filePartCopy := fileToUploadPart
		go func() {
			file, err := f.uploadFilePart(filePartCopy, fileToUpload.Name(), progressBar)
			if err != nil {
				errorChannel <- err
				return
			}
			uploadedFilesChannel <- file
		}()
		if f.sequentialUpload {
			if err := waitForChannelData(uploadedFilesChannel, errorChannel, &uploadedFileParts); err != nil {
				return nil, err
			}
		}
	}

	for len(uploadedFileParts) < len(fileToUploadParts) {
		if err := waitForChannelData(uploadedFilesChannel, errorChannel, &uploadedFileParts); err != nil {
			return nil, err
		}
	}
	return uploadedFileParts, nil
}

func waitForChannelData(fileChan <-chan *models.FileMetadata, errChan <-chan error, result *[]*models.FileMetadata) error {
	select {
	case uploadedFile := <-fileChan:
		*result = append(*result, uploadedFile)
	case err := <-errChan:
		return err
	}
	return nil
}

func attemptToRemoveFileParts(fileParts []*os.File) {
	// If more than one file parts exists, then remove them.
	// If there is only one, then this is the archive itself
	if len(fileParts) <= 1 {
		return
	}
	for _, filePart := range fileParts {
		filePartAbsPath, err := filepath.Abs(filePart.Name())
		if err != nil {
			ui.Warn("Error retrieving absolute file path of %q", filePart.Name())
		}
		err = os.Remove(filePartAbsPath)
		if err != nil {
			ui.Warn("Error cleaning up temporary files")
		}
	}
}

func (f *FileUploader) uploadFilePart(filePart *os.File, fileName string, pb *pb.ProgressBar) (*models.FileMetadata, error) {
	defer filePart.Close()
	fileInfo, err := filePart.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not stat file part %q of file %q", filePart.Name(), fileName)
	}

	file := &progressBarReader{file: filePart, fileName: fileInfo.Name(), pb: pb}
	uploadedFile, err := f.mtaClient.UploadMtaFile(file, fileInfo.Size(), &f.namespace)
	if err != nil {
		return nil, fmt.Errorf("could not upload file %s: %s", terminal.EntityNameColor(fileName), err)
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
