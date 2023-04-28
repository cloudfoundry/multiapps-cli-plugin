package util

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
)

const deploymentDescriptorYamlName string = "mtad.yaml"

// MtaArchiveBuilder builds mta archive
type MtaArchiveBuilder struct {
	modules   []string
	resources []string
}

// NewMtaArchiveBuilder constructs new MtaArchiveBuilder
func NewMtaArchiveBuilder(modules, resources []string) MtaArchiveBuilder {
	return MtaArchiveBuilder{
		modules:   modules,
		resources: resources,
	}
}

// Build creates deployment archive from the provided deployment descriptor
func (builder MtaArchiveBuilder) Build(deploymentDescriptorLocation string) (string, error) {
	descriptor, deploymentDescriptorFile, err := ParseDeploymentDescriptor(deploymentDescriptorLocation)
	if err != nil {
		return "", err
	}

	modulesPaths, err := builder.getModulesPaths(descriptor.Modules)
	if err != nil {
		return "", fmt.Errorf("Error building MTA Archive: %s", err.Error())
	}
	resourcesPaths, err := builder.getResourcesPaths(descriptor.Resources)
	if err != nil {
		return "", fmt.Errorf("Error building MTA Archive: %s", err.Error())
	}
	bindingParametersPaths := builder.getBindingParametersPaths(descriptor.Modules)

	modulesSections := buildSection(normalizePathsUsingUnixSeparator(modulesPaths), MtaModule)
	resourcesSections := buildSection(normalizePathsUsingUnixSeparator(resourcesPaths), MtaResource)
	bindingParametersSections := buildSection(normalizePathsUsingUnixSeparator(bindingParametersPaths), MtaRequires)

	manifestBuilder := NewMtaManifestBuilder()
	manifestBuilder.ManifestSections(modulesSections)
	manifestBuilder.ManifestSections(resourcesSections)
	manifestBuilder.ManifestSections(bindingParametersSections)
	manifestBuilder.ManifestSections([]ManifestSection{NewMtaManifestSectionBuilder().Name(MtadAttribute).Build()})

	manifestLocation, err := manifestBuilder.Build()
	if err != nil {
		return "", err
	}
	defer os.Remove(manifestLocation)

	mtaAssembly, err := os.MkdirTemp("", "mta-assembly")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(mtaAssembly)

	metaInfLocation := filepath.Join(mtaAssembly, "META-INF")
	err = os.Mkdir(metaInfLocation, os.ModePerm)
	if err != nil {
		return "", err
	}

	manifestInfo, err := os.Stat(manifestLocation)
	if err != nil {
		return "", err
	}

	err = copyFile(manifestLocation, filepath.Join(metaInfLocation, manifestInfo.Name()))
	if err != nil {
		return "", err
	}

	err = copyFile(deploymentDescriptorFile, filepath.Join(metaInfLocation, "mtad.yaml"))
	if err != nil {
		return "", err
	}
	// TODO: modify the deployment descriptor after copying it in order not to contain any path parameters...

	err = copyContent(deploymentDescriptorLocation, modulesPaths, mtaAssembly)
	if err != nil {
		return "", err
	}

	err = copyContent(deploymentDescriptorLocation, resourcesPaths, mtaAssembly)
	if err != nil {
		return "", err
	}

	err = copyContent(deploymentDescriptorLocation, bindingParametersPaths, mtaAssembly)
	if err != nil {
		return "", err
	}

	mtaArchiveName := descriptor.ID + ".mtar"
	mtaArchiveLocation := filepath.Join(deploymentDescriptorLocation, mtaArchiveName)
	err = CreateMtaArchive(mtaAssembly, mtaArchiveLocation)
	if err != nil {
		return "", err
	}

	mtaArchiveAbsolutePath, err := filepath.Abs(mtaArchiveLocation)
	if err != nil {
		return "", err
	}

	return mtaArchiveAbsolutePath, nil
}

func copyContent(sourceDirectory string, elementsPaths map[string]string, targetLocation string) error {
	for name, path := range elementsPaths {
		if path != "" {
			sourceLocation := filepath.Join(sourceDirectory, path)
			filesInSourceInfo, err := os.Stat(sourceLocation)
			if err != nil {
				return fmt.Errorf("Error building MTA Archive: file path %s not found", sourceLocation)
			}
			destinationLocation := filepath.Join(targetLocation, name, filepath.Base(path))
			if filesInSourceInfo.IsDir() {
				err = copyDirectory(sourceLocation, destinationLocation)
			} else {
				os.MkdirAll(filepath.Dir(destinationLocation), os.ModePerm)
				err = copyFile(sourceLocation, destinationLocation)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func copyDirectory(src, dest string) error {
	var err error
	var filesInDestinationInfo []os.DirEntry
	var sourceInfo os.FileInfo

	if sourceInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dest, sourceInfo.Mode()); err != nil {
		return err
	}

	if filesInDestinationInfo, err = os.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range filesInDestinationInfo {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dest, fd.Name())

		if fd.IsDir() {
			if err = copyDirectory(srcfp, dstfp); err != nil {
				return err
			}
		} else {
			if err = copyFile(srcfp, dstfp); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dest string) error {
	fileFrom, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fileFrom.Close()

	toFile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer toFile.Close()

	_, err = io.Copy(toFile, fileFrom)
	if err != nil {
		return err
	}

	return nil
}

func buildSection(elements map[string]string, locatorName string) []ManifestSection {
	result := make([]ManifestSection, 0)
	for key, value := range concatenateElementsWithSameValue(elements) {
		manifestSectionBuilder := NewMtaManifestSectionBuilder()
		manifestSectionBuilder.Name(key)
		manifestSectionBuilder.Attribute(locatorName, strings.Join(value, ","))
		result = append(result, manifestSectionBuilder.Build())
	}
	return result
}

func concatenateElementsWithSameValue(elements map[string]string) map[string][]string {
	result := make(map[string][]string)
	for key, value := range elements {
		if value == "" {
			continue
		}
		if len(result[value]) != 0 {
			result[value] = append(result[value], key)
		} else {
			result[value] = []string{key}
		}
	}
	return result
}

func (builder MtaArchiveBuilder) getBindingParametersPaths(deploymentDescriptorResources []Module) map[string]string {
	result := make(map[string]string, 0)
	modulesToAdd := filterModules(deploymentDescriptorResources, func(module Module) bool {
		return Contains(builder.modules, module.Name)
	})
	for _, module := range modulesToAdd {
		requiredDependenciesConfigPaths := getRequiredDependenciesConfigPaths(module.RequiredDependencies)
		for requiredDependencyName, configFile := range requiredDependenciesConfigPaths {
			result[module.Name+"/"+requiredDependencyName] = configFile
		}
	}
	return result
}

func getRequiredDependenciesConfigPaths(requiredDependencies []RequiredDependency) map[string]string {
	result := make(map[string]string, 0)
	for _, requiredDependency := range requiredDependencies {
		if requiredDependency.Parameters["path"] != nil {
			result[requiredDependency.Name] = getString(requiredDependency.Parameters["path"])
		}
	}
	return result
}

func (builder MtaArchiveBuilder) getResourcesPaths(deploymentDescriptorResources []Resource) (map[string]string, error) {
	err := validateSpecifiedResources(builder.resources, deploymentDescriptorResources)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	resourcesToAdd := filterResources(deploymentDescriptorResources, func(resource Resource) bool {
		return Contains(builder.resources, resource.Name)
	})

	for _, resource := range resourcesToAdd {
		result[resource.Name] = getString(resource.Parameters["path"])
	}
	return result, nil
}

func validateSpecifiedResources(resourcesForDeployment []string, deploymentDescriptorResources []Resource) error {
	deploymentDescriptorResourceNames := make([]string, 0)
	for _, deploymentDescriptorResource := range deploymentDescriptorResources {
		deploymentDescriptorResourceNames = append(deploymentDescriptorResourceNames, deploymentDescriptorResource.Name)
	}
	specifiedResourcesNotPartOfDeploymentDescriptor := make([]string, 0)
	for _, resourceForDeployment := range resourcesForDeployment {
		if !Contains(deploymentDescriptorResourceNames, resourceForDeployment) {
			specifiedResourcesNotPartOfDeploymentDescriptor = append(specifiedResourcesNotPartOfDeploymentDescriptor, resourceForDeployment)
		}
	}

	if len(specifiedResourcesNotPartOfDeploymentDescriptor) == 0 {
		return nil
	}

	return fmt.Errorf("Resources %s are specified for deployment but are not part of deployment descriptor resources", strings.Join(specifiedResourcesNotPartOfDeploymentDescriptor, ", "))
}

func getString(value interface{}) string {
	if value == nil {
		return ""
	}
	return value.(string)
}

func (builder MtaArchiveBuilder) getModulesPaths(deploymentDescriptorResources []Module) (map[string]string, error) {
	err := validateSpecifiedModules(builder.modules, deploymentDescriptorResources)
	if err != nil {
		return nil, err
	}
	modulesToAdd := filterModules(deploymentDescriptorResources, func(module Module) bool {
		return Contains(builder.modules, module.Name)
	})
	specifiedModulesWithoutPaths := filterModules(modulesToAdd, func(module Module) bool {
		return module.Path == ""
	})
	moduleNamesWithoutPaths := make([]string, 0)
	for _, moduleWithoutPath := range specifiedModulesWithoutPaths {
		moduleNamesWithoutPaths = append(moduleNamesWithoutPaths, moduleWithoutPath.Name)
	}
	if len(moduleNamesWithoutPaths) > 0 {
		ui.Warn("Modules %s do not have a path, specified for their binaries and will be ignored", strings.Join(moduleNamesWithoutPaths, ", "))
	}
	result := make(map[string]string)
	for _, moduleToAdd := range modulesToAdd {
		if moduleToAdd.Path != "" {
			result[moduleToAdd.Name] = moduleToAdd.Path
		}
	}

	return result, nil
}

func normalizePathsUsingUnixSeparator(elementsPaths map[string]string) map[string]string {
	normalizedPaths := make(map[string]string, len(elementsPaths))
	for key, value := range elementsPaths {
		if value != "" {
			zipEntryNormalizedPath := fmt.Sprintf("%s/%s", key, filepath.Base(value))
			normalizedPaths[key] = zipEntryNormalizedPath
		}
	}
	return normalizedPaths
}

func validateSpecifiedModules(modulesForDeployment []string, deploymentDescriptorResources []Module) error {
	deploymentDescriptorResourceNames := make([]string, 0)
	for _, deploymentDescriptorResource := range deploymentDescriptorResources {
		deploymentDescriptorResourceNames = append(deploymentDescriptorResourceNames, deploymentDescriptorResource.Name)
	}
	specifiedResourcesNotPartOfDeploymentDescriptor := make([]string, 0)
	for _, moduleForDeployment := range modulesForDeployment {
		if !Contains(deploymentDescriptorResourceNames, moduleForDeployment) {
			specifiedResourcesNotPartOfDeploymentDescriptor = append(specifiedResourcesNotPartOfDeploymentDescriptor, moduleForDeployment)
		}
	}

	if len(specifiedResourcesNotPartOfDeploymentDescriptor) == 0 {
		return nil
	}

	return fmt.Errorf("Modules %s are specified for deployment but are not part of deployment descriptor modules", strings.Join(specifiedResourcesNotPartOfDeploymentDescriptor, ", "))
}

func filterModules(modulesSlice []Module, prediacte func(element Module) bool) []Module {
	result := make([]Module, 0)
	for _, sliceElement := range modulesSlice {
		if prediacte(sliceElement) {
			result = append(result, sliceElement)
		}
	}
	return result
}

func filterResources(resourcesSlice []Resource, predicate func(element Resource) bool) []Resource {
	result := make([]Resource, 0)
	for _, sliceElement := range resourcesSlice {
		if predicate(sliceElement) {
			result = append(result, sliceElement)
		}
	}
	return result
}

func Contains(slice []string, element string) bool {
	for _, sliceElement := range slice {
		if sliceElement == element {
			return true
		}
	}
	return false
}
