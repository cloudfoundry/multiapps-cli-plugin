package util

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
)

const ContentTypeAttribute string = "Content-Type"
const ManifestVersion string = "Manifest-Version"
const Name string = "Name"
const ManifestName string = "MANIFEST.MF"
const MtaResource string = "MTA-Resource"
const MtaRequires string = "MTA-Requires"
const MtaModule string = "MTA-Module"
const MtadAttribute string = "META-INF/mtad.yaml"
const SectionSeparator string = "\n"

type MtaManifest struct {
	ManifestVersion  string
	ManifestSections []ManifestSection
}

type ManifestSection struct {
	Name       string
	Attributes map[string]string
}

type MtaManifestBuilder struct {
	manifest MtaManifest
}

func NewMtaManifestBuilder() *MtaManifestBuilder {
	return &MtaManifestBuilder{}
}

func (builder *MtaManifestBuilder) ManifestSections(sections []ManifestSection) *MtaManifestBuilder {
	builder.manifest.ManifestSections = append(builder.manifest.ManifestSections, sections...)
	return builder
}

func (builder *MtaManifestBuilder) Build() (string, error) {
	location, err := ioutil.TempDir("", "mta-manifest")
	if err != nil {
		return "", err
	}
	fileLocation := filepath.Join(location, ManifestName)
	file, err := os.Create(fileLocation)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileWriter := bufio.NewWriter(file)
	_, err = fileWriter.WriteString(getManifestVersion(builder.manifest))
	if err != nil {
		return "", err
	}
	_, err = fileWriter.WriteString(SectionSeparator)
	if err != nil {
		return "", err
	}
	_, err = fileWriter.WriteString(SectionSeparator)
	if err != nil {
		return "", err
	}

	for _, section := range builder.manifest.ManifestSections {
		_, err = fileWriter.WriteString(section.Name)
		if err != nil {
			return "", err
		}
		_, err = fileWriter.WriteString(SectionSeparator)
		if err != nil {
			return "", err
		}
		err = writeSectionAttributes(fileWriter, section.Attributes)
		if err != nil {
			return "", err
		}
		fileWriter.WriteString(SectionSeparator)
		fileWriter.WriteString(SectionSeparator)
	}

	fileWriter.Flush()
	return fileLocation, nil
}

func getManifestVersion(manifest MtaManifest) string {
	if manifest.ManifestVersion == "" {
		return ManifestVersion + ": 1.0"
	}
	return ManifestVersion + ": " + manifest.ManifestVersion
}

func writeSectionAttributes(fileWriter *bufio.Writer, attributes map[string]string) error {
	for attrName, attrValue := range attributes {
		_, err := fileWriter.WriteString(attrName + ": " + attrValue)
		if err != nil {
			return err
		}
		_, err = fileWriter.WriteString(SectionSeparator)
		if err != nil {
			return err
		}
	}
	return nil
}

type MtaManifestSectionBuilder struct {
	section ManifestSection
}

func NewMtaManifestSectionBuilder() *MtaManifestSectionBuilder {
	return &MtaManifestSectionBuilder{
		section: ManifestSection{
			Attributes: make(map[string]string),
		},
	}
}

func (builder *MtaManifestSectionBuilder) Name(name string) *MtaManifestSectionBuilder {
	builder.section.Name = Name + ": " + name
	return builder
}

func (builder *MtaManifestSectionBuilder) Attribute(attributeName, attributeValue string) *MtaManifestSectionBuilder {
	builder.section.Attributes[attributeName] = attributeValue
	return builder
}

func (builder *MtaManifestSectionBuilder) Build() ManifestSection {
	return builder.section
}
