package transformers

import (
	"fmt"
	"github.com/kubernauts/parameterizer/pkg/parameterizer"
	"k8s.io/api/core/v1"
	"strings"
)

func helmFetch(t *parameterizer.HelmChart, dest string) ([]string, []string) {
	var addRepo []string
	var fetchChart []string
	if dest == "" {
		dest = "."
	}
	if t.Repo.Name != "" {
		addRepo = []string{"helm", "repo", "add", t.Repo.Name, t.Repo.URL}
		fetchChart = []string{"helm", "fetch", "--untar", "--untardir",
			dest + "/" + strings.Split(t.Name, "/")[0], t.Name}
		if t.Version != "" {
			fetchChart = append(fetchChart, []string{"--version", t.Version}...)
		}

	}
	return addRepo, fetchChart

}

func helmTemplate(t *parameterizer.TransformationSpec, chartPath string) []string {
	args := []string{"helm", "template", chartPath}

	if len(t.Helm.ValueFiles) > 0 {
		args = append(args, []string{"-f", strings.Join(t.Helm.ValueFiles, ",")}...)
	}

	if len(t.Helm.ExtraOpts) > 0 {
		for _, opt := range t.Helm.ExtraOpts {
			args = append(args, opt)
		}
	}

	setArgs := []string{}
	for _, v := range t.Helm.SetValues {
		setArgs = append(setArgs, fmt.Sprintf("%s=%s", v.Name, v.Value))
	}

	if len(setArgs) > 0 {
		args = append(args, "--set "+strings.Join(setArgs, ","))
	}

	if t.Helm.ReleaseName != "" {
		args = append(args, []string{"-n", t.Helm.ReleaseName}...)
	}
	return args
}

func shellScript(cmds [][]string) string {
	command := []string{}
	for _, cmd := range cmds {
		command = append(command, strings.Join(cmd, " "))
	}
	return strings.Join(command, " && ")
}

func HelmTransform(t *parameterizer.TransformationSpec) v1.Container {
	addRepoCmd, fetchCmd := helmFetch(&t.Helm.Chart, ".")
	tplCmd := helmTemplate(t, t.Helm.Chart.Name)
	cmd := []string{"sh", "-c"}
	if addRepoCmd != nil && fetchCmd != nil {
		cmd = append(cmd, shellScript([][]string{addRepoCmd, fetchCmd, tplCmd}))
	} else {
		cmd = append(cmd, strings.Join(tplCmd, " "))
	}

	container := v1.Container{
		Image:        "quay.io/wire/alpine-helm",
		Command:      cmd,
		Name:         t.Name,
		VolumeMounts: t.Helm.VolumeMounts,
	}
	return container
}
