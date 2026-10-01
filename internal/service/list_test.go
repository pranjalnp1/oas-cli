package service

import (
	"reflect"
	"testing"
)

func TestList_ValidSpec(t *testing.T) {
	result, err := List("../../testdata/valid/test.yml")
	if err != nil {
		t.Fatalf("List returned unexpected error: %v", err)
	}

	want := []OperationSummary{
		{Method: "GET", Path: "/pets"},
		{Method: "POST", Path: "/pets"},
		{Method: "GET", Path: "/pets/{petId}"},
	}

	if !reflect.DeepEqual(result.Operations, want) {
		t.Errorf("Operations = %+v, want %+v", result.Operations, want)
	}
}

func TestList_StableOrderAcrossRuns(t *testing.T) {
	first, err := List("../../testdata/valid/test.yml")
	if err != nil {
		t.Fatalf("List returned unexpected error: %v", err)
	}
	second, err := List("../../testdata/valid/test.yml")
	if err != nil {
		t.Fatalf("List returned unexpected error: %v", err)
	}

	if !reflect.DeepEqual(first.Operations, second.Operations) {
		t.Errorf("List produced different operation order across runs:\nfirst:  %+v\nsecond: %+v", first.Operations, second.Operations)
	}
}

func TestList_FileNotFound(t *testing.T) {
	_, err := List("../../testdata/valid/does-not-exist.yml")
	if err == nil {
		t.Fatal("List returned nil error for a nonexistent file, want an error")
	}
}
