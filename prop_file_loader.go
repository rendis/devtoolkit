package devtoolkit

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

// configFileType represents the supported configuration file formats.
type configFileType int

const (
	ymlType  configFileType = iota // YAML file type
	jsonType                       // JSON file type
)

// LoadPropFile loads configuration properties from a file into the provided
// slice of structs. The file format can be either YAML or JSON.
// The 'filePath' parameter specifies the path to the configuration file.
// The 'props' parameter is a slice of pointers to struct instances that
// should be populated with the loaded properties.
// Returns an error if the file cannot be loaded, parsed, or is of an unsupported format.
func LoadPropFile[T any](filePath string, props []T) error {
	// get the configuration file type (yml or json).
	fileType, err := getConfigFileType(filePath)
	if err != nil {
		return fmt.Errorf("error getting config file type of file '%s': %w", filePath, err)
	}

	// read the configuration file.
	propArr, err := readPropFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading property file '%s': %w", filePath, err)
	}

	// select the appropriate parsing function based on the file type.
	var parseFn func([]byte, interface{}) error
	switch fileType {
	case ymlType:
		parseFn = parseFromYml
	case jsonType:
		parseFn = parseFromJson
	default:
		return fmt.Errorf("invalid config file '%s' type. only 'yml' and 'json' are supported", filePath)
	}

	// parse the configuration file and populate the provided structs.
	var parseErr error
	for _, prop := range props {
		if err := parseFn(propArr, prop); err != nil {
			if parseErr == nil {
				parseErr = err
			} else {
				parseErr = fmt.Errorf("%v; %w", parseErr, err)
			}
		}
	}

	return parseErr
}

// readPropFile reads a file from the provided 'filePath' and returns its contents
// as a byte slice. Environment variables within the file are expanded.
// Returns an error if the file does not exist or cannot be read.
func readPropFile(filePath string) ([]byte, error) {
	// check if the file exists.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file '%s' does not exist: %w", filePath, err)
	}

	// read the file contents.
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %w", filePath, err)
	}

	// expand environment variables in the file contents.
	propWithEnv := []byte(os.ExpandEnv(string(b)))

	return propWithEnv, nil
}

// getConfigFileType determines the file type of the configuration file specified by 'path'.
// It returns ymlType for .yml files and jsonType for .json files.
// An error is returned if the file extension is unsupported.
func getConfigFileType(path string) (configFileType, error) {
	ext := filepath.Ext(path)
	switch ext {
	case ".yml":
		return ymlType, nil
	case ".json":
		return jsonType, nil
	default:
		return 0, errors.New("invalid config file type. only '.yml' and '.json' are supported")
	}
}

// parseFromYml parses the contents of a YAML file represented by 'propArr' into
// the provided struct 'prop'. Returns an error if the parsing fails.
func parseFromYml(propArr []byte, prop interface{}) error {
	if err := yaml.Unmarshal(propArr, prop); err != nil {
		return fmt.Errorf("error parsing YAML file to struct '%v': %v", prop, err)
	}
	return nil
}

// parseFromJson parses the contents of a JSON file represented by 'propArr' into
// the provided struct 'prop'. Logs and exits if the parsing fails.
func parseFromJson(propArr []byte, prop interface{}) error {
	if err := json.Unmarshal(propArr, prop); err != nil {
		log.Fatalf("error parsing JSON file to struct '%v': %v", prop, err)
	}
	return nil
}
