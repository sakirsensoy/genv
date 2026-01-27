package genv_test

import (
	"fmt"
	"os"
	"sync"
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

// ==================== Multi-threaded Tests ====================

func TestConcurrentKeyAccess(t *testing.T) {
	// Test concurrent access to genv.Key() for the same key
	os.Clearenv()
	os.Setenv("CONCURRENT_KEY", "concurrent_value")

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	errors := make(chan error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			val := genv.Key("CONCURRENT_KEY").String()
			if val != "concurrent_value" {
				errors <- fmt.Errorf("expected 'concurrent_value', got '%s'", val)
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

func TestConcurrentDifferentKeys(t *testing.T) {
	// Test concurrent access to different keys
	os.Clearenv()

	const numKeys = 50
	for i := 0; i < numKeys; i++ {
		os.Setenv(fmt.Sprintf("DIFF_KEY_%d", i), fmt.Sprintf("value_%d", i))
	}

	var wg sync.WaitGroup
	wg.Add(numKeys * 2) // Each key accessed twice concurrently

	errors := make(chan error, numKeys*2)

	for i := 0; i < numKeys; i++ {
		keyNum := i
		// First goroutine
		go func() {
			defer wg.Done()
			key := fmt.Sprintf("DIFF_KEY_%d", keyNum)
			expected := fmt.Sprintf("value_%d", keyNum)
			val := genv.Key(key).String()
			if val != expected {
				errors <- fmt.Errorf("key %s: expected '%s', got '%s'", key, expected, val)
			}
		}()
		// Second goroutine for same key
		go func() {
			defer wg.Done()
			key := fmt.Sprintf("DIFF_KEY_%d", keyNum)
			expected := fmt.Sprintf("value_%d", keyNum)
			val := genv.Key(key).String()
			if val != expected {
				errors <- fmt.Errorf("key %s: expected '%s', got '%s'", key, expected, val)
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

func TestConcurrentReadWrite(t *testing.T) {
	// Test concurrent read and write operations
	os.Clearenv()
	os.Setenv("RW_KEY", "initial")

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	// Readers
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_ = genv.Key("RW_KEY").String()
		}()
	}

	// Writers
	for i := 0; i < goroutines; i++ {
		iteration := i
		go func() {
			defer wg.Done()
			genv.Key("RW_KEY").Update(fmt.Sprintf("value_%d", iteration))
		}()
	}

	wg.Wait()

	// Final value should be one of the written values
	finalVal := genv.Key("RW_KEY").String()
	if finalVal == "" {
		t.Error("Final value should not be empty")
	}
}

func TestConcurrentTypeConversions(t *testing.T) {
	// Test concurrent access with different type conversions
	os.Clearenv()
	os.Setenv("CONV_INT", "42")
	os.Setenv("CONV_FLOAT", "3.14")
	os.Setenv("CONV_BOOL", "true")
	os.Setenv("CONV_STRING", "hello")

	const goroutines = 25
	var wg sync.WaitGroup
	wg.Add(goroutines * 4)

	errors := make(chan error, goroutines*4)

	// Int conversions
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if val := genv.Key("CONV_INT").Int(); val != 42 {
				errors <- fmt.Errorf("expected int 42, got %d", val)
			}
		}()
	}

	// Float conversions
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if val := genv.Key("CONV_FLOAT").Float(); val != 3.14 {
				errors <- fmt.Errorf("expected float 3.14, got %f", val)
			}
		}()
	}

	// Bool conversions
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if val := genv.Key("CONV_BOOL").Bool(); val != true {
				errors <- fmt.Errorf("expected bool true, got %v", val)
			}
		}()
	}

	// String conversions
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if val := genv.Key("CONV_STRING").String(); val != "hello" {
				errors <- fmt.Errorf("expected string 'hello', got '%s'", val)
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

func TestConcurrentWithDefaults(t *testing.T) {
	// Test concurrent access with default values
	os.Clearenv()

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines * 4)

	errors := make(chan error, goroutines*4)

	// String with default
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			val := genv.Key("UNDEF_STR").Default("default_string").String()
			if val != "default_string" {
				errors <- fmt.Errorf("expected 'default_string', got '%s'", val)
			}
		}()
	}

	// Int with default
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			val := genv.Key("UNDEF_INT").Default(999).Int()
			if val != 999 {
				errors <- fmt.Errorf("expected 999, got %d", val)
			}
		}()
	}

	// Float with default
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			val := genv.Key("UNDEF_FLOAT").Default(1.23).Float()
			if val != 1.23 {
				errors <- fmt.Errorf("expected 1.23, got %f", val)
			}
		}()
	}

	// Bool with default
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			val := genv.Key("UNDEF_BOOL").Default(true).Bool()
			if val != true {
				errors <- fmt.Errorf("expected true, got %v", val)
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

// Benchmark concurrent access
func BenchmarkConcurrentKeyAccess(b *testing.B) {
	os.Setenv("BENCH_KEY", "bench_value")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = genv.Key("BENCH_KEY").String()
		}
	})
}

func BenchmarkConcurrentDifferentKeys(b *testing.B) {
	for i := 0; i < 100; i++ {
		os.Setenv(fmt.Sprintf("BENCH_KEY_%d", i), fmt.Sprintf("value_%d", i))
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("BENCH_KEY_%d", i%100)
			_ = genv.Key(key).String()
			i++
		}
	})
}
