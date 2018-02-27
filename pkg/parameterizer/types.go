package parameterizer

import (
	"bytes"
	"fmt"
)

// Resource represents the `Parameterizer` resource.
type Resource struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name   string            `yaml:"name"`
		Labels map[string]string `yaml:"labels"`
	} `yaml:"metadata"`
	Spec struct {
		Resources  []presource  `yaml:"resources"`
		UserInputs []puserinput `yaml:"userInputs"`
		Volumes    []pvolume    `yaml:"volumes"`
		Apply      []papply     `yaml:"apply"`
	} `yaml:"spec"`
}

// presource represents the `resources` sub-resource.
type presource struct {
	Name   string `yaml:"name"`
	Source struct {
		URLs []string `yaml:"urls"`
	} `yaml:"source"`
	Volume pvolume `yaml:"volume"`
}

// puserinput represents the `userInputs` sub-resource.
type puserinput struct {
	Name   string `yaml:"name"`
	Source struct {
		HostPath struct {
			Path string `yaml:"path"`
		} `yaml:"hostPath,omitempty"`
	} `yaml:"source"`
	Volume pvolume `yaml:"volume"`
}

// pvolume represents the `volume` sub-resource.
type pvolume struct {
	Name     string `yaml:"name"`
	HostPath struct {
		Path string `yaml:"path"`
	} `yaml:"hostPath,omitempty"`
	EmptyDir struct {
		Path string `yaml:"path"`
	} `yaml:"emptyDir,omitempty"`
}

// papply represents the `apply` sub-resource.
type papply struct {
	Image        string   `yaml:"image"`
	Commands     []string `yaml:"commands"`
	VolumeMounts []struct {
		Name      string `yaml:"name"`
		mountPath string `yaml:"mountPath"`
	} `yaml:"volumeMounts"`
}

func (parameterizer Resource) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Parameterizer {\n name: %v\n", parameterizer.Metadata.Name))
	buffer.WriteString(fmt.Sprintf(" resources:\n"))
	for _, r := range parameterizer.Spec.Resources {
		buffer.WriteString(fmt.Sprintf("  - %v\n", r.Name))
		for _, u := range r.Source.URLs {
			buffer.WriteString(fmt.Sprintf("    - %v\n", u))
		}
	}
	buffer.WriteString(fmt.Sprintf(" user inputs:\n"))
	for _, u := range parameterizer.Spec.UserInputs {
		buffer.WriteString(fmt.Sprintf("  - %v\n", u.Name))
	}
	buffer.WriteString(fmt.Sprintf(" volumes:\n"))
	for _, v := range parameterizer.Spec.Volumes {
		buffer.WriteString(fmt.Sprintf("  - %v\n", v.Name))
	}
	buffer.WriteString(fmt.Sprintf(" apply:\n"))
	for _, a := range parameterizer.Spec.Apply {
		buffer.WriteString(fmt.Sprintf("  - %v\n", a.Image))
		for _, c := range a.Commands {
			buffer.WriteString(fmt.Sprintf("    - %v\n", c))
		}
	}
	buffer.WriteString(fmt.Sprintf("}\n"))
	return buffer.String()
}
