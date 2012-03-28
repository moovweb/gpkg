package tool

import "os"
import "os/exec"
import "path/filepath"

type GoinstallTool struct {
	gopath string
	target string
}

func NewGoinstallTool(gopath string, target string) Tool {
	gopath, err := filepath.Abs(gopath)
	if err != nil {
		return nil
	}
	return Tool(GoinstallTool{gopath: gopath, target: target})
}

func (g GoinstallTool) Clean() (string, error) {
	return "", nil
}

func (g GoinstallTool) Build() (string, error) {
	pushd := os.Getenv("GOPATH")
	os.Setenv("GOPATH", g.gopath)
	out, err := exec.Command("goinstall", g.target).CombinedOutput()
	if err != nil {
		return "", NewToolError(string(out))
	}
	os.Setenv("GOPATH", pushd)
	return string(out), nil
}

func (g GoinstallTool) Test() (string, error) {
	return "", nil
}

func (g GoinstallTool) Install() (string, error) {
	return "", nil
}
