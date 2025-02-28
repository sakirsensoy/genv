package dotenv_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestAutoload(t *testing.T) {
	// Create a temporary .env file in the current directory
	err := os.WriteFile(".env", []byte("AUTOLOAD_TEST_KEY=autoload_value\nAUTOLOAD_TEST_NUMBER=12345\nAUTOLOAD_TEST_BOOL=true"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary .env file: %v", err)
	}
	defer os.Remove(".env") // Clean up after test

	// Run the test program that imports the autoload package
	cmd := exec.Command("go", "run", "../testdata/autoload_test_program.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run test program: %v\nOutput: %s", err, output)
	}

	// Check the output for expected environment variables
	outputStr := string(output)
	expectedValues := map[string]string{
		"AUTOLOAD_TEST_KEY":    "autoload_value",
		"AUTOLOAD_TEST_NUMBER": "12345",
		"AUTOLOAD_TEST_BOOL":   "true",
	}

	for key, expectedValue := range expectedValues {
		expectedOutput := key + "=" + expectedValue
		if !strings.Contains(outputStr, expectedOutput) {
			t.Errorf("Expected output to contain %q, but got: %s", expectedOutput, outputStr)
		}
	}
}
