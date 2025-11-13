package logging

import (
	"os"
	"testing"
)

// TestLogging_GenerateID checks if the generated ID has the correct length.
func TestLogging_GenerateID(t *testing.T) {
	id := GenerateID()
	if len(id) != 32 {
		t.Errorf("Expected ID length 32, got %d", len(id))
	}
}

// TestLogging_CreateAppDataFolder verifies that the application data folder is created successfully.
func TestLogging_CreateAppDataFolder(t *testing.T) {
	dir, err := CreateAppDataFolder("testapp")
	if err != nil {
		t.Fatalf("Failed to create app data folder: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("App data folder does not exist: %s", dir)
	}
	// Cleanup
	_ = os.Remove(dir)
}

// TestLogging_OpenLogFile tests the OpenLogFile function.
func TestLogging_OpenLogFile(t *testing.T) {
	fileName := "testlog.log"
	f, err := OpenLogFile(fileName)
	if err != nil {
		t.Fatalf("Failed to open log file: %v", err)
	}
	if f == nil {
		t.Error("Expected file handle, got nil")
	}
	f.Close()
	_ = os.Remove(fileName)
}

// TestLogging_LoggerOptions checks the LoggerOptions function.
func TestLogging_LoggerOptions(t *testing.T) {
	opts := LoggerOptions()
	if opts.AddSource {
		t.Error("Expected AddSource to be false")
	}
}
