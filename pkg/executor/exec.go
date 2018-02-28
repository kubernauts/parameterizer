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
		fmt.Printf("Using image %v to run commands %v\n", a.Image, a.Commands)
		// kubectl run -it --rm pexecutor --overrides='
		// {
		//   "apiVersion": "batch/v1",
		//   "spec": {
		//     "template": {
		//       "spec": {
		//         "containers": [
		//           {
		//             "name": "$RANDOM",
		//             "image": "$IMAGE",
		//             "args": [
		//               "$COMMAND"
		//             ],
		//             "stdin": true,
		//             "stdinOnce": true,
		//             "tty": true,
		//             "volumeMounts": [{
		//               "mountPath": "$VMP0",
		//               "name": "$VMN0"
		//              },
		//						  ...
		//            ]
		//           }
		//         ],
		//         "volumes": [{
		//           "name":"$VMN0",
		//           "emptyDir":{}
		//         }]
		//       }
		//     }
		//   }
		// }
		// '  -image=$IMAGE --restart=Never -- $COMMAND

	}
	return nil
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