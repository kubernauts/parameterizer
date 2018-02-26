package parameterizer

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Parse parses a Parameterizer resource from a YAML file.
func Parse(filename string) (p Resource, err error) {
	p = Resource{}
	pfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return p, err
	}
	err = yaml.Unmarshal(pfile, &p)
	if err != nil {
		return p, err
	}
	return p, nil
}
