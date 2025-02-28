package genv_test

import (
	"os"
	"testing"

	"github.com/sakirsensoy/genv"
)

func TestUndefinedVar(t *testing.T) {

	os.Clearenv()

	key := "UNDEF_VAR"
	cases := []struct {
		Val       interface{}
		ExpectVal interface{}
	}{
		{Val: "UPV", ExpectVal: "UPV"},
		{Val: 11.332, ExpectVal: 11.332},
		{Val: false, ExpectVal: false},
		{Val: 1234, ExpectVal: 1234},
	}

	for _, td := range cases {

		envVar := genv.Key(key).Default(td.Val)

		switch td.ExpectVal.(type) {
		case bool:
			if val := envVar.Bool(); td.ExpectVal != val {
				t.Errorf("Expected '%v' got '%v'", td.ExpectVal, val)
			}
		case float64:
			if val := envVar.Float(); td.ExpectVal != val {
				t.Errorf("Expected '%v' got '%v'", td.ExpectVal, val)
			}
		case int:
			if val := envVar.Int(); td.ExpectVal != val {
				t.Errorf("Expected '%v' got '%v'", td.ExpectVal, val)
			}
		case string:
			if val := envVar.String(); td.ExpectVal != val {
				t.Errorf("Expected '%v' got '%v'", td.ExpectVal, val)
			}
		}
	}
}

func TestDefinedVar(t *testing.T) {

	os.Clearenv()

	key := "DEF_VAR"
	expectVal := "TR"
	os.Setenv(key, expectVal)

	val := genv.Key(key).String()
	if val != expectVal {
		t.Errorf("Expected '%v' got '%v'", expectVal, val)
	}
}

func TestEnvGetSet(t *testing.T) {

	os.Clearenv()

	key := "VAR"
	cases := []struct {
		Val       interface{}
		ExpectVal interface{}
	}{
		{Val: "UPV", ExpectVal: "UPV"},
		{Val: 11.332, ExpectVal: 11.332},
		{Val: false, ExpectVal: false},
		{Val: 1234, ExpectVal: 1234},
	}

	for _, td := range cases {

		genv.Key(key).Update(td.Val)

		envVar := genv.Key(key)

		switch td.ExpectVal.(type) {
		case bool:
			if val := envVar.Bool(); td.ExpectVal != val {
				t.Errorf("Expected '%v' got '%v'", td.ExpectVal, val)
			}
		case float64:
			if val := envVar.Float(); td.ExpectVal != val {
				t.Errorf("Expected '%v' got '%v'", td.ExpectVal, val)
			}
		case int:
			if val := envVar.Int(); td.ExpectVal != val {
				t.Errorf("Expected '%v' got '%v'", td.ExpectVal, val)
			}
		case string:
			if val := envVar.String(); td.ExpectVal != val {
				t.Errorf("Expected '%v' got '%v'", td.ExpectVal, val)
			}
		}
	}
}

func TestQuoteTrimming(t *testing.T) {
	os.Clearenv()

	cases := []struct {
		Key       string
		EnvValue  string
		ExpectVal string
	}{
		{Key: "DOUBLE_QUOTES", EnvValue: `"value"`, ExpectVal: "value"},
		{Key: "SINGLE_QUOTES", EnvValue: `'value'`, ExpectVal: "value"},
		{Key: "NO_QUOTES", EnvValue: "value", ExpectVal: "value"},
		{Key: "MIXED_QUOTES", EnvValue: `"value'`, ExpectVal: `"value'`}, // Should not trim when quotes don't match
		{Key: "ONLY_START_QUOTE", EnvValue: `"value`, ExpectVal: `"value`},
		{Key: "ONLY_END_QUOTE", EnvValue: `value"`, ExpectVal: `value"`},
		{Key: "NESTED_QUOTES", EnvValue: `"'value'"`, ExpectVal: `'value'`}, // Outer quotes are trimmed
		{Key: "EMPTY_QUOTES", EnvValue: `""`, ExpectVal: ``},
		{Key: "EMPTY_SINGLE_QUOTES", EnvValue: `''`, ExpectVal: ``},
	}

	for _, tc := range cases {
		t.Run(tc.Key, func(t *testing.T) {
			os.Setenv(tc.Key, tc.EnvValue)
			val := genv.Key(tc.Key).String()
			if val != tc.ExpectVal {
				t.Errorf("Expected '%v' got '%v'", tc.ExpectVal, val)
			}
		})
	}
}

func TestQuoteTrimmingWithDefault(t *testing.T) {
	os.Clearenv()

	key := "UNDEFINED_WITH_QUOTED_DEFAULT"
	cases := []struct {
		DefaultVal string
		ExpectVal  string
	}{
		{DefaultVal: `"value"`, ExpectVal: "value"},
		{DefaultVal: `'value'`, ExpectVal: "value"},
		{DefaultVal: "value", ExpectVal: "value"},
		{DefaultVal: `"value'`, ExpectVal: `"value'`},
		{DefaultVal: `""`, ExpectVal: ``},
	}

	for i, tc := range cases {
		t.Run(key+string(rune('A'+i)), func(t *testing.T) {
			val := genv.Key(key).Default(tc.DefaultVal).String()
			if val != tc.ExpectVal {
				t.Errorf("Expected '%v' got '%v'", tc.ExpectVal, val)
			}
		})
	}
}
