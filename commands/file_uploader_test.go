package commands_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/SAP/cf-mta-plugin/clients/models"
	"github.com/SAP/cf-mta-plugin/clients/slmpclient/fakes"
	"github.com/SAP/cf-mta-plugin/commands"
	"github.com/SAP/cf-mta-plugin/testutil"
	"github.com/SAP/cf-mta-plugin/ui"
	"github.com/SAP/cf-mta-plugin/util"
)

var _ = Describe("FileUploader", func() {
	Describe("UploadFiles", func() {
		const testFileName = "test.mtar"

		var fileUploader *commands.FileUploader
		var testFile *os.File
		var testFileAbsolutePath string
		var testFileDigest string
		var oc = testutil.NewUIOutputCapturer()
		var ex = testutil.NewUIExpector()

		fakeSlmpClientBuilder := fakes.NewFakeSlmpClientBuilder()

		BeforeEach(func() {
			ui.DisableTerminalOutput(true)
			testFile, _ = os.Create(testFileName)
			testFileAbsolutePath, _ = filepath.Abs(testFile.Name())
			testFileDigest, _ = util.ComputeFileChecksum(testFileAbsolutePath, "MD5")
			testFileDigest = strings.ToUpper(testFileDigest)
		})
		var uploadedFiles []*models.File
		var status commands.ExecutionStatus
		Context("with non-existing service files and no files to upload", func() {
			It("should return no uploaded files", func() {
				client := fakeSlmpClientBuilder.GetServiceFiles("xs2-deploy", models.Files{}, nil).Build()

				output := oc.CaptureOutput(func() {
					fileUploader = commands.NewFileUploader("xs2-deploy", []string{}, client)
					uploadedFiles, status = fileUploader.UploadFiles()
				})
				ex.ExpectSuccess(status.ToInt(), output)
				Expect(uploadedFiles).To(Equal([]*models.File{}))
			})
		})

		Context("with existing service files and no files to upload", func() {
			It("should return no uploaded files", func() {
				client := fakeSlmpClientBuilder.GetServiceFiles("xs2-deploy", testutil.FilesResult, nil).Build()
				var uploadedFiles []*models.File
				output := oc.CaptureOutput(func() {
					fileUploader = commands.NewFileUploader("xs2-deploy", []string{}, client)
					uploadedFiles, status = fileUploader.UploadFiles()
				})
				ex.ExpectSuccess(status.ToInt(), output)
				Expect(uploadedFiles).To(Equal([]*models.File{}))
			})
		})

		Context("with non-existing service files and one file to upload", func() {
			It("should return the uploaded file", func() {
				files := []*models.File{testutil.GetFile("xs2-deploy", *testFile, testFileDigest)}
				client := fakeSlmpClientBuilder.
					GetMetadata(&testutil.SlmpMetadataResult, nil).
					GetServiceFiles("xs2-deploy", models.Files{}, nil).
					CreateServiceFile("xs2-deploy", testFile, testutil.GetFiles(files), nil).Build()
				var uploadedFiles []*models.File
				output := oc.CaptureOutput(func() {
					fileUploader = commands.NewFileUploader("xs2-deploy", []string{testFileAbsolutePath}, client)
					uploadedFiles, status = fileUploader.UploadFiles()
				})
				Expect(len(uploadedFiles)).To(Equal(1))
				fullPath, _ := filepath.Abs(testFile.Name())
				ex.ExpectSuccessWithOutput(status.ToInt(), output, []string{
					"Uploading 1 files...\n",
					"  " + fullPath + "\n",
					"OK\n",
				})
				Expect(uploadedFiles).To(Equal(files))
			})
		})

		Context("with existing service files and one file to upload", func() {
			It("should display a message that the file upload will be skipped", func() {
				files := []*models.File{testutil.GetFile("xs2-deploy", *testFile, testFileDigest)}
				client := fakeSlmpClientBuilder.
					GetServiceFiles("xs2-deploy", testutil.FilesResult, nil).
					CreateServiceFile("xs2-deploy", testFile, testutil.GetFiles(files), nil).Build()
				var uploadedFiles []*models.File
				output := oc.CaptureOutput(func() {
					fileUploader = commands.NewFileUploader("xs2-deploy", []string{testFileAbsolutePath}, client)
					uploadedFiles, status = fileUploader.UploadFiles()
				})
				ex.ExpectSuccessWithOutput(status.ToInt(), output, []string{
					"Previously uploaded file test.mtar with same digest detected, new upload will be skipped.\n"})
				Expect(len(uploadedFiles)).To(Equal(1))
				Expect(uploadedFiles).To(Equal(files))
			})
		})

		Context("with non-existing service files and one file to upload and service versions returned from the backend", func() {
			It("should return the uploaded file", func() {
				files := []*models.File{testutil.GetFile("xs2-deploy", *testFile, testFileDigest)}
				client := fakeSlmpClientBuilder.
					GetMetadata(&testutil.SlmpMetadataResult, nil).
					GetServiceFiles("xs2-deploy", models.Files{}, nil).
					CreateServiceFile("xs2-deploy", testFile, testutil.GetFiles(files), nil).
					GetServiceVersions("xs2-deploy", testutil.ServiceVersion1_1, nil).Build()
				var uploadedFiles []*models.File
				output := oc.CaptureOutput(func() {
					fileUploader = commands.NewFileUploader("xs2-deploy", []string{testFileAbsolutePath}, client)
					uploadedFiles, status = fileUploader.UploadFiles()
				})
				Expect(len(uploadedFiles)).To(Equal(1))
				fullPath, _ := filepath.Abs(testFile.Name())
				ex.ExpectSuccessWithOutput(status.ToInt(), output, []string{
					"Uploading 1 files...\n",
					"  " + fullPath + "\n",
					"OK\n",
				})
				Expect(uploadedFiles).To(Equal(files))
			})
		})

		Context("with error returned from the backend", func() {
			It("should return the uploaded file", func() {
				// files := []*models.File{testutil.GetFile("xs2-deploy", *testFile, testFileDigest)}
				client := fakeSlmpClientBuilder.
					GetMetadata(&testutil.SlmpMetadataResult, nil).
					GetServiceFiles("xs2-deploy", models.Files{}, nil).
					CreateServiceFile("xs2-deploy", testFile, models.Files{}, errors.New("Unexpected error from the backend")).
					GetServiceVersions("xs2-deploy", testutil.ServiceVersion1_1, nil).Build()
				var uploadedFiles []*models.File
				output := oc.CaptureOutput(func() {
					fileUploader = commands.NewFileUploader("xs2-deploy", []string{testFileAbsolutePath}, client)
					uploadedFiles, status = fileUploader.UploadFiles()
				})
				// Expect(len(uploadedFiles)).To(Equal(1))
				// fullPath, _ := filepath.Abs(testFile.Name())
				ex.ExpectFailureOnLine(status.ToInt(), output, "Could not upload file "+testFileAbsolutePath, 2)
				// Expect(uploadedFiles).To(Equal(files))
			})
		})

		AfterEach(func() {
			os.RemoveAll(testFileName)
		})
	})
})
