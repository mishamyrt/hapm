package hapm

import (
	"errors"
	"testing"
)

func TestHandledErrorContract(t *testing.T) {
	err := errors.New("boom")
	wrapped := HandledError(err)
	if !IsHandledError(wrapped) {
		t.Fatalf("expected wrapped error to be handled")
	}
	if IsHandledError(err) {
		t.Fatalf("expected plain error to not be handled")
	}
}
