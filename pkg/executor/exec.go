package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kubernauts/parameterizer/pkg/parameterizer"
)

// Run executes the Parameterizer resource's transformation
// steps as defined in the apply sub-resource.
func Run(p parameterizer.Resource) (err error) {
	for _, a := range p.Spec.Apply {
		cmd := []string{"run", "-it", "--rm", "pexecutor",
			"--image=" + a.Image, "--restart=Never",
			"--", buildcmds(a.Commands)}
		fmt.Printf("Executing following command: %v\n", cmd)
		// res, err := kubectl(true, cmd[0], cmd[1:]...)
		// if err != nil {
		// 	return err
		// }
		// fmt.Printf("%v", res)
	}
	return nil
}

func buildcmds(cmds []string) string {
	var res string
	for _, cmd := range cmds {
		wocmd := strings.Split(cmd, " ")[1:]
		res += strings.Join(wocmd, " ") + " && "
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
