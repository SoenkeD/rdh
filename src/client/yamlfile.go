package client

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SoenkeD/rdh/src/core"
	"gopkg.in/yaml.v3"
)

type ParseSpecFileFn func(specFile SpecFile) (core.CreationSpecDefinition, error)

type YamlFile struct {
	SourceDir   string
	ParseFileFn map[string]ParseSpecFileFn
}

func NewYamlFile(sourceDir string, parseFileFn map[string]ParseSpecFileFn) *YamlFile {
	return &YamlFile{
		SourceDir:   sourceDir,
		ParseFileFn: parseFileFn,
	}
}

func (yf *YamlFile) LoadSpecs(rdhBe *core.ResourceDefinitionHandler) error {

	files, err := ReadFilesInDir(yf.SourceDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		specFile, err := ReadSpecFile(file)
		if err != nil {
			return err
		}

		entry, ok := yf.ParseFileFn[specFile.Kind]
		if !ok {
			return fmt.Errorf("resource kind %s in %s is not supported", specFile.Kind, specFile.Id)
		}

		var specs core.CreationSpecDefinition
		specs, err = entry(specFile)

		if err != nil {
			return err
		}

		err = rdhBe.CreateSpec(specFile.Id, specFile.Kind, specs)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadFilesInDir(dir string) ([]string, error) {
	var files []string

	fInfo, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, fInf := range fInfo {
		if !fInf.IsDir() {
			files = append(files, filepath.Join(dir, fInf.Name()))
		}
	}

	return files, nil
}

type SpecFile struct {
	Id    string
	Kind  string
	Specs map[string]any
}

func ParseSpecFile[T any](specFile SpecFile) (resSpec T, err error) {

	contentYamlBytes, err := yaml.Marshal(specFile.Specs)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(contentYamlBytes, &resSpec)
	if err != nil {
		return
	}

	return
}

func ReadSpecFile(location string) (specFile SpecFile, err error) {
	fBytes, err := os.ReadFile(location)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(fBytes, &specFile)
	if err != nil {
		return
	}

	return
}
