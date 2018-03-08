package executor

import (
	"github.com/kubernauts/parameterizer/pkg/parameterizer"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {

}

func TestCreatePod(t *testing.T) {
	p, err := parameterizer.Parse("../../test/basic.yaml")
	require.NoError(t, err)
	pod := createPod(&p)
	want := "krm-exec-" + p.ObjectMeta.Name
	if !strings.HasPrefix(pod.ObjectMeta.Name, want) {
		t.Fatal("Pod name must start with " + want + ", got: " + pod.ObjectMeta.Name)
	}
	PrintObj(pod, "yaml")
}
