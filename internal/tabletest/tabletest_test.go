package tabletest

import (
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	type want struct {
		value int
	}
	cases := []Case[string, want]{
		{Name: "first", In: "one", Want: want{value: 1}},
		{Name: "second", In: "two", Want: want{value: 2}},
	}

	var got []string
	Run(t, cases, func(t *testing.T, in string, want want) {
		got = append(got, in)
		if want.value != len(got) {
			t.Fatalf("value = %d, want %d", want.value, len(got))
		}
	})

	if want := []string{"one", "two"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("inputs = %v, want %v", got, want)
	}
}
