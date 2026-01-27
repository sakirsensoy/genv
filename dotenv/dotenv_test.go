package dotenv_test

import (
	"os"
	"testing"

	"github.com/sakirsensoy/genv/dotenv"
)

func TestEnvData(t *testing.T) {
	// Test cases following motdotla/dotenv parsing rules exactly
	// https://github.com/motdotla/dotenv
	// Source: https://raw.githubusercontent.com/motdotla/dotenv/refs/heads/master/lib/main.js
	expectedValues := map[string]string{
		// Basic values
		"KEY_1": "rubic",
		"KEY_2": "rubic",
		"KEY_3": "rubic",
		"KEY_4": "rubic",
		"KEY_5": "123",
		"KEY_6": "127.0.0.1",
		"KEY_7": "127.0.0.1",
		"KEY_8": "1.23",
		"KEY_9": "1.23",

		// Comments
		"KEY_10": "value",                             // inline comment stripped
		"KEY_11": "http://localhost:4443/storage/v1/", // inline comment stripped
		"KEY_12": "value with # hash inside",          // hash inside double quotes preserved
		"KEY_13": "value with # hash inside",          // hash inside single quotes preserved
		"KEY_14": "value",                             // # starts comment even without space (motdotla behavior)
		"KEY_15": "",                                  // empty value
		"KEY_16": "",                                  // empty value with only comment
		"KEY_17": "true",                              // boolean with inline comment

		// Export prefix
		"KEY_18": "exported_value",
		"KEY_19": "exported quoted",

		// Escape sequences in double quotes (ONLY \n and \r are expanded per motdotla/dotenv)
		"KEY_20": "line1\nline2",  // \n expanded to newline
		"KEY_21": "tab\\there",    // \t NOT expanded (stays as literal \t)
		"KEY_22": "quote\\\"inside", // \" NOT expanded (stays as literal \")
		"KEY_23": "backslash\\\\", // \\ NOT expanded (stays as literal \\)

		// Single quotes are raw
		"KEY_24": "no\\nescape", // no escape processing

		// Backticks are raw
		"KEY_25": "backtick value",

		// JSON inner quotes maintained
		"KEY_26": "{\"foo\":\"bar\"}",
		"KEY_27": "{\"foo\": \"bar\", \"baz\": 123}",

		// Whitespace handling
		"KEY_28": "  spaced  ",      // double quotes preserve whitespace
		"KEY_29": "  spaced  ",      // single quotes preserve whitespace
		"KEY_30": "unquoted spaced", // unquoted trims whitespace

		// Empty quoted values
		"KEY_31": "",
		"KEY_32": "",

		// Colon separator support
		"KEY_33": "colon_value",
		"KEY_34": "colon quoted",
	}

	os.Clearenv()

	err := dotenv.Load("./../testdata/test.env")
	if err != nil {
		t.Error("Error reading file")
	}

	for key, expectVal := range expectedValues {
		val, exists := os.LookupEnv(key)
		if !exists {
			t.Errorf("Key '%v' not found in environment", key)
			continue
		}
		if val != expectVal {
			t.Errorf("Key '%v': expected %q got %q", key, expectVal, val)
		}
	}
}
