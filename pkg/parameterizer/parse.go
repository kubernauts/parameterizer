package parameterizer

import (
	"io/ioutil"

	yaml "github.com/ghodss/yaml"
)

// Parse parses a Parameterizer resource from a YAML file.
func Parse(filename string) (p Parameterizer, err error) {
	p = Parameterizer{}
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
