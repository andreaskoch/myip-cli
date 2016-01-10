// +build ignore

// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program builds myip.
//
// $ go run make.go -install
//
// View the README.md for further details.
//
// The output binaries go into the ./bin/ directory (under the GOPATH, where make.go is)
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ProjectName contains the name of the project
const ProjectName = "myip"

// GOPATH environment variable name.
const GOPATH_ENVIRONMENT_VARIABLE = "GOPATH"

// GOBIN environment variable name
const GOBIN_ENVIRONMENT_VARIBALE = "GOBIN"

// GOOS environment variable name
const GOOS_ENVIRONMENT_VARIBALE = "GOOS"

// GOARCH environment variable name
const GOARCH_ENVIRONMENT_VARIBALE = "GOARCH"

var (

	// command line flags
	verboseFlagIsSet      = flag.Bool("v", false, "Verbose mode")
	installFlagIsSet      = flag.Bool("install", false, fmt.Sprintf("Compiles the %s binaries", ProjectName))
	crossCompileFlagIsSet = flag.Bool("crosscompile", false, fmt.Sprintf("Compile %s binaries for all platforms and architectures", ProjectName))
	versionFlagIsSet      = flag.Bool("version", false, "Get the current version number of the repository")

	// The GOPATH for the current project
	goPath = getWorkingDirectory()

	// The GOBIN for the current project
	goBin = filepath.Join(goPath, "bin")

	// A list of all supported compilation targets (e.g. "windows/amd64")
	compilationTargets = []compilationTarget{
		compilationTarget{"darwin", "amd64", []string{}},
		compilationTarget{"dragonfly", "amd64", []string{}},
		compilationTarget{"freebsd", "amd64", []string{}},
		compilationTarget{"freebsd", "arm", []string{}},
		compilationTarget{"freebsd", "arm", []string{"GOARM=5"}},
		compilationTarget{"freebsd", "arm", []string{"GOARM=6"}},
		compilationTarget{"freebsd", "arm", []string{"GOARM=7"}},
		compilationTarget{"linux", "amd64", []string{}},
		compilationTarget{"linux", "arm", []string{}},
		compilationTarget{"linux", "arm", []string{"GOARM=5"}},
		compilationTarget{"linux", "arm", []string{"GOARM=6"}},
		compilationTarget{"linux", "arm", []string{"GOARM=7"}},
		compilationTarget{"netbsd", "amd64", []string{}},
		compilationTarget{"netbsd", "arm", []string{"GOARM=5"}},
		compilationTarget{"netbsd", "arm", []string{"GOARM=6"}},
		compilationTarget{"netbsd", "arm", []string{"GOARM=7"}},
		compilationTarget{"openbsd", "amd64", []string{}},
		compilationTarget{"solaris", "amd64", []string{}},
		compilationTarget{"windows", "amd64", []string{}},
	}

	// The current version number (e.g. "2015-01-11-284c030+")
	version = gitVersion()
)

// Compilation Target Definition
type compilationTarget struct {
	OperatingSystem string
	Architecture    string
	OtherVariables  []string
}

func (target *compilationTarget) String() string {
	return fmt.Sprintf("%s/%s", target.OperatingSystem, target.Architecture)
}

func init() {

	executableName := "make.go"

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s provides functions for compiling %s.\n", executableName, ProjectName)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Version: %s\n", version)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  go run %s [options]\n", executableName)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	if *verboseFlagIsSet {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("%s: %s\n", GOPATH_ENVIRONMENT_VARIABLE, goPath)
		fmt.Printf("%s: %s\n", GOBIN_ENVIRONMENT_VARIBALE, goBin)
	}

	if *installFlagIsSet {
		if err := install(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		return
	}

	if *crossCompileFlagIsSet {
		if err := crossCompile(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		return
	}

	if *versionFlagIsSet {
		printProjectVersionNumber()
		return
	}

	flag.Usage()
}

// Install all parts of allmark using go install.
func install() error {

	// prepare the environment variables
	environmentVariables := cleanGoEnv()
	environmentVariables = setEnv(environmentVariables, GOPATH_ENVIRONMENT_VARIABLE, goPath)
	environmentVariables = setEnv(environmentVariables, GOBIN_ENVIRONMENT_VARIBALE, goBin)

	return runCommand(os.Stdout, os.Stderr, goPath, environmentVariables, "go", "install", getBuildVersionFlag())

}

// Cross-compile all parts of allmark for all supported platforms.
func crossCompile() error {

	// prepare the environment variables
	environmentVariables := cleanGoEnv()
	environmentVariables = setEnv(environmentVariables, GOPATH_ENVIRONMENT_VARIABLE, goPath)
	// environmentVariables = setEnv(environmentVariables, GOBIN_ENVIRONMENT_VARIBALE, goBin)

	// iterate over all supported compilation targets
	for _, target := range compilationTargets {

		// assemble the target path
		targetFile := filepath.Join(goBin, fmt.Sprintf("%s_%s_%s", ProjectName, target.OperatingSystem, target.Architecture))

		// prepare environment variables for cross-compilation
		crossCompileEnvironemntVariables := environmentVariables
		crossCompileEnvironemntVariables = setEnv(crossCompileEnvironemntVariables, GOOS_ENVIRONMENT_VARIBALE, target.OperatingSystem)
		crossCompileEnvironemntVariables = setEnv(crossCompileEnvironemntVariables, GOARCH_ENVIRONMENT_VARIBALE, target.Architecture)

		// add additional environment variables
		for _, additionalEnvVariable := range target.OtherVariables {
			components := strings.Split(additionalEnvVariable, "=")
			name := components[0]
			value := components[1]

			// append additional environment variables to target file name
			targetFile += fmt.Sprintf("_%s_%s", strings.ToLower(name), strings.ToLower(value))

			crossCompileEnvironemntVariables = setEnv(crossCompileEnvironemntVariables, name, value)
		}

		// build the package for the specified os and arch
		if *verboseFlagIsSet {
			fmt.Printf("Compiling %s for %s\n", ProjectName, target.String())
		}

		// add .exe extension for windows
		if target.OperatingSystem == "windows" {
			targetFile += ".exe"
		}

		err := runCommand(os.Stdout,
			os.Stderr,
			goPath,
			crossCompileEnvironemntVariables,
			"go",
			"build",
			"-o",
			targetFile,
			"-x",
			getBuildVersionFlag())

		if err != nil {
			return err
		}

	}

	return nil
}

// Print the current version number of the project.
func printProjectVersionNumber() {
	fmt.Fprintf(os.Stdout, gitVersion())
}

// getWorkingDirectory returns the current working directory path or fails.
func getWorkingDirectory() string {
	goPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	return goPath
}

// Execute go in the specified go path with the supplied command arguments.
func runCommand(stdout, stderr io.Writer, workingDirectory string, environmentVariables []string, command string, args ...string) error {

	// Create the command
	cmdName := fmt.Sprintf("%s %s", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)

	// Set the working directory
	cmd.Dir = workingDirectory

	// set environment variables
	cmd.Env = environmentVariables

	// Capture the output
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if *verboseFlagIsSet {
		log.Printf("Running %s", cmdName)
	}

	// execute the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error running %s: %v", cmdName, err.Error())
	}

	return nil
}

// cleanGoEnv returns a copy of the current environment with GOPATH_ENVIRONMENT_VARIABLE and GOBIN_ENVIRONMENT_VARIBALE removed.
func cleanGoEnv() (clean []string) {
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, GOPATH_ENVIRONMENT_VARIABLE+"=") || strings.HasPrefix(env, GOBIN_ENVIRONMENT_VARIBALE+"=") {
			continue
		}

		clean = append(clean, env)
	}

	return
}

// setEnv sets the given key & value in the provided environment.
// Each value in the env list should be of the form key=value.
func setEnv(env []string, key, value string) []string {
	for i, s := range env {
		if strings.HasPrefix(s, fmt.Sprintf("%s=", key)) {
			env[i] = envPair(key, value)
			return env
		}
	}
	env = append(env, envPair(key, value))
	return env
}

// Create an environment variable of the form key=value.
func envPair(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}

// gitVersion returns the git version of the git repo at goPath as a
// string of the form "yyyy-mm-dd-xxxxxxx", with an optional trailing
// '+' if there are any local uncomitted modifications to the tree.
func gitVersion() string {

	// get the latest git hash
	gitHeadCommand := exec.Command("git", "rev-parse", "--short", "HEAD")
	gitHeadCommand.Dir = goPath

	gitHeadCommandOutput, err := gitHeadCommand.Output()
	if err != nil {
		log.Fatalf("Error while retrieving the latest git hash: %v", err)
	}

	gitHash := strings.TrimSpace(string(gitHeadCommandOutput))

	// check if the working directory is dirty
	dirtyFlag := ""
	gitDiffCommand := exec.Command("git", "diff", "--exit-code")
	gitDiffCommand.Dir = goPath
	if err := gitDiffCommand.Run(); err != nil {
		dirtyFlag = "+"
	}

	// get the latest tag
	gitLatestTagCommand := exec.Command("git", "describe", "--abbrev=0", "--tags")
	gitLatestTagCommand.Dir = goPath

	latestTag := "v0.0.0"
	gitLatestTagCommandOutput, _ := gitLatestTagCommand.Output()
	latestTag = strings.TrimSpace(string(gitLatestTagCommandOutput))
	date := getCurrentDateString()

	versionNumber := fmt.Sprintf("%s-%s-%s%s", latestTag, date, gitHash, dirtyFlag)

	return versionNumber
}

// Get the build version flag for the go linker (e.g. -X allmark.io/cmd/allmark 2015-01-11-284c030+).
func getBuildVersionFlag() string {
	return fmt.Sprintf("--ldflags=-X %s.GitInfo=%s", "main", version)
}

// getCurrentDateString returns the current date as a formatted string (e.g. "2006-01-02").
func getCurrentDateString() string {
	return time.Now().Format("2006-01-02")
}
