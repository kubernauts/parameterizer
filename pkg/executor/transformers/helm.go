package transformers

import (
	"fmt"
	"github.com/kubernauts/parameterizer/pkg/parameterizer"
	"k8s.io/api/core/v1"
	"strings"
)

func HelmTransform(t *parameterizer.TransformationSpec) v1.Container {
	args := []string{"helm", "template"}

	if t.Helm.ValueFile != "" {
		args = append(args, []string{"-f", t.Helm.ValueFile}...)
	}

	if len(t.Helm.ExtraOpts) > 0 {
		for _, opt := range t.Helm.ExtraOpts {
			args = append(args, opt)
		}
	}

	setArgs := []string{}
	for k, v := range t.Helm.SetArgs {
		setArgs = append(setArgs, fmt.Sprintf("%s=%s", k, v))
	}

	if len(setArgs) > 0 {
		args = append(args, "--set "+strings.Join(setArgs, ","))
	}

	container := v1.Container{
		Image:        "quay.io/wire/alpine-helm",
		Command:      args,
		Name:         t.Name,
		VolumeMounts: t.Helm.VolumeMounts,
	}
	return container
}
