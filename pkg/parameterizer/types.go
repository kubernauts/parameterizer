package parameterizer

// Resource encodes the Parameterizer resource.
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

type presource struct {
	Name   string `yaml:"name"`
	Source struct {
		URLs []string `yaml:"urls"`
	} `yaml:"source"`
	Volume pvolume `yaml:"volume"`
}

type puserinput struct {
	Name   string `yaml:"name"`
	Source struct {
		HostPath struct {
			Path string `yaml:"path"`
		} `yaml:"hostPath,omitempty"`
	} `yaml:"source"`
	Volume pvolume `yaml:"volume"`
}

type pvolume struct {
	Name     string `yaml:"name"`
	HostPath struct {
		Path string `yaml:"path"`
	} `yaml:"hostPath,omitempty"`
	EmptyDir struct {
		Path string `yaml:"path"`
	} `yaml:"emptyDir,omitempty"`
}

type papply struct {
	Image string `yaml:"image"`
}
