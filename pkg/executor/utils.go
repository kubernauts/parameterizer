package executor

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"os"
)

// PrintObj marshall an object and print it to stdout
func PrintObj(v interface{}, outputFormat string) {
	output, err := MarshallObj(v, outputFormat)
	if err != nil {
		fmt.Println("error:", err)
	}

	if outputFormat == "yaml" {
		os.Stdout.Write([]byte("---\n"))
	}
	os.Stdout.Write(output)
}

// MarshallObj marshall a struct to json or yaml
func MarshallObj(v interface{}, outputFormat string) ([]byte, error) {
	jsonBytes, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return nil, err
	}
	if outputFormat == "json" {
		return jsonBytes, nil
	}
	yamlBytes, err := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		return nil, err
	}
	return yamlBytes, nil
}
