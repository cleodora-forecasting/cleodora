//go:build mage

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/carolynvs/magex/mgx"
	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/shx"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var must = shx.CommandBuilder{StopOnError: true}
var mustFrontend = shx.CommandBuilder{StopOnError: true, Dir: "frontend"}
var mustE2E = shx.CommandBuilder{StopOnError: true, Dir: "e2e_tests"}

// Clean up build artifacts.
func Clean() {
	mgx.Must(sh.Rm("cleosrv/cleosrv/frontend_build/"))
	mgx.Must(sh.Rm("dist/"))
	mgx.Must(sh.Rm("frontend/build/"))
	mgx.Must(sh.Rm("website/public/"))
}

// Generate code (for example after changing the schema).
func Generate() {
	// Run cleosrv generate first because it updates the GraphQL schema which
	// is used everywhere else.
	_ = must.RunV("go", "generate", "./cleosrv/...")
	_ = must.RunV("go", "generate", "./...")
	_ = mustFrontend.RunV("npm", "run", "generate")
}

// Lint all code, fixing things automatically where possible.
func Lint() {
	// Create an empty file to avoid linter error:
	// cleosrv/cleosrv/production.go: pattern frontend_build: no matching files found (typecheck)
	emptyFile := filepath.Join(
		"cleosrv",
		"cleosrv",
		"frontend_build",
		"empty_for_lint.html",
	)
	mgx.Must(os.MkdirAll(filepath.Dir(emptyFile), os.ModePerm))
	f, err := os.OpenFile(
		emptyFile,
		os.O_RDONLY|os.O_CREATE,
		0644,
	)
	mgx.Must(err)
	mgx.Must(f.Close())
	defer func() {
		if err := sh.Rm(emptyFile); err != nil {
			fmt.Printf("ERROR deleting emtpyFile: %v\n", err)
		}
	}()

	_ = must.RunV(
		"go",
		"run",
		"github.com/golangci/golangci-lint/cmd/golangci-lint",
		"run",
		"--fix",
	)
	_ = must.RunV(
		"go",
		"run",
		"github.com/golangci/golangci-lint/cmd/golangci-lint",
		"run",
		"--fix",
		"--build-tags",
		"production",
	)
	_ = must.RunV(
		"go",
		"run",
		"github.com/golangci/golangci-lint/cmd/golangci-lint",
		"run",
		"--fix",
		"--build-tags",
		"mage",
	)
	_ = mustFrontend.RunV("npm", "run", "lint")
	_ = mustE2E.RunV("npm", "run", "lint")
}

// Build cleosrv and cleoc binaries for the current platform.
func Build() {
	_ = must.RunV(
		"go",
		"run",
		"github.com/goreleaser/goreleaser",
		"build",
		"--snapshot",
		"--single-target",
	)
}

// EnsureMage installs mage globally if it's not already installed.
func EnsureMage() error {
	return pkg.EnsureMage("")
}

// Test executes all tests except the e2e tests.
func Test() {
	_ = must.RunV("go", "test", "./...")
	_ = mustFrontend.RunV("npm", "test", "a", "--", "--watchAll=false")
}

// InstallDeps installs all dependencies.
func InstallDeps() {
	_ = must.RunV("go", "mod", "tidy")
	_ = must.RunV("go", "mod", "download")
	_ = mustFrontend.RunV("npm", "install")
	_ = mustE2E.RunV("npm", "install")
}

// MergeDependabot merges all open dependabot PRs.
func MergeDependabot() error {
	out, err := shx.Output("git", "rev-parse", "--abbrev-ref", "HEAD")
	mgx.Must(err)
	if out != "main" {
		return fmt.Errorf("Not on main! Exiting")
	}
	err = shx.Run("git", "diff", "--quiet", "--exit-code")
	if err != nil {
		return fmt.Errorf("There are uncommitted changes! Exiting")
	}

	_ = must.RunV("git", "fetch")
	_ = must.RunV("git", "remote", "prune", "origin")

	out, err = shx.Output(
		"git",
		"for-each-ref",
		"--format=%(refname)",
		"refs/remotes/origin/dependabot/",
	)
	mgx.Must(err)

	for _, pr := range strings.Split(out, "\n") {
		fmt.Printf("PR: %v\n", pr)
		err = shx.Run("git", "merge-base", "--is-ancestor", pr, "HEAD")
		if err == nil {
			fmt.Println("Already merged")
			continue
		}
		_ = must.RunV("git", "merge", pr, "-m", "Merge dependabot update")
		fmt.Println()
	}
	fmt.Println("All PRs merged")

	InstallDeps()
	Lint()
	Generate()

	err = shx.Run("git", "diff", "--exit-code")
	if err != nil {
		return fmt.Errorf("Code was changed via lint/generate: %w", err)
	}

	Test()
	E2ETest()

	fmt.Println("Successfully done. You must run 'git push' to publish the changes.")
	return nil
}

// E2ETest starts cleosrv and runs end-to-end tests with Cypress and Firefox.
func E2ETest() {
	e2eTestHelper("firefox", true, "")
}

// E2ETestB is like E2ETest but allows choosing a browser.
// Note that if the browser is 'electron' it will run headless.
func E2ETestB(browser string) {
	e2eTestHelper(browser, true, "")
}

// E2ETestC is like E2ETest but allow to configure multiple things.
// If the browser is 'electron' it will run headless. If the baseURL is
// specified (as opposed to an empty string) then no cleosrv is started, but
// it's assumed one is already running with that baseURL.
// shouldRebuild specifies whether the binaries should be re-built.
func E2ETestC(browser string, shouldRebuild bool, baseURL string) {
	e2eTestHelper(browser, shouldRebuild, baseURL)
}

func e2eTestHelper(browser string, shouldRebuild bool, baseURL string) {
	dbPath := "e2e_tests.db"
	mgx.Must(sh.Rm(dbPath))

	if shouldRebuild {
		mg.Deps(Clean)
		mg.Deps(Build)
	}
	buildTarget, err := getCurrentBuildTarget()
	mgx.Must(err)
	cleosrvPath, err := filepath.Abs(filepath.Join(
		"dist",
		"cleosrv_"+buildTarget,
		"cleosrv",
	))
	mgx.Must(err)
	cleocPath, err := filepath.Abs(filepath.Join(
		"dist",
		"cleoc_"+buildTarget,
		"cleoc",
	))
	mgx.Must(err)

	if baseURL == "" {
		cleosrvCmd := exec.Command(cleosrvPath, "--database", dbPath)
		var cleosrvStdout, cleosrvStderr bytes.Buffer
		cleosrvCmd.Stdout = &cleosrvStdout
		cleosrvCmd.Stderr = &cleosrvStderr
		mgx.Must(cleosrvCmd.Start())
		defer func() {
			if err := cleosrvCmd.Process.Kill(); err != nil {
				fmt.Printf("error stopping cleosrv: %v\n", err)
			}
			fmt.Println(strings.Repeat("=", 80))
			fmt.Println("cleosrv stdout:")
			fmt.Println(cleosrvStdout.String())
			fmt.Println(strings.Repeat("=", 80))
			fmt.Println("cleosrv stderr:")
			fmt.Println(cleosrvStderr.String())
			fmt.Println(strings.Repeat("=", 80))
		}()
	}
	cypressArgs := []string{
		"cypress",
		"run",
		"-b",
		browser,
	}
	if strings.ToLower(browser) == "electron" {
		cypressArgs = append(cypressArgs, "--headless")
	} else {
		cypressArgs = append(cypressArgs, "--headed")
	}

	var cypressConfig []string
	var cypressEnv []string

	if baseURL == "" {
		cypressConfig = append(cypressConfig, "baseUrl=http://localhost:8080")
	} else {
		cypressConfig = append(cypressConfig, "baseUrl="+baseURL)
	}
	cypressEnv = append(cypressEnv, "cleocPath="+cleocPath)

	if len(cypressConfig) > 0 {
		cypressArgs = append(cypressArgs, "--config", strings.Join(cypressConfig, ","))
	}
	if len(cypressEnv) > 0 {
		cypressArgs = append(cypressArgs, "--env", strings.Join(cypressEnv, ","))
	}
	_ = shx.Command("npx", cypressArgs...).In("e2e_tests").Must().RunV()
}

// getCurrentBuildTarget is a helper function to determine where the compiled
// binaries have been placed by goreleaser
func getCurrentBuildTarget() (string, error) {
	goos := os.Getenv("GOOS")
	if goos == "" {
		goos = runtime.GOOS
	}
	goarch := os.Getenv("GOARCH")
	if goarch == "" {
		goarch = runtime.GOARCH
	}
	if goarch == "amd64" {
		return fmt.Sprintf("%s_%s_v1", goos, goarch), nil
	} else if goarch == "arm64" {
		return fmt.Sprintf("%s_%s", goos, goarch), nil
	}
	return "", fmt.Errorf("unknown goarch: %v", goarch)
}

// DeployDemo deploy and overwrite demo.cleodora.org
func DeployDemo() error {
	buildTarget, err := getCurrentBuildTarget()
	mgx.Must(err)
	if buildTarget != "linux_amd64_v1" {
		return fmt.Errorf(
			"the Dockerfile is hardcoded to use 'linux_amd64_v1' and not '%v'",
			buildTarget,
		)
	}
	mg.Deps(Clean)
	mg.Deps(InstallDeps)
	mg.Deps(Lint)
	mg.Deps(Generate)
	err = shx.Run("git", "diff", "--quiet", "--exit-code")
	mg.Deps(Test)
	mg.Deps(E2ETest)
	if err != nil {
		return fmt.Errorf("There are uncommitted changes! Exiting")
	}
	mg.Deps(Build)
	_ = must.RunV(
		"flyctl",
		"deploy",
		"--local-only", // use local Docker to build
	)
	fmt.Println("Sleeping 30s so demo.cleodora.org comes up")
	// sometimes it takes a little until the new instance is fully active
	time.Sleep(30 * time.Second)
	_ = must.RunV("./scripts/demoDummyData.sh")
	return nil
}
