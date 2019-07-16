package dotenv_test

import (
	"os"
	"testing"

	"github.com/sakirsensoy/genv/dotenv"
)

func TestEnvData(t *testing.T) {

	expectedValues := map[string]string{
		"KEY_1": "rubic",
		"KEY_2": "rubic",
		"KEY_3": "rubic",
		"KEY_4": "rubic",
		"KEY_5": "123",
		"KEY_6": "127.0.0.1",
		"KEY_7": "127.0.0.1",
		"KEY_8": "1.23",
		"KEY_9": "1.23",
	}

	os.Clearenv()

	err := dotenv.Load("./../testdata/test.env")

	if err != nil {
		t.Error("Error reading file")
	}

	for key, expectVal := range expectedValues {
		if val := os.Getenv(key); val != expectVal {
			t.Errorf("Expected '%v' got '%v'", expectVal, val)
		}
	}
}
