//go:build tools
// +build tools

// The only purpose of this file is to track tool dependencies.
// It's never built or included anywhere.
// See: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package cleoutils

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/99designs/gqlgen/graphql/introspection"
	_ "github.com/Khan/genqlient"
)
