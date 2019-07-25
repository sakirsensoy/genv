package dotenv

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

const keyValRegex = `^\s*([\w.-]+)\s*=\s*(.*)?\s*$`

// Load method allows environment variables to be loaded from the desired file.
//
//				dotenv.Load(".envfile")
//
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

	r, err := regexp.Compile(keyValRegex)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		key, val := parseLine(scanner.Text(), r)
		if key != "" && val != "" {
			variables[key] = val
		}
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return
}

func parseLine(line string, r *regexp.Regexp) (key string, val string) {

	matches := r.FindStringSubmatch(line)

	if len(matches) == 3 {
		key, val = matches[1], matches[2]

		if val != "" {
			end := len(val) - 1
			isDoubleQuoted := val[0] == '"' && val[end] == '"'
			isSingleQuoted := val[0] == '\'' && val[end] == '\''

			if isSingleQuoted || isDoubleQuoted {
				val = val[1:end]
			} else {
				val = strings.Trim(val, " ")
			}
		}
	}

	return
}
