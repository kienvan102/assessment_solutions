package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadJSONFile_Success(t *testing.T) {
	// Create a temporary JSON file
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "test.json")
	
	testData := `{"name": "test", "value": 123}`
	if err := os.WriteFile(jsonFile, []byte(testData), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Test loading
	var result map[string]interface{}
	loadedPath, err := LoadJSONFile(&result, jsonFile)
	
	if err != nil {
		t.Fatalf("LoadJSONFile failed: %v", err)
	}
	
	if loadedPath != jsonFile {
		t.Errorf("Expected loadedPath %s, got %s", jsonFile, loadedPath)
	}
	
	if result["name"] != "test" {
		t.Errorf("Expected name 'test', got %v", result["name"])
	}
	
	if result["value"].(float64) != 123 {
		t.Errorf("Expected value 123, got %v", result["value"])
	}
}

func TestLoadJSONFile_MultiplePaths(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "test.json")
	
	testData := `{"found": true}`
	if err := os.WriteFile(jsonFile, []byte(testData), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Try multiple paths, only the last one exists
	var result map[string]interface{}
	loadedPath, err := LoadJSONFile(&result, 
		"/nonexistent/path1.json",
		"/nonexistent/path2.json",
		jsonFile,
	)
	
	if err != nil {
		t.Fatalf("LoadJSONFile failed: %v", err)
	}
	
	if loadedPath != jsonFile {
		t.Errorf("Expected loadedPath %s, got %s", jsonFile, loadedPath)
	}
	
	if result["found"] != true {
		t.Errorf("Expected found=true, got %v", result["found"])
	}
}

func TestLoadJSONFile_NoPathsProvided(t *testing.T) {
	var result map[string]interface{}
	_, err := LoadJSONFile(&result)
	
	if err == nil {
		t.Error("Expected error when no paths provided, got nil")
	}
}

func TestLoadJSONFile_FileNotFound(t *testing.T) {
	var result map[string]interface{}
	_, err := LoadJSONFile(&result, "/nonexistent/file.json")
	
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestLoadJSONFile_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "invalid.json")
	
	// Write invalid JSON
	if err := os.WriteFile(jsonFile, []byte("not valid json {"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	var result map[string]interface{}
	_, err := LoadJSONFile(&result, jsonFile)
	
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestLoadJSONFile_StructTarget(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "struct.json")
	
	type TestStruct struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	
	testData := `{"name": "example", "count": 42}`
	if err := os.WriteFile(jsonFile, []byte(testData), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	var result TestStruct
	_, err := LoadJSONFile(&result, jsonFile)
	
	if err != nil {
		t.Fatalf("LoadJSONFile failed: %v", err)
	}
	
	if result.Name != "example" {
		t.Errorf("Expected name 'example', got %s", result.Name)
	}
	
	if result.Count != 42 {
		t.Errorf("Expected count 42, got %d", result.Count)
	}
}
