package parameterizer

import "testing"

func TestParse(t *testing.T) {
	p, err := Parse("../../test/install-ghost-with-helm.yaml")
	if err != nil {
		t.Errorf(err.Error())
	}
	got := p
	want := "Parameterizer"
	if got.Kind != want {
		t.Errorf("parameterizer.Parse(\"install-ghost-with-helm.yaml\") => %q, want %q", got.Kind, want)
	}
	want = "kubernetes.sh/v1alpha1"
	if got.APIVersion != want {
		t.Errorf("parameterizer.Parse(\"install-ghost-with-helm.yaml\") => %q, want %q", got.APIVersion, want)
	}
	want = "install-ghost"
	if got.Metadata.Name != want {
		t.Errorf("parameterizer.Parse(\"install-ghost-with-helm.yaml\") => %q, want %q", got.Metadata.Name, want)
	}
}
