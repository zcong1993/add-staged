package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	version = "master"
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(buildVersion(version, commit, date, builtBy))
	}

	app := &cli.App{
		Name:    "add-staged",
		Usage:   "git add only staged files",
		Action:  addStaged,
		Version: version,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func addStaged(c *cli.Context) error {
	files, err := getStagedFiles(mustCwd())
	if err != nil {
		return fmt.Errorf("get staged files err %v", err)
	}

	args := []string{"add", "--"}
	args = append(args, files...)

	return exec.Command("git", args...).Run()
}

func mustCwd() string {
	cwd, _ := os.Getwd()
	return cwd
}

func getStagedFiles(cwd string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--staged", "--diff-filter=ACMR", "--name-only", "-z")
	cmd.Dir = cwd
	resp, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	r1 := regexp.MustCompile("\u0000$")
	resp = r1.ReplaceAll(resp, []byte(""))

	files := strings.Split(string(resp), "\u0000")
	nfs := make([]string, len(files))
	for i, f := range files {
		nfs[i] = filepath.Join(cwd, f)
	}

	return nfs, nil
}

func buildVersion(version, commit, date, builtBy string) string {
	result := version

	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}

	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}

	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}

	return result
}
