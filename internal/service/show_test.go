package service

import "testing"

func TestShow_ValidOperation(t *testing.T) {
	result, err := Show("../../testdata/valid/test.yml", "GET", "/pets/{petId}")
	if err != nil {
		t.Fatalf("Show returned unexpected error: %v", err)
	}

	if result.Method != "GET" {
		t.Errorf("Method = %q, want %q", result.Method, "GET")
	}
	if result.Summary != "Info for a specific pet" {
		t.Errorf("Summary = %q, want %q", result.Summary, "Info for a specific pet")
	}

	if len(result.Parameters) != 1 {
		t.Fatalf("len(Parameters) = %d, want 1", len(result.Parameters))
	}
	param := result.Parameters[0]
	if param.Name != "petId" || param.In != "path" || !param.Required || param.Type != "string" {
		t.Errorf("Parameters[0] = %+v, want {Name:petId In:path Required:true Type:string}", param)
	}

	foundOK := false
	for _, r := range result.Responses {
		if r.Code == "200" {
			foundOK = true
			if r.SchemaName != "Pets" {
				t.Errorf("200 response SchemaName = %q, want %q", r.SchemaName, "Pets")
			}
		}
	}
	if !foundOK {
		t.Errorf("Responses = %+v, want a 200 response to be present", result.Responses)
	}
}

func TestShow_MethodNotDefined(t *testing.T) {
	_, err := Show("../../testdata/valid/test.yml", "DELETE", "/pets/{petId}")
	if err == nil {
		t.Fatal("Show returned nil error for a method not defined on the path, want an error")
	}
}

func TestShow_PathNotFound(t *testing.T) {
	_, err := Show("../../testdata/valid/test.yml", "GET", "/does-not-exist")
	if err == nil {
		t.Fatal("Show returned nil error for an undefined path, want an error")
	}
}
