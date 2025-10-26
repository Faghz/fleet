//go:build tools
// +build tools

package tools

// This file imports packages that are used when running go generate.
// These tools are not included in the final binary.

//go:generate go install github.com/swaggo/swag/cmd/swag@latest
//go:generate go install github.com/golang/mock/mockgen@v1.6.0
//go:generate go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

import (
	// Build and documentation tools
	_ "github.com/swaggo/swag/cmd/swag"

	// Testing tools
	_ "github.com/golang/mock/mockgen"

	// Linting tools
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
