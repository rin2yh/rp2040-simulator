package tabletestlinttest

import "testing"

func productionValidate(int) error { return nil }

func TestNameOutsideTestFile(t *testing.T) {
	if err := productionValidate(1); err == nil {
		t.Fatal("accepted 1")
	}
	if err := productionValidate(2); err == nil {
		t.Fatal("accepted 2")
	}
	if err := productionValidate(3); err == nil {
		t.Fatal("accepted 3")
	}
}
