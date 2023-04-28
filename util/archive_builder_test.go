package util_test

import (
	"archive/zip"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/testutil"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("ArchiveBuilder", func() {
	Describe("Build", func() {
		var tempDirLocation string
		BeforeEach(func() {
			tempDirLocation, _ = os.MkdirTemp("", "archive-builder")
		})
		Context("With not existing resources", func() {
			It("should try to find the directory and fail with error", func() {
				_, err := util.NewMtaArchiveBuilder([]string{}, []string{}).Build("not-existing-location")
				Expect(err).To(MatchError("Deployment descriptor location does not exist not-existing-location"))
			})
			It("should try to find the deployment descriptor in the provided location and fail with error", func() {
				_, err := util.NewMtaArchiveBuilder([]string{}, []string{}).Build(tempDirLocation)
				Expect(err).To(MatchError("No deployment descriptor with name mtad.yaml was found in location " + tempDirLocation))
			})
		})

		Context("With different paths relative to the deployment descriptor", func() {
			var (
				currentWorkingDirectory string
				err                     error
				mtaArchiveLocation      string
			)
			const requiredDependencyContent = "test-module-content"
			const testDeploymentDescriptor = "mtad.yaml"
			const moduleName = "TestModule"

			BeforeEach(func() {
				// need to cd into the tempDir in order to simulate the relative path
				currentWorkingDirectory, err = os.Getwd()
				Expect(err).To(BeNil())
				err = os.Chdir(tempDirLocation)
				Expect(err).To(BeNil())

				os.Create(requiredDependencyContent)
				os.WriteFile(requiredDependencyContent, []byte("this is a test module content"), os.ModePerm)
				descriptor := util.MtaDeploymentDescriptor{SchemaVersion: "100", ID: "test", Modules: []util.Module{
					util.Module{Name: moduleName, Path: requiredDependencyContent},
				}}
				generatedYamlBytes, _ := yaml.Marshal(descriptor)

				os.WriteFile(testDeploymentDescriptor, generatedYamlBytes, os.ModePerm)
			})

			It("Should find deployment descriptor with \".\" baseDirectory path", func() {
				mtaArchiveLocation, err = util.NewMtaArchiveBuilder([]string{moduleName}, []string{}).Build(".")
			})

			It("Should find deployment descriptor with \"./\" baseDirectory path", func() {
				mtaArchiveLocation, err = util.NewMtaArchiveBuilder([]string{moduleName}, []string{}).Build("./")
			})

			It("Should find deployment descriptor with \"../\" baseDirectory path", func() {
				// create and cd into new dir
				err = os.Mkdir("test", 0700)
				Expect(err).To(BeNil())
				err = os.Chdir("test")
				Expect(err).To(BeNil())
				mtaArchiveLocation, err = util.NewMtaArchiveBuilder([]string{moduleName}, []string{}).Build("../")
			})

			AfterEach(func() {
				expectedModulePath := filepathUnixJoin(moduleName, requiredDependencyContent)
				Expect(err).To(BeNil())
				_, err = os.Stat(mtaArchiveLocation)
				Expect(err).To(BeNil())
				Expect(isInArchive(expectedModulePath, mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/MANIFEST.MF", mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/mtad.yaml", mtaArchiveLocation)).To(BeTrue())
				Expect(isManifestValid("META-INF/MANIFEST.MF", map[string]string{"MTA-Module": "TestModule", "Name": expectedModulePath}, mtaArchiveLocation)).To(Equal(map[string]string{"MTA-Module": "TestModule", "Name": expectedModulePath}))
				defer os.Remove(mtaArchiveLocation)
				defer os.Chdir(currentWorkingDirectory)
			})
		})

		Context("With deployment descriptor which contains some modules and resources", func() {
			It("Try to parse the specified modules and fail as the paths are not existing", func() {
				descriptor := util.MtaDeploymentDescriptor{SchemaVersion: "100", ID: "test", Modules: []util.Module{
					util.Module{Name: "TestModule", Path: "not-existing-path"},
				}}
				generatedYamlBytes, _ := yaml.Marshal(descriptor)
				testDeploymentDescriptor := tempDirLocation + string(os.PathSeparator) + "mtad.yaml"
				os.WriteFile(testDeploymentDescriptor, generatedYamlBytes, os.ModePerm)
				_, err := util.NewMtaArchiveBuilder([]string{"TestModule"}, []string{}).Build(tempDirLocation)
				Expect(err.Error()).To(MatchRegexp("Error building MTA Archive: file path .*?not-existing-path not found"))
			})

			It("Try to parse the specified resources and fail as the paths are not existing", func() {
				descriptor := util.MtaDeploymentDescriptor{SchemaVersion: "100", ID: "test", Modules: []util.Module{}, Resources: []util.Resource{
					util.Resource{Name: "foo", Type: "Some type", Parameters: map[string]interface{}{
						"path": "not-existing-resource-path",
					}},
				}}
				generatedYamlBytes, _ := yaml.Marshal(descriptor)
				testDeploymentDescriptor := tempDirLocation + string(os.PathSeparator) + "mtad.yaml"
				os.WriteFile(testDeploymentDescriptor, generatedYamlBytes, os.ModePerm)
				_, err := util.NewMtaArchiveBuilder([]string{}, []string{"foo"}).Build(tempDirLocation)
				Expect(err.Error()).To(MatchRegexp("Error building MTA Archive: file path .*?not-existing-resource-path not found"))
			})

			It("Try to parse the specified required dependencies config paths and fail as the paths are not existing", func() {
				requiredDependencyContent := filepath.Join(tempDirLocation, "test-module-1-content")
				os.Create(requiredDependencyContent)
				os.WriteFile(requiredDependencyContent, []byte("this is a test module content"), os.ModePerm)
				descriptor := util.MtaDeploymentDescriptor{SchemaVersion: "100", ID: "test", Modules: []util.Module{
					util.Module{Name: "TestModule", Path: "test-module-1-content", RequiredDependencies: []util.RequiredDependency{
						util.RequiredDependency{Name: "foo", Parameters: map[string]interface{}{
							"path": "not-existing-required-dependency-path",
						}},
					}},
				}}
				generatedYamlBytes, _ := yaml.Marshal(descriptor)
				testDeploymentDescriptor := tempDirLocation + string(os.PathSeparator) + "mtad.yaml"
				os.WriteFile(testDeploymentDescriptor, generatedYamlBytes, os.ModePerm)
				_, err := util.NewMtaArchiveBuilder([]string{"TestModule"}, []string{}).Build(tempDirLocation)
				Expect(err.Error()).To(MatchRegexp("Error building MTA Archive: file path .*?not-existing-required-dependency-path not found"))
			})
		})

		Context("With deployment descriptor which contains some modules and resources and not valid modules or resources", func() {
			It("Try to parse the specified modules and fail as the modules are not presented in the descriptor", func() {
				descriptor := util.MtaDeploymentDescriptor{SchemaVersion: "100", ID: "test", Modules: []util.Module{
					util.Module{Name: "foo", Path: "not-existing-path"},
					util.Module{Name: "bar", Path: "not-existing-path"},
					util.Module{Name: "baz", Path: "not-existing-path"},
					util.Module{Name: "baz-foo", Path: "not-existing-path"},
				}}
				generatedYamlBytes, _ := yaml.Marshal(descriptor)
				testDeploymentDescriptor := tempDirLocation + string(os.PathSeparator) + "mtad.yaml"
				os.WriteFile(testDeploymentDescriptor, generatedYamlBytes, os.ModePerm)
				_, err := util.NewMtaArchiveBuilder([]string{"foo", "bar", "test-1", "test-2"}, []string{}).Build(tempDirLocation)
				Expect(err.Error()).To(MatchRegexp("Error building MTA Archive: Modules test-1, test-2 are specified for deployment but are not part of deployment descriptor modules"))
			})

			It("Try to parse the specified resources and fail as the resources are not part of deployment descriptor", func() {
				descriptor := util.MtaDeploymentDescriptor{SchemaVersion: "100", ID: "test", Modules: []util.Module{}, Resources: []util.Resource{
					util.Resource{Name: "foo", Type: "Some type", Parameters: map[string]interface{}{
						"path": "not-existing-resource-path",
					}},
					util.Resource{Name: "bar", Type: "Some type", Parameters: map[string]interface{}{
						"path": "not-existing-resource-path",
					}},
					util.Resource{Name: "baz", Type: "Some type", Parameters: map[string]interface{}{
						"path": "not-existing-resource-path",
					}},
					util.Resource{Name: "baz-foo", Type: "Some type", Parameters: map[string]interface{}{
						"path": "not-existing-resource-path",
					}},
				}}
				generatedYamlBytes, _ := yaml.Marshal(descriptor)
				testDeploymentDescriptor := tempDirLocation + string(os.PathSeparator) + "mtad.yaml"
				os.WriteFile(testDeploymentDescriptor, generatedYamlBytes, os.ModePerm)
				_, err := util.NewMtaArchiveBuilder([]string{}, []string{"foo", "bar", "testing", "not-existing"}).Build(tempDirLocation)
				Expect(err.Error()).To(MatchRegexp("Error building MTA Archive: Resources testing, not-existing are specified for deployment but are not part of deployment descriptor resources"))
			})
		})

		Context("With deployment descriptor which does not contain any path path param", func() {
			var oc = testutil.NewUIOutputCapturer()
			It("Should try to resolve the modules and report that they do not have path params.", func() {
				descriptor := util.MtaDeploymentDescriptor{SchemaVersion: "100", ID: "test", Modules: []util.Module{
					util.Module{Name: "TestModule"},
					util.Module{Name: "TestModule1"},
				}}

				generatedYamlBytes, _ := yaml.Marshal(descriptor)
				testDeploymentDescriptor := tempDirLocation + string(os.PathSeparator) + "mtad.yaml"
				os.WriteFile(testDeploymentDescriptor, generatedYamlBytes, os.ModePerm)
				output := oc.CaptureOutput(func() {
					util.NewMtaArchiveBuilder([]string{"TestModule", "TestModule1"}, []string{}).Build(tempDirLocation)
				})
				Expect(output[0]).To(Equal("Modules TestModule, TestModule1 do not have a path, specified for their binaries and will be ignored"))
			})

		})

		Context("With deployment descriptor which contains only valid modules", func() {
			It("Should build the MTA Archive containing the valid modules", func() {
				requiredDependencyContent := filepath.Join(tempDirLocation, "test-module-1-content")
				os.Create(requiredDependencyContent)
				os.WriteFile(requiredDependencyContent, []byte("this is a test module content"), os.ModePerm)
				descriptor := util.MtaDeploymentDescriptor{SchemaVersion: "100", ID: "test", Modules: []util.Module{
					util.Module{Name: "TestModule", Path: "test-module-1-content"},
				}}
				generatedYamlBytes, _ := yaml.Marshal(descriptor)
				testDeploymentDescriptor := tempDirLocation + string(os.PathSeparator) + "mtad.yaml"
				os.WriteFile(testDeploymentDescriptor, generatedYamlBytes, os.ModePerm)
				mtaArchiveLocation, err := util.NewMtaArchiveBuilder([]string{"TestModule"}, []string{}).Build(tempDirLocation)
				defer os.Remove(mtaArchiveLocation)
				Expect(err).To(BeNil())
				_, err = os.Stat(mtaArchiveLocation)
				Expect(err).To(BeNil())
				expectedModulePath := filepathUnixJoin("TestModule", "test-module-1-content")
				Expect(isInArchive(expectedModulePath, mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/MANIFEST.MF", mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/mtad.yaml", mtaArchiveLocation)).To(BeTrue())
				Expect(isManifestValid("META-INF/MANIFEST.MF", map[string]string{"MTA-Module": "TestModule", "Name": expectedModulePath}, mtaArchiveLocation)).To(Equal(map[string]string{"MTA-Module": "TestModule", "Name": expectedModulePath}))
			})
		})

		Context("With deployment descriptor which contains valid modules with non-normalized paths", func() {
			It("Should build the MTA Archive containing the valid modules", func() {
				requiredDependencyContent := filepath.Join(tempDirLocation, "test-module-1-content")
				os.Create(requiredDependencyContent)
				os.WriteFile(requiredDependencyContent, []byte("this is a test module content"), os.ModePerm)
				descriptor := util.MtaDeploymentDescriptor{SchemaVersion: "100", ID: "test", Modules: []util.Module{
					util.Module{Name: "TestModule", Path: "../test-module-1-content"},
				}}
				generatedYamlBytes, _ := yaml.Marshal(descriptor)
				mtadDirectory := filepath.Join(tempDirLocation, "test")
				os.MkdirAll(mtadDirectory, os.ModePerm)
				testDeploymentDescriptor := filepath.Join(mtadDirectory, "mtad.yaml")
				os.WriteFile(testDeploymentDescriptor, generatedYamlBytes, os.ModePerm)
				mtaArchiveLocation, err := util.NewMtaArchiveBuilder([]string{"TestModule"}, []string{}).Build(mtadDirectory)
				defer os.Remove(mtaArchiveLocation)
				Expect(err).To(BeNil())
				_, err = os.Stat(mtaArchiveLocation)
				Expect(err).To(BeNil())
				expectedModulePath := filepathUnixJoin("TestModule", "test-module-1-content")
				Expect(isInArchive(expectedModulePath, mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/MANIFEST.MF", mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/mtad.yaml", mtaArchiveLocation)).To(BeTrue())
				Expect(isManifestValid("META-INF/MANIFEST.MF", map[string]string{"MTA-Module": "TestModule", "Name": expectedModulePath}, mtaArchiveLocation)).To(Equal(map[string]string{"MTA-Module": "TestModule", "Name": expectedModulePath}))
			})
		})

		Context("With deployment descriptor which contains only valid modules with same paths", func() {
			It("should build the MTA Archive containing the valid modules", func() {
				requiredDependencyContent := filepath.Join(tempDirLocation, "test-module-1-content")
				os.Create(requiredDependencyContent)
				os.WriteFile(requiredDependencyContent, []byte("this is a test module content"), os.ModePerm)
				descriptor := util.MtaDeploymentDescriptor{SchemaVersion: "100", ID: "test", Modules: []util.Module{
					util.Module{Name: "TestModule", Path: "test-module-1-content"},
					util.Module{Name: "TestModule1", Path: "test-module-1-content"},
				}}
				generatedYamlBytes, _ := yaml.Marshal(descriptor)
				testDeploymentDescriptor := tempDirLocation + string(os.PathSeparator) + "mtad.yaml"
				os.WriteFile(testDeploymentDescriptor, generatedYamlBytes, os.ModePerm)
				mtaArchiveLocation, err := util.NewMtaArchiveBuilder([]string{"TestModule", "TestModule1"}, []string{}).Build(tempDirLocation)
				defer os.Remove(mtaArchiveLocation)
				Expect(err).To(BeNil())
				_, err = os.Stat(mtaArchiveLocation)
				Expect(err).To(BeNil())
				Expect(isInArchive("TestModule/test-module-1-content", mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("TestModule1/test-module-1-content", mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/MANIFEST.MF", mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/mtad.yaml", mtaArchiveLocation)).To(BeTrue())
			})
		})

		Context("With deployment descriptor which contains only valid resources", func() {
			It("Should build the MTA Archive containing the valid resources", func() {
				resourceContent := filepath.Join(tempDirLocation, "test-resource-1-content")
				os.Create(resourceContent)
				os.WriteFile(resourceContent, []byte("this is a test resource content"), os.ModePerm)
				descriptor := util.MtaDeploymentDescriptor{SchemaVersion: "100", ID: "test", Resources: []util.Resource{
					util.Resource{Name: "TestResource", Parameters: map[string]interface{}{"path": "test-resource-1-content"}},
				}}
				generatedYamlBytes, _ := yaml.Marshal(descriptor)
				testDeploymentDescriptor := tempDirLocation + string(os.PathSeparator) + "mtad.yaml"
				os.WriteFile(testDeploymentDescriptor, generatedYamlBytes, os.ModePerm)
				mtaArchiveLocation, err := util.NewMtaArchiveBuilder([]string{}, []string{"TestResource"}).Build(tempDirLocation)
				Expect(err).To(BeNil())
				_, err = os.Stat(mtaArchiveLocation)
				Expect(err).To(BeNil())
				expectedResourcePath := filepathUnixJoin("TestResource", "test-resource-1-content")
				Expect(isInArchive(expectedResourcePath, mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/MANIFEST.MF", mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/mtad.yaml", mtaArchiveLocation)).To(BeTrue())
				Expect(isManifestValid("META-INF/MANIFEST.MF", map[string]string{"MTA-Resource": "TestResource", "Name": expectedResourcePath}, mtaArchiveLocation)).To(Equal(map[string]string{"MTA-Resource": "TestResource", "Name": expectedResourcePath}))
				defer os.Remove(mtaArchiveLocation)
			})

			It("Should build the MTA Archive containing the valid resources with non-normalized paths", func() {
				resourceContent := filepath.Join(tempDirLocation, "test-resource-1-content")
				os.Create(resourceContent)
				os.WriteFile(resourceContent, []byte("this is a test resource content"), os.ModePerm)
				descriptor := util.MtaDeploymentDescriptor{SchemaVersion: "100", ID: "test", Resources: []util.Resource{
					util.Resource{Name: "TestResource", Parameters: map[string]interface{}{"path": "../test-resource-1-content"}},
				}}
				generatedYamlBytes, _ := yaml.Marshal(descriptor)
				mtadDirectory := filepath.Join(tempDirLocation, "test")
				os.MkdirAll(mtadDirectory, os.ModePerm)
				testDeploymentDescriptor := filepath.Join(mtadDirectory, "mtad.yaml")
				os.WriteFile(testDeploymentDescriptor, generatedYamlBytes, os.ModePerm)
				mtaArchiveLocation, err := util.NewMtaArchiveBuilder([]string{}, []string{"TestResource"}).Build(mtadDirectory)
				Expect(err).To(BeNil())
				_, err = os.Stat(mtaArchiveLocation)
				Expect(err).To(BeNil())
				expectedResourcePath := filepathUnixJoin("TestResource", "test-resource-1-content")
				Expect(isInArchive(expectedResourcePath, mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/MANIFEST.MF", mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/mtad.yaml", mtaArchiveLocation)).To(BeTrue())
				Expect(isManifestValid("META-INF/MANIFEST.MF", map[string]string{"MTA-Resource": "TestResource", "Name": expectedResourcePath}, mtaArchiveLocation)).To(Equal(map[string]string{"MTA-Resource": "TestResource", "Name": expectedResourcePath}))
				defer os.Remove(mtaArchiveLocation)
			})

			It("Should build the MTA Archive containing the resources and add them in the MANIFEST.MF only", func() {
				descriptor := util.MtaDeploymentDescriptor{SchemaVersion: "100", ID: "test", Resources: []util.Resource{
					util.Resource{Name: "TestResource"},
				}}
				generatedYamlBytes, _ := yaml.Marshal(descriptor)
				testDeploymentDescriptor := tempDirLocation + string(os.PathSeparator) + "mtad.yaml"
				os.WriteFile(testDeploymentDescriptor, generatedYamlBytes, os.ModePerm)
				mtaArchiveLocation, err := util.NewMtaArchiveBuilder([]string{}, []string{"TestResource"}).Build(tempDirLocation)
				Expect(err).To(BeNil())
				_, err = os.Stat(mtaArchiveLocation)
				Expect(err).To(BeNil())
				Expect(isInArchive("test-resource-1-content", mtaArchiveLocation)).To(BeFalse())
				Expect(isInArchive("META-INF/MANIFEST.MF", mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/mtad.yaml", mtaArchiveLocation)).To(BeTrue())
				Expect(isManifestValid("META-INF/MANIFEST.MF", map[string]string{"MTA-Resource": "TestResource"}, mtaArchiveLocation)).To(Equal(map[string]string{}))
				defer os.Remove(mtaArchiveLocation)
			})

		})

		Context("With deployment descriptor which contains only valid modules with required dependencies", func() {
			It("Should build the MTA Archive containing the valid modules and required dependencies configuration", func() {
				requiredDependencyContent := filepath.Join(tempDirLocation, "test-required-dep-1-content")
				os.Create(requiredDependencyContent)
				os.WriteFile(requiredDependencyContent, []byte("this is a test module content"), os.ModePerm)
				descriptor := util.MtaDeploymentDescriptor{SchemaVersion: "100", ID: "test", Modules: []util.Module{
					util.Module{Name: "TestModule", RequiredDependencies: []util.RequiredDependency{
						util.RequiredDependency{
							Name: "TestRequired",
							Parameters: map[string]interface{}{
								"path": "test-required-dep-1-content",
							},
						},
					}},
				}}
				generatedYamlBytes, _ := yaml.Marshal(descriptor)
				testDeploymentDescriptor := tempDirLocation + string(os.PathSeparator) + "mtad.yaml"
				os.WriteFile(testDeploymentDescriptor, generatedYamlBytes, os.ModePerm)
				mtaArchiveLocation, err := util.NewMtaArchiveBuilder([]string{"TestModule"}, []string{}).Build(tempDirLocation)
				defer os.Remove(mtaArchiveLocation)
				Expect(err).To(BeNil())
				_, err = os.Stat(mtaArchiveLocation)
				Expect(err).To(BeNil())
				expectedRequiredDependenciesPath := filepathUnixJoin("TestModule", "TestRequired", "test-required-dep-1-content")
				Expect(isInArchive(expectedRequiredDependenciesPath, mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/MANIFEST.MF", mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/mtad.yaml", mtaArchiveLocation)).To(BeTrue())
				Expect(isManifestValid("META-INF/MANIFEST.MF", map[string]string{"MTA-Requires": "TestModule/TestRequired", "Name": expectedRequiredDependenciesPath}, mtaArchiveLocation)).To(Equal(map[string]string{"MTA-Requires": "TestModule/TestRequired", "Name": expectedRequiredDependenciesPath}))
			})

			It("Should build the MTA Archive containing the valid modules and required dependencies configuration with non-normalized paths", func() {
				requiredDependencyContent := filepath.Join(tempDirLocation, "test-required-dep-1-content")
				os.Create(requiredDependencyContent)
				os.WriteFile(requiredDependencyContent, []byte("this is a test module content"), os.ModePerm)
				descriptor := util.MtaDeploymentDescriptor{SchemaVersion: "100", ID: "test", Modules: []util.Module{
					util.Module{Name: "TestModule", RequiredDependencies: []util.RequiredDependency{
						util.RequiredDependency{
							Name: "TestRequired",
							Parameters: map[string]interface{}{
								"path": "../test-required-dep-1-content",
							},
						},
					}},
				}}
				generatedYamlBytes, _ := yaml.Marshal(descriptor)
				mtadDirectory := filepath.Join(tempDirLocation, "test")
				os.MkdirAll(mtadDirectory, os.ModePerm)
				testDeploymentDescriptor := filepath.Join(mtadDirectory, "mtad.yaml")
				os.WriteFile(testDeploymentDescriptor, generatedYamlBytes, os.ModePerm)
				mtaArchiveLocation, err := util.NewMtaArchiveBuilder([]string{"TestModule"}, []string{}).Build(mtadDirectory)
				defer os.Remove(mtaArchiveLocation)
				Expect(err).To(BeNil())
				_, err = os.Stat(mtaArchiveLocation)
				Expect(err).To(BeNil())
				expectedRequiredDependenciesPath := filepathUnixJoin("TestModule", "TestRequired", "test-required-dep-1-content")
				Expect(isInArchive(expectedRequiredDependenciesPath, mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/MANIFEST.MF", mtaArchiveLocation)).To(BeTrue())
				Expect(isInArchive("META-INF/mtad.yaml", mtaArchiveLocation)).To(BeTrue())
				Expect(isManifestValid("META-INF/MANIFEST.MF", map[string]string{"MTA-Requires": "TestModule/TestRequired", "Name": expectedRequiredDependenciesPath}, mtaArchiveLocation)).To(Equal(map[string]string{"MTA-Requires": "TestModule/TestRequired", "Name": expectedRequiredDependenciesPath}))
			})
		})

		AfterEach(func() {
			os.RemoveAll(tempDirLocation)
		})
	})
})

func isInArchive(fileName, archiveLocation string) bool {
	mtaArchiveReader, err := zip.OpenReader(archiveLocation)
	if err != nil {
		return false
	}
	defer mtaArchiveReader.Close()
	for _, file := range mtaArchiveReader.File {
		if file.Name == fileName {
			return true
		}
	}
	return false
}

func isManifestValid(manifestLocation string, searchCriteria map[string]string, archiveLocation string) map[string]string {
	mtaArchiveReader, err := zip.OpenReader(archiveLocation)
	if err != nil {
		return map[string]string{}
	}
	defer mtaArchiveReader.Close()
	searchCriteriaResult := make(map[string]string)
	for _, file := range mtaArchiveReader.File {
		if file.Name == manifestLocation {
			reader, err := file.Open()
			if err != nil {
				return map[string]string{}
			}
			defer reader.Close()
			manifestBytes, _ := io.ReadAll(reader)
			manifestSplittedByNewLine := strings.Split(string(manifestBytes), "\n")
			for _, manifestSectionElement := range manifestSplittedByNewLine {
				if strings.Trim(manifestSectionElement, " ") == "" {
					continue
				}
				separatorIndex := strings.Index(manifestSectionElement, ":")
				key := manifestSectionElement[:separatorIndex]
				value := manifestSectionElement[separatorIndex+1:]
				if searchCriteria[key] != "" {
					delete(searchCriteria, key)
					searchCriteriaResult[key] = strings.Trim(value, " ")
				}
			}
			break
		}
	}
	return searchCriteriaResult
}

func filepathUnixJoin(elements ...string) string {
	return strings.Join(elements, "/")
}
