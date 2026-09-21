package tabletestlinttest

import "testing"

func validate(int) error      { return nil }
func validateOther(int) error { return nil }
func parse(int) int           { return 0 }

func TestRepeated(t *testing.T) {
	if err := validate(27); err == nil { // want "repeated test cases can be expressed as a table-driven test"
		t.Fatal("accepted 27")
	}
	if err := validate(20001); err == nil {
		t.Fatal("accepted 20001")
	}
	if err := validate(1 << 30); err == nil {
		t.Fatal("accepted maximum")
	}
}

func TestOnlyTwice(t *testing.T) {
	if err := validate(27); err == nil {
		t.Fatal("accepted 27")
	}
	if err := validate(20001); err == nil {
		t.Fatal("accepted 20001")
	}
}

func TestRepeatedStatementGroup(t *testing.T) {
	var got int
	got = parse(1) // want "repeated test cases can be expressed as a table-driven test"
	if got != 10 {
		t.Fatalf("got %d", got)
	}
	got = parse(2)
	if got != 20 {
		t.Fatalf("got %d", got)
	}
	got = parse(3)
	if got != 30 {
		t.Fatalf("got %d", got)
	}
}

func TestSeparated(t *testing.T) {
	if err := validate(27); err == nil {
		t.Fatal("accepted 27")
	}
	t.Log("separator")
	if err := validate(20001); err == nil {
		t.Fatal("accepted 20001")
	}
	t.Log("separator")
	if err := validate(1 << 30); err == nil {
		t.Fatal("accepted maximum")
	}
}

func TestDifferentCalls(t *testing.T) {
	if err := validate(1); err == nil {
		t.Fatal("accepted 1")
	}
	if err := validateOther(2); err == nil {
		t.Fatal("accepted 2")
	}
	if err := validate(3); err == nil {
		t.Fatal("accepted 3")
	}
}

func TestCallsWithoutAssertions(t *testing.T) {
	_ = validate(1)
	_ = validate(2)
	_ = validate(3)
}

type sequence struct {
	value int
}

func (s *sequence) advance(value int) { s.value = value }
func (s *sequence) current() int      { return s.value }

func TestStateTransition(t *testing.T) {
	var state sequence
	state.advance(1)
	if got := state.current(); got != 1 {
		t.Fatal(got)
	}
	state.advance(2)
	if got := state.current(); got != 2 {
		t.Fatal(got)
	}
	state.advance(3)
	if got := state.current(); got != 3 {
		t.Fatal(got)
	}
}

func helperRepeated(t *testing.T) {
	if err := validate(1); err == nil {
		t.Fatal("accepted 1")
	}
	if err := validate(2); err == nil {
		t.Fatal("accepted 2")
	}
	if err := validate(3); err == nil {
		t.Fatal("accepted 3")
	}
}
