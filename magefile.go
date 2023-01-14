//go:build mage

package main

import (
	"github.com/carolynvs/magex/mgx"
	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/shx"
	"github.com/magefile/mage/sh"
)

var must = shx.CommandBuilder{StopOnError: true}
var mustFrontend = shx.CommandBuilder{StopOnError: true, Dir: "frontend"}

// Clean up build artifacts.
func Clean() {
	mgx.Must(sh.Rm("cleosrv/cleosrv/frontend_build/"))
	mgx.Must(sh.Rm("dist/"))
	mgx.Must(sh.Rm("frontend/build/"))
	mgx.Must(sh.Rm("website/public/"))
}

// Generate code (for example after changing the schema).
func Generate() {
	must.RunV("go", "generate", "./...")
	mustFrontend.RunV("npm", "run", "generate")
}

// Lint all code, fixing things automatically where possible.
func Lint() {
	must.RunV(
		"go",
		"run",
		"github.com/golangci/golangci-lint/cmd/golangci-lint",
		"run",
		"--fix",
	)
	must.RunV(
		"go",
		"run",
		"github.com/golangci/golangci-lint/cmd/golangci-lint",
		"run",
		"--fix",
		"--build-tags",
		"production",
	)
	mustFrontend.RunV("npm", "run", "lint")
}

// Build cleosrv and cleoc binaries for the current platform.
func Build() {
	must.RunV(
		"go",
		"run",
		"github.com/goreleaser/goreleaser",
		"build",
		"--snapshot",
		"--single-target",
	)
}

// EnsureMage installs mage globally if it's not already installed.
func EnsureMage() {
	pkg.EnsureMage("")
}

// Test executes all tests except the e2e tests.
func Test() {
	must.RunV("go", "test", "./...")
	mustFrontend.RunV("npm", "test", "a", "--", "--watchAll=false")
}

// InstallDeps installs all dependencies.
func InstallDeps() {
	must.RunV("go", "mod", "tidy")
	must.RunV("go", "mod", "download")
	mustFrontend.RunV("npm", "install")
	shx.Command("npm", "install").In("e2e_tests").Must().RunV()
}
