package executor

import (
	"bytes"
	"fmt"
	"github.com/kubernauts/parameterizer/pkg/parameterizer"
	"io/ioutil"
	// "k8s.io/api/apps/v1beta1"
	"github.com/ghodss/yaml"
	"github.com/kubernauts/parameterizer/pkg/executor/transformers"
	"github.com/pborman/uuid"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"os/exec"
	"strings"
)

// returns the pod name once the job is completed
func waitJobToComplete(job string) error {
	completed := false
	var o batchv1.Job

	for !completed {
		res, err := kubectl(true, "get", "jobs", "-o", "yaml", job)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal([]byte(res), &o)
		if err != nil {
			return err
		}
		completed = o.Status.Succeeded > 0 || o.Status.Failed > 0
	}

	return nil
}

func getJobPodName(job string) string {
	res, _ := kubectl(true, "get", "pods",
		"--show-all",
		"--selector=job-name="+job,
		"--output=jsonpath={.items..metadata.name}")
	return res
}

// Run executes the Parameterizer resource's transformation
// steps as defined in the apply sub-resource.
func Run(p parameterizer.Parameterizer) (err error) {
	// we create a temporary manifest file with all
	// the necessary settings in there
	jobName, mf, mc, err := createmanifest(p)
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
	fmt.Fprintf(os.Stderr, "Using manifest:\n%v\n", mc)
	cmd := []string{"create", "-f", mfn}
	fmt.Fprintf(os.Stderr, "Executing command: %v\n", strings.Join(cmd, " "))
	res, err := kubectl(true, cmd[0], cmd[1:]...)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%v\n", res)
	fmt.Fprintf(os.Stderr, "Waiting for Job %s to complete....", jobName)
	err = waitJobToComplete(jobName)
	if err != nil {
		return err
	}
	podName := getJobPodName(jobName)

	fmt.Fprintf(os.Stderr, "%s\n", podName)
	res, err = kubectl(true, "logs", podName)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%v\n", res)
	res, err = kubectl(true, "logs", podName, "-c", res)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", res)
	// res, err = kubectl(true, "delete", "po", podName, "--force")
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("%v\n", res)
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

func CreatePod(p *parameterizer.Parameterizer) *batchv1.Job {
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
		Command: []string{"echo", initContainers[len(initContainers)-1].Name},
	}
	pod := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					RestartPolicy:  "Never",
					InitContainers: initContainers,
					Volumes:        volumes,
					Containers:     []v1.Container{container},
				},
			},
		},
	}
	return pod
}

func createmanifest(p parameterizer.Parameterizer) (string, *os.File, string, error) {
	pod := CreatePod(&p)
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
