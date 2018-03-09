package executor

import (
	"bytes"
	"fmt"
	"github.com/kubernauts/parameterizer/pkg/parameterizer"
	"io/ioutil"
	// "k8s.io/api/apps/v1beta1"
	"github.com/kubernauts/parameterizer/pkg/executor/transformers"
	"github.com/pborman/uuid"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"os"
	"os/exec"
	"strings"
	"time"
)

// Run executes the Parameterizer resource's transformation
// steps as defined in the apply sub-resource.
func Run(p parameterizer.Parameterizer) (err error) {
	// we create a temporary manifest file with all
	// the necessary settings in there
	podName, mf, mc, err := createmanifest(p)
	if err != nil {
		return err
	}
	mfn := mf.Name()
	defer func() {
		e := os.Remove(mfn)
		if e != nil {
			fmt.Printf("Couldn't clean up temporary manifest %v", mfn)
		}
	}()
	fmt.Printf("Using manifest:\n%v\n", mc)
	cmd := []string{"create", "-f", mfn}
	fmt.Printf("Executing command: %v\n", strings.Join(cmd, " "))
	res, err := kubectl(true, cmd[0], cmd[1:]...)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", res)
	time.Sleep(1 * time.Minute)
	res, err = kubectl(true, "get", "po", "-a")
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", res)
	res, err = kubectl(true, "logs", podName)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", res)
	res, err = kubectl(true, "logs", strings.Split(res, " ")...)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", res)
	res, err = kubectl(true, "delete", "po", podName, "--force")
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", res)
	return nil
}

func generateName(name string) string {
	return "krm-exec-" + name + "-" + uuid.NewUUID().String()[0:8]
}

func fetchSourceContainer(source parameterizer.SourceSpec) (string, []string) {
	image := ""
	command := []string{}

	if source.Container.Image != "" {
		image = source.Container.Image
		command = source.Container.Command
	} else if len(source.Fetch.URLs) > 0 {
		image = "alpine:3.7"
		command = []string{"sh", "-c", "wget -P " + source.Fetch.Dest + " " + source.Fetch.URLs[0]}
	}

	return image, command
}

func createResourceContainers(name string, resources []parameterizer.ResourceSpec) []v1.Container {
	initContainers := []v1.Container{}

	for _, resource := range resources {
		image, command := fetchSourceContainer(resource.Source)
		container := v1.Container{
			Name:         resource.Name,
			Image:        image,
			Command:      command,
			VolumeMounts: resource.VolumeMounts}
		initContainers = append(initContainers, container)
	}
	return initContainers
}

func createTransformationContainers(name string, transformations []parameterizer.TransformationSpec) []v1.Container {
	initContainers := []v1.Container{}
	container := v1.Container{}
	for _, transformation := range transformations {
		if transformation.Helm.Chart.Name != "" {
			container = transformers.HelmTransform(&transformation)
		} else {
			container = transformation.Container
		}
		initContainers = append(initContainers, container)
	}
	return initContainers
}

func createPod(p *parameterizer.Parameterizer) *v1.Pod {
	name := generateName(p.ObjectMeta.Name)
	sourceContainers := createResourceContainers(name, p.Spec.Resources)
	userInputContainers := createResourceContainers(name, p.Spec.UserInputs)
	transformationContainers := createTransformationContainers(name, p.Spec.Transformations)
	initContainers := append(sourceContainers, userInputContainers...)
	initContainers = append(initContainers, transformationContainers...)
	volumes := p.Spec.Volumes
	container := v1.Container{
		Name:    "krm-result",
		Image:   "alpine:3.7",
		Command: []string{"echo", name + " -c " + initContainers[len(initContainers)-1].Name},
	}
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.PodSpec{
			InitContainers: initContainers,
			Volumes:        volumes,
			Containers:     []v1.Container{container},
		},
	}
	return pod
}

func createmanifest(p parameterizer.Parameterizer) (string, *os.File, string, error) {
	pod := createPod(&p)
	podName := pod.ObjectMeta.Name
	content, err := MarshallObj(pod, "yaml")
	if err != nil {
		return podName, nil, "", err
	}
	tmpf, err := ioutil.TempFile("/tmp", "krm")
	if err != nil {
		return podName, nil, "", err
	}
	if _, err := tmpf.Write(content); err != nil {
		return podName, nil, "", err
	}
	if err := tmpf.Close(); err != nil {
		return podName, nil, "", err
	}
	return podName, tmpf, string(content), nil
}

func buildcmds(cmds []string) string {
	var res string
	for _, cmd := range cmds {
		wocmd := strings.Split(cmd, " ")[1:]
		res += strings.Join(wocmd, " ")
	}
	return res
}

func kubectl(withstderr bool, cmd string, args ...string) (string, error) {
	kubectlbin, err := executecmd(false, "which", "kubectl")
	if err != nil {
		return "", err
	}
	all := append([]string{cmd}, args...)
	result, err := executecmd(withstderr, kubectlbin, all...)
	if err != nil {
		return "", err
	}
	return result, nil
}

func executecmd(withstderr bool, cmd string, args ...string) (string, error) {
	result := ""
	var out bytes.Buffer
	c := exec.Command(cmd, args...)
	c.Env = os.Environ()
	if withstderr {
		c.Stderr = os.Stderr
	}
	c.Stdout = &out
	err := c.Run()
	if err != nil {
		return result, err
	}
	result = strings.TrimSpace(out.String())
	return result, nil
}
