package commands_test

import (
	"errors"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/mtaclient/fakes"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/commands"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/configuration/properties"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"strings"
)

var _ = Describe("FileUploader", func() {
	Describe("UploadFiles", func() {
		const testFileName = "test.mtar"
		const namespace = "namespace"

		var fileUploader *commands.FileUploader
		var testFile *os.File
		var testFileAbsolutePath string
		var testFileDigest string
		var oc = testutil.NewUIOutputCapturer()
		var ex = testutil.NewUIExpector()

		fakeMtaClientBuilder := fakes.NewFakeMtaClientBuilder()

		BeforeEach(func() {
			ui.DisableTerminalOutput(true)
			testFile, _ = os.Create(testFileName)
			testFileAbsolutePath, _ = filepath.Abs(testFile.Name())
			testFileDigest, _ = util.ComputeFileChecksum(testFileAbsolutePath, "MD5")
			testFileDigest = strings.ToUpper(testFileDigest)
		})
		var uploadedFiles []*models.FileMetadata
		var status commands.ExecutionStatus
		Context("with non-existing service files and no files to upload", func() {
			It("should return no uploaded files", func() {
				client := fakeMtaClientBuilder.GetMtaFiles([]*models.FileMetadata{}, nil).Build()

				output := oc.CaptureOutput(func() {
					fileUploader = commands.NewFileUploader(client, namespace, properties.DefaultUploadChunkSizeInMB,
						properties.DefaultUploadChunksSequentially, properties.DefaultDisableProgressBar)
					uploadedFiles, status = fileUploader.UploadFiles([]string{})
				})
				ex.ExpectSuccessWithOutput(status.ToInt(), output, []string{""})
				Expect(uploadedFiles).To(BeNil())
			})
		})

		Context("with existing service files and no files to upload", func() {
			It("should return no uploaded files", func() {
				client := fakeMtaClientBuilder.GetMtaFiles([]*models.FileMetadata{&testutil.SimpleFile}, nil).Build()
				var uploadedFiles []*models.FileMetadata
				output := oc.CaptureOutput(func() {
					fileUploader = commands.NewFileUploader(client, namespace, properties.DefaultUploadChunkSizeInMB,
						properties.DefaultUploadChunksSequentially, properties.DefaultDisableProgressBar)
					uploadedFiles, status = fileUploader.UploadFiles([]string{})
				})
				ex.ExpectSuccessWithOutput(status.ToInt(), output, []string{""})
				Expect(uploadedFiles).To(BeNil())
			})
		})

		Context("with non-existing service files and one file to upload", func() {
			It("should return the uploaded file", func() {
				files := []*models.FileMetadata{testutil.GetFile(testFile, testFileDigest, namespace)}
				client := fakeMtaClientBuilder.
					GetMtaFiles([]*models.FileMetadata{}, nil).
					UploadMtaFile(testFile, testutil.GetFile(testFile, testFileDigest, namespace), nil).Build()
				var uploadedFiles []*models.FileMetadata
				output := oc.CaptureOutput(func() {
					fileUploader = commands.NewFileUploader(client, namespace, properties.DefaultUploadChunkSizeInMB,
						properties.DefaultUploadChunksSequentially, properties.DefaultDisableProgressBar)
					uploadedFiles, status = fileUploader.UploadFiles([]string{testFileAbsolutePath})
				})
				Expect(len(uploadedFiles)).To(Equal(1))
				fullPath, _ := filepath.Abs(testFile.Name())
				ex.ExpectSuccessWithOutput(status.ToInt(), output, []string{
					"Uploading 1 files...",
					"  " + fullPath,
					"OK",
				})
				Expect(uploadedFiles).To(Equal(files))
			})
		})

		Context("with existing service files and one file to upload", func() {
			It("should display a message that the file upload will be skipped", func() {
				files := []*models.FileMetadata{testutil.GetFile(testFile, testFileDigest, "namespace")}
				client := fakeMtaClientBuilder.
					GetMtaFiles([]*models.FileMetadata{&testutil.SimpleFile}, nil).
					UploadMtaFile(testFile, testutil.GetFile(testFile, testFileDigest, "namespace"), nil).Build()
				var uploadedFiles []*models.FileMetadata
				output := oc.CaptureOutput(func() {
					fileUploader = commands.NewFileUploader(client, namespace, properties.DefaultUploadChunkSizeInMB,
						properties.DefaultUploadChunksSequentially, properties.DefaultDisableProgressBar)
					uploadedFiles, status = fileUploader.UploadFiles([]string{testFileAbsolutePath})
				})
				ex.ExpectSuccessWithOutput(status.ToInt(), output, []string{
					"Previously uploaded file test.mtar with same digest detected, new upload will be skipped."})
				Expect(len(uploadedFiles)).To(Equal(1))
				Expect(uploadedFiles).To(Equal(files))
			})
		})

		Context("with non-existing service files and one file to upload and service versions returned from the backend", func() {
			It("should return the uploaded file", func() {
				fileMetadata := testutil.GetFile(testFile, testFileDigest, namespace)
				client := fakeMtaClientBuilder.
					GetMtaFiles([]*models.FileMetadata{}, nil).
					UploadMtaFile(testFile, fileMetadata, nil).Build()
				var uploadedFiles []*models.FileMetadata
				output := oc.CaptureOutput(func() {
					fileUploader = commands.NewFileUploader(client, namespace, properties.DefaultUploadChunkSizeInMB,
						properties.DefaultUploadChunksSequentially, properties.DefaultDisableProgressBar)
					uploadedFiles, status = fileUploader.UploadFiles([]string{testFileAbsolutePath})
				})
				Expect(len(uploadedFiles)).To(Equal(1))
				fullPath, _ := filepath.Abs(testFile.Name())
				ex.ExpectSuccessWithOutput(status.ToInt(), output, []string{
					"Uploading 1 files...",
					"  " + fullPath,
					"OK",
				})
				Expect(uploadedFiles).To(Equal([]*models.FileMetadata{fileMetadata}))
			})
		})

		Context("with error returned from the backend", func() {
			It("should return the uploaded file", func() {
				// files := []*models.File{testutil.GetFile("xs2-deploy", *testFile, testFileDigest)}
				client := fakeMtaClientBuilder.
					GetMtaFiles([]*models.FileMetadata{}, nil).
					UploadMtaFile(testFile, &models.FileMetadata{}, errors.New("Unexpected error from the backend")).Build()
				// var uploadedFiles []*models.FileMetadata
				output := oc.CaptureOutput(func() {
					fileUploader = commands.NewFileUploader(client, namespace, properties.DefaultUploadChunkSizeInMB,
						properties.DefaultUploadChunksSequentially, properties.DefaultDisableProgressBar)
					_, status = fileUploader.UploadFiles([]string{testFileAbsolutePath})
				})
				// Expect(len(uploadedFiles)).To(Equal(1))
				// fullPath, _ := filepath.Abs(testFile.Name())
				ex.ExpectFailureOnLine(status.ToInt(), output, "Could not upload file "+testFileAbsolutePath, 3)
				// Expect(uploadedFiles).To(Equal(files))
			})
		})

		AfterEach(func() {
			testFile.Close()
			os.Remove(testFileName)
		})
	})
})
