package main

import (
	"fmt"
	"os"

	// Import the autoload package to trigger init() function
	_ "github.com/sakirsensoy/genv/dotenv/autoload"
)

func main() {
	// Print the environment variables that should be loaded by autoload
	fmt.Printf("AUTOLOAD_TEST_KEY=%s\n", os.Getenv("AUTOLOAD_TEST_KEY"))
	fmt.Printf("AUTOLOAD_TEST_NUMBER=%s\n", os.Getenv("AUTOLOAD_TEST_NUMBER"))
	fmt.Printf("AUTOLOAD_TEST_BOOL=%s\n", os.Getenv("AUTOLOAD_TEST_BOOL"))
}
