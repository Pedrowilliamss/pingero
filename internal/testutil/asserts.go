package testutil

import (
	"reflect"
	"testing"
)

func ExpectEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func ExpectDeepEqual[T any](t *testing.T, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v want: %v", got, want)
	}
}

func ExpectNotNil(t *testing.T, value any) {
	t.Helper()
	if value == nil {
		t.Errorf("value: %v sould not be nil", value)
	}
}

func ExpectTrue(t *testing.T, value bool) {
	t.Helper()
	if !value {
		t.Errorf("value: %v should be true", value)
	}
}

func ExpectFalse(t *testing.T, value bool) {
	t.Helper()
	if value {
		t.Errorf("value: %v should be false", value)
	}
}
