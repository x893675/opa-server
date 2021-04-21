package names

import (
	"strings"
	"testing"
)

func TestSimpleNameGenerator(t *testing.T) {
	name := SimpleNameGenerator.GenerateName("foo")
	if !strings.HasPrefix(name, "foo") || name == "foo" {
		t.Errorf("unexpected name: %s", name)
	}
}
