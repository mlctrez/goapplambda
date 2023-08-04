package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"os"
	"reflect"
	"runtime"
	"strings"
)

var Default = Run

type buildStep func() error

const WasmPath = "web/app.wasm"

func Clean() (err error) { return sh.Run("rm", "-rf", "temp", WasmPath) }

func Build() (err error) {

	buildSteps := []buildStep{
		mkTemp, modDownload, parseGitInfo, findModuleName,
		buildWasm, buildBinary,
	}
	for i, step := range buildSteps {
		name := runtime.FuncForPC(reflect.ValueOf(step).Pointer()).Name()
		name = strings.TrimPrefix(name, "main.")
		fmt.Printf("%02d %s\n", i, name)
		if err = step(); err != nil {
			return
		}
	}

	return nil
}

func Run() error {
	if err := Build(); err != nil {
		return err
	}
	env := map[string]string{"DEV": "1"}

	return sh.RunWith(env, "temp/goapp")
}

func Deploy() (err error) {
	if err = Build(); err != nil {
		return
	}
	if err = sh.Run("zip", "-j", "temp/goapp.zip", "temp/goapp"); err != nil {
		return
	}

	s3Args := []string{"s3", "sync", "web", "s3://mlctrez-goapplambda/web"}
	if err = sh.Run("aws", s3Args...); err != nil {
		return
	}

	// E1T21UEDW4RGGJ

	if err = sh.Run("aws", "cloudfront", "create-invalidation",
		"--distribution-id", "E1T21UEDW4RGGJ",
		"--paths", "/*",
	); err != nil {
		return
	}

	lambdaArgs := []string{
		"lambda", "update-function-code", "--function-name", "goapplambda",
		"--zip-file", "fileb://temp/goapp.zip", "--output", "text",
	}

	// 	aws lambda update-function-code --function-name goapplambda --zip-file fileb://temp/goapp.zip --output text
	if err = sh.Run("aws", lambdaArgs...); err != nil {
		return
	}

	return nil
}

func mkTemp() error      { return os.MkdirAll("temp", 0755) }
func modDownload() error { return sh.Run("go", "mod", "download") }

type GitInfo struct {
	Version string
	Commit  string
}

var gitInfo GitInfo

func parseGitInfo() error {
	gitInfo.Version = "v0.0.0"
	gitInfo.Commit = "HEAD"
	_, err := os.Stat(".git")
	if os.IsNotExist(err) {
		return nil
	}
	var output string
	output, err = sh.Output("git", "describe", "--abbrev=0", "--tags")
	if err == nil {
		gitInfo.Version = output
	}
	output, err = sh.Output("git", "rev-parse", "--short", "HEAD")
	if err == nil {
		gitInfo.Commit = output
	}
	return nil
}

var moduleName string

func findModuleName() (err error) {
	moduleName, err = sh.Output("go", "list", "-m")
	return
}

func buildWasm() (err error) {

	if err = sh.Run("rm", "-rf", WasmPath); err != nil {
		return
	}
	env := map[string]string{"GOARCH": "wasm", "GOOS": "js"}
	// -ldflags

	var ldFlags string
	if ldFlags, err = buildLdFlags(false); err != nil {
		return
	}

	return sh.RunWith(env, "go", "build",
		"-o", WasmPath,
		"-ldflags", ldFlags,
		"goapp/service/main/main.go")

}

func buildLdFlags(withSize bool) (string, error) {
	var ldFlags string
	ldFlags += "-w"
	ldFlags += fmt.Sprintf(" -X %s/goapp.Version=%s", moduleName, gitInfo.Version)
	ldFlags += fmt.Sprintf(" -X %s/goapp.Commit=%s", moduleName, gitInfo.Commit)
	if withSize {
		stat, err := os.Stat(WasmPath)
		if err != nil {
			return "", err
		}
		ldFlags += fmt.Sprintf(" -X %s/goapp.WasmSize=%d", moduleName, stat.Size())
	}
	//fmt.Println(ldFlags)
	return ldFlags, nil
}

func buildBinary() error {
	ldFlags, err := buildLdFlags(true)
	if err != nil {
		return err
	}
	env := map[string]string{"CGO_ENABLED": "0"}
	return sh.RunWith(env, "go", "build",
		"-o", "temp/goapp",
		"-ldflags", ldFlags,
		"goapp/service/main/main.go")
}
