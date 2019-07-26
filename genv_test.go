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
