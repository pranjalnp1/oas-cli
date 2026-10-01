package service

import (
	"encoding/json"
	"testing"
)

func TestCurl_PathAndQueryParams(t *testing.T) {
	result, err := Curl("../../testdata/valid/test.yml", "GET", "/pets", "")
	if err != nil {
		t.Fatalf("Curl returned unexpected error: %v", err)
	}

	want := "http://petstore.swagger.io/v1/pets?limit=0"
	if result.URL != want {
		t.Errorf("URL = %q, want %q", result.URL, want)
	}
	if result.Body != "" {
		t.Errorf("Body = %q, want empty (operation has no request body)", result.Body)
	}
}

func TestCurl_PathParamSubstitution(t *testing.T) {
	result, err := Curl("../../testdata/valid/test.yml", "GET", "/pets/{petId}", "")
	if err != nil {
		t.Fatalf("Curl returned unexpected error: %v", err)
	}

	want := "http://petstore.swagger.io/v1/pets/string"
	if result.URL != want {
		t.Errorf("URL = %q, want %q", result.URL, want)
	}
}

func TestCurl_RelativeServerRequiresBaseURLOverride(t *testing.T) {
	_, err := Curl("../../testdata/validLiveAPI/petStore3.json", "GET", "/pet/findByStatus", "")
	if err == nil {
		t.Fatal("Curl returned nil error for a relative server URL with no override, want an error")
	}
}

func TestCurl_BaseURLOverride(t *testing.T) {
	result, err := Curl("../../testdata/validLiveAPI/petStore3.json", "GET", "/pet/findByStatus", "https://petstore3.swagger.io/api/v3")
	if err != nil {
		t.Fatalf("Curl returned unexpected error: %v", err)
	}

	want := "https://petstore3.swagger.io/api/v3/pet/findByStatus?status=available"
	if result.URL != want {
		t.Errorf("URL = %q, want %q (expected the enum's first value to be used)", result.URL, want)
	}
}

func TestCurl_NestedRefAndEnumInBody(t *testing.T) {
	result, err := Curl("../../testdata/validLiveAPI/petStore3.json", "POST", "/pet", "https://petstore3.swagger.io/api/v3")
	if err != nil {
		t.Fatalf("Curl returned unexpected error: %v", err)
	}

	if result.Headers["Content-Type"] != "application/json" {
		t.Errorf("Content-Type header = %q, want %q", result.Headers["Content-Type"], "application/json")
	}

	var body map[string]interface{}
	if err := json.Unmarshal([]byte(result.Body), &body); err != nil {
		t.Fatalf("generated body is not valid JSON: %v\nbody: %s", err, result.Body)
	}

	category, ok := body["category"].(map[string]interface{})
	if !ok {
		t.Fatalf("category = %v (%T), want a nested object (expected $ref inside properties to be resolved)", body["category"], body["category"])
	}
	if _, hasName := category["name"]; !hasName {
		t.Errorf("category = %+v, want a \"name\" field from the resolved Category schema", category)
	}

	status, _ := body["status"].(string)
	if status != "available" {
		t.Errorf("status = %q, want %q (expected the enum's first value, not a generic placeholder)", status, "available")
	}
}

func TestCurl_PathNotFound(t *testing.T) {
	_, err := Curl("../../testdata/valid/test.yml", "GET", "/does-not-exist", "")
	if err == nil {
		t.Fatal("Curl returned nil error for an undefined path, want an error")
	}
}
