package parameterizer

import "testing"

func TestParse(t *testing.T) {
	p, err := Parse("../../test/install-ghost-with-helm.yaml")
	if err != nil {
		t.Errorf(err.Error())
	}
	got := p.Kind
	want := "Parameterizer"
	if got != want {
		t.Errorf("parameterizer.Parse(\"install-ghost-with-helm.yaml\") => %q, want %q", got, want)
	}
}
