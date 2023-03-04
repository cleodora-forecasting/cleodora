//go:build mage

package main

import (
	"bufio"
	"bytes"
	"errors"
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
	ensureOnMainBranch()
	ensureGitDiffEmpty()

	_ = must.RunV("git", "fetch")
	_ = must.RunV("git", "remote", "prune", "origin")

	out, err := shx.Output(
		"git",
		"for-each-ref",
		"--format=%(refname)",
		"refs/remotes/origin/dependabot/",
	)
	mgx.Must(err)
	if out == "" {
		fmt.Println("Nothing to merge")
		return nil
	}

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
	mg.Deps(All)
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

// DeployDemo deploy and overwrite demo.cleodora.org .
func DeployDemo() error {
	buildTarget, err := getCurrentBuildTarget()
	mgx.Must(err)
	if buildTarget != "linux_amd64_v1" {
		return fmt.Errorf(
			"the Dockerfile is hardcoded to use 'linux_amd64_v1' and not '%v'",
			buildTarget,
		)
	}
	mg.Deps(All)
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

// All includes the most important targets such as Lint, Test, E2ETest.
func All() {
	mg.Deps(Clean)
	mg.Deps(InstallDeps)
	mg.Deps(Lint)
	mg.Deps(Generate)

	// run it once here to exit early because the tests can take a while
	ensureGitDiffEmpty()

	mg.Deps(Test)
	mg.Deps(E2ETest)
	mg.Deps(Build)

	ensureGitDiffEmpty()
}

// Release builds and releases all packages
func Release() error {
	ensureGitDiffEmpty()

	// Changelog
	changelogPath := "temp_changelog.md"
	if _, err := os.Stat(changelogPath); errors.Is(err, os.ErrNotExist) {
		msg := `'%v' does not exist. You should:
* Update:
    vim website/content/docs/changelog.md
* Commit the changes
* Execute:
    cp website/content/docs/changelog.md temp_changelog.md
    vim temp_changelog.md
* Remove everything except the current release from that file. Also
  remove the version title because it becomes redundant`
		return fmt.Errorf(msg, changelogPath)
	}
	changelogContent, err := os.ReadFile(changelogPath)
	mgx.Must(err)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Print(string(changelogContent))
	fmt.Println(strings.Repeat("=", 80))
	if !ask("Is the changelog correct?", false) {
		return fmt.Errorf("'%v' is not ready", changelogPath)
	}

	// Download links
	if !ask(
		"Are the download links and 'latest release' on the website up-to-date?",
		false,
	) {
		msg := `the download links are not ready. You should:
* Update:
    vim website/content/docs/user/*
* Commit the changes`
		return errors.New(msg)
	}

	// GITHUB_TOKEN
	if _, present := os.LookupEnv("GITHUB_TOKEN"); !present {
		msg := `GITHUB_TOKEN for goreleaser is missing
Instructions to create it:
  * https://github.com/settings/personal-access-tokens/new
  * Access on the cleodora-forecasting organization and repository
    cleodora-forecasting/cleodora
  * Give the token no organization permissions and the following repository
    permissions:
    * Read access to metadata
    * Read and Write access to code
Then set is as an ENV variable:
    GITHUB_TOKEN=asdf mage release`
		return errors.New(msg)
	}

	// Ensure tests run etc.
	mg.Deps(All)

	// Tag ready and checked out
	out, err := shx.Output("git", "tag", "--points-at", "HEAD")
	mgx.Must(err)
	if out == "" {
		msg := `no tag is checked out. You should:
* Create a tag:
    git tag vX.Y.Z
* OR check out the correct tag (if changelog etc. is ready there):
    git checkout vX.Y.Z`
		mgx.Must(errors.New(msg))
	}

	// Explicitly call Clean() because it's not enough that it was already
	// called as an indirect dependency earlier because goreleaser expects
	// dist/ to be empty.
	Clean()

	_ = must.RunV(
		"go",
		"run",
		"github.com/goreleaser/goreleaser",
		"release",
		"--release-notes",
		changelogPath,
	)

	mgx.Must(sh.Rm(changelogPath))
	return nil
}

func ensureGitDiffEmpty() {
	err := shx.Run("git", "diff", "--quiet", "--exit-code")
	if err != nil {
		mgx.Must(errors.New("There are uncommitted changes! Exiting"))
	}
}

func ensureOnMainBranch() {
	out, err := shx.Output("git", "rev-parse", "--abbrev-ref", "HEAD")
	mgx.Must(err)
	if out != "main" {
		mgx.Must(errors.New("Not on main! Exiting"))
	}
}

// ask a yes/no question.
func ask(question string, defaultYes bool) bool {
	choices := "Y/n"
	if !defaultYes {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Printf("%s (%s) ", question, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return defaultYes
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}
