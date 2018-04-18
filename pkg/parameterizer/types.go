package parameterizer

import (
	"bytes"
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Parameterizer represents the `Parameterizer` resource.
type Parameterizer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PSpec `json:"spec"`
}

// PSpec parameterizer Spec
type PSpec struct {
	Resources       []ResourceSpec       `json:"resources"`
	UserInputs      []ResourceSpec       `json:"userInputs,omitempty"`
	Volumes         []v1.Volume          `json:"volumes,omitempty"`
	Transformations []TransformationSpec `json:"transformations"`
}

// SourceSpec represents resource location
type SourceSpec struct {
	Container v1.Container `json:"container,omitempty"`
	Files     []struct {
		Dest    string `json:"dest"`
		Content string `json:"content"`
	} `json:"files,omitempty"`
	Fetch struct {
		URLs []string `json:"urls,omitempty"`
		Dest string   `json:"dest"`
	} `json:"fetch,omitempty"`
}

// ResourceSpec represents the `resources` sub-resource.
type ResourceSpec struct {
	Name         string           `json:"name"`
	Source       SourceSpec       `json:"source"`
	VolumeMounts []v1.VolumeMount `json:"volumeMounts"`
}

// UserInputSpec represents the `userInputs` sub-resource.
type UserInputSpec struct {
	ResourceSpec `json:",inline"`
}

// TransformationSpec represents the `apply` sub-resource.
type TransformationSpec struct {
	Name      string                 `json:"name"`
	Container v1.Container           `json:"container,omitempty"`
	Helm      HelmTransformationSpec `json:"helm,omitempty"`
}

type NamedValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HelmRepo struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type HelmChart struct {
	Name    string   `json:"name,omitempty"`
	Path    string   `json:"path,omitempty"`
	Version string   `json:"version,omitempty"`
	Repo    HelmRepo `json:"repo,omitempty"`
}

type HelmTransformationSpec struct {
	Chart        HelmChart        `json:"chart"`
	ReleaseName  string           `json:"releaseName,omitempty"`
	ValueFiles   []string         `json:"valueFiles,omitempty"`
	SetValues    []NamedValue     `json:"setValues,omitempty"`
	ExtraOpts    []string         `json:"extraOpts,omitempty"`
	VolumeMounts []v1.VolumeMount `json:"volumeMounts,omitempty"`
	OutputFile   string           `json:"outputFile,omitempty"`
}

func (parameterizer Parameterizer) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Parameterizer {\n name: %v\n", parameterizer.ObjectMeta.Name))
	buffer.WriteString(fmt.Sprintf(" resources:\n"))
	for _, r := range parameterizer.Spec.Resources {
		buffer.WriteString(fmt.Sprintf("  - %v\n", r.Name))
		for _, u := range r.Source.Fetch.URLs {
			buffer.WriteString(fmt.Sprintf("    - %v\n", u))
		}
	}
	buffer.WriteString(fmt.Sprintf(" user inputs:\n"))
	for _, u := range parameterizer.Spec.UserInputs {
		buffer.WriteString(fmt.Sprintf("  - %v\n", u.Name))
		for _, url := range u.Source.Fetch.URLs {
			buffer.WriteString(fmt.Sprintf("    - %v\n", url))
		}
	}
	buffer.WriteString(fmt.Sprintf(" volumes:\n"))
	for _, v := range parameterizer.Spec.Volumes {
		buffer.WriteString(fmt.Sprintf("  - %v\n", v.Name))
	}
	buffer.WriteString(fmt.Sprintf(" apply:\n"))
	for _, a := range parameterizer.Spec.Transformations {
		buffer.WriteString(fmt.Sprintf("  - %v:\n", a.Container.Name))
		buffer.WriteString(fmt.Sprintf("    with image [%v] executing commands:\n", a.Container.Image))
		for _, c := range a.Container.Command {
			buffer.WriteString(fmt.Sprintf("    - %v\n", c))
		}
	}
	buffer.WriteString(fmt.Sprintf("}\n"))
	return buffer.String()
}
