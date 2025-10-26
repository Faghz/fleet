package main

import (
	"fmt"
	"os"

	"aidanwoods.dev/go-paseto"
)

func main() {
	key := paseto.NewV4SymmetricKey()

	// Check if we should output to file or stdout
	if len(os.Args) > 1 && os.Args[1] == "--save" {
		// Save to .env file or create a new auth key file
		filename := ".auth_key"
		if len(os.Args) > 2 {
			filename = os.Args[2]
		}

		err := os.WriteFile(filename, []byte(fmt.Sprintf("AUTH_KEY=%s\n", key)), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to file %s: %v\n", filename, err)
			os.Exit(1)
		}
		fmt.Printf("Auth key saved to %s\n", filename)
	} else {
		// Output to stdout
		fmt.Println(key.ExportHex())
	}
}
