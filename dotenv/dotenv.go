package dotenv

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// Load method allows environment variables to be loaded from the desired file.
//
//	dotenv.Load(".envfile")
func Load(path ...string) (err error) {
	var dotenvPath string
	if len(path) > 0 {
		dotenvPath = path[0]
	} else {
		dotenvPath = ".env"
	}

	source, err := os.Open(dotenvPath)
	if err != nil {
		log.Print(err)
		return
	}
	defer source.Close()

	variables, err := parse(source)
	if err != nil {
		log.Print(err)
		return
	}

	for key, val := range variables {
		if _, ok := os.LookupEnv(key); !ok {
			os.Setenv(key, val)
		}
	}

	return
}

func parse(source *os.File) (variables map[string]string, err error) {
	variables = make(map[string]string)

	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		key, val, ok := parseLine(scanner.Text())
		if ok {
			variables[key] = val
		}
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return
}

// parseLine parses a single line from a .env file.
// It follows the motdotla/dotenv parsing rules exactly:
// https://github.com/motdotla/dotenv
//
// Rules:
// - BASIC=basic becomes {BASIC: 'basic'}
// - empty lines are skipped
// - lines beginning with # are treated as comments
// - # marks the beginning of a comment (unless when the value is wrapped in quotes)
// - empty values become empty strings (EMPTY= becomes {EMPTY: ''})
// - inner quotes are maintained (think JSON) (JSON={"foo": "bar"})
// - whitespace is removed from both ends of unquoted values
// - single and double quoted values have quotes removed
// - single and double quoted values maintain whitespace from both ends
// - double quoted values expand \n and \r only
// - single quoted and backtick values are raw (no escape processing)
// - export prefix is supported
// - colon separator is supported (KEY: value)
func parseLine(line string) (key string, val string, ok bool) {
	// Trim leading/trailing whitespace
	line = strings.TrimSpace(line)

	// Skip empty lines and comment lines
	if line == "" || strings.HasPrefix(line, "#") {
		return
	}

	// Handle 'export' prefix (export KEY=value)
	if strings.HasPrefix(line, "export ") {
		line = strings.TrimPrefix(line, "export ")
		line = strings.TrimSpace(line)
	}

	// Find separator (= or :)
	// The original regex: (?:\s*=\s*?|:\s+?)
	// = can have optional whitespace on both sides
	// : must be followed by at least one whitespace
	sepIndex := -1
	sepLen := 1

	for i := 0; i < len(line); i++ {
		if line[i] == '=' {
			sepIndex = i
			break
		}
		if line[i] == ':' && i+1 < len(line) && (line[i+1] == ' ' || line[i+1] == '\t') {
			sepIndex = i
			break
		}
	}

	if sepIndex == -1 {
		return
	}

	key = strings.TrimSpace(line[:sepIndex])
	if key == "" {
		return
	}

	// Get the raw value part (everything after separator)
	rawValue := ""
	if sepIndex+sepLen < len(line) {
		rawValue = line[sepIndex+sepLen:]
	}

	// Parse the value
	val = parseValue(rawValue)
	ok = true

	return
}

// parseValue parses the value part of a key=value pair.
// Follows motdotla/dotenv behavior exactly:
// https://raw.githubusercontent.com/motdotla/dotenv/refs/heads/master/lib/main.js
//
// The original logic:
// 1. value = value.trim()
// 2. maybeQuote = value[0]
// 3. value = value.replace(/^(['"`])([\s\S]*)\1$/mg, '$2')  <- removes matching quotes
// 4. if maybeQuote === '"': expand \n and \r
func parseValue(rawValue string) string {
	// Trim the value first (like the original does with value.trim())
	value := strings.TrimSpace(rawValue)

	if value == "" {
		return ""
	}

	// Remember the first character to check if it was double-quoted
	maybeQuote := value[0]

	// Remove surrounding quotes if they match (same type at start and end)
	// This mimics: value.replace(/^(['"`])([\s\S]*)\1$/mg, '$2')
	if len(value) >= 2 {
		lastChar := value[len(value)-1]
		if (maybeQuote == '"' || maybeQuote == '\'' || maybeQuote == '`') && maybeQuote == lastChar {
			// Remove the surrounding quotes
			value = value[1 : len(value)-1]

			// Only expand \n and \r for double-quoted values
			if maybeQuote == '"' {
				value = strings.ReplaceAll(value, "\\n", "\n")
				value = strings.ReplaceAll(value, "\\r", "\r")
			}
			return value
		}
	}

	// Unquoted value - strip comments and trim whitespace
	return parseUnquotedValue(value)
}

// parseUnquotedValue parses an unquoted value.
// Unquoted values:
// - Have whitespace trimmed from both ends
// - Have comments stripped (everything after #)
func parseUnquotedValue(s string) string {
	// Find the comment start
	commentIndex := strings.Index(s, "#")
	if commentIndex != -1 {
		s = s[:commentIndex]
	}

	return strings.TrimSpace(s)
}
