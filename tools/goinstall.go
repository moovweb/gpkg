package tool

import "os"
import "exec"
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
	return Tool(GoinstallTool{gopath:gopath, target:target})
}

func (g GoinstallTool) Clean() (string, *ToolError) {
	return "", nil
}

func (g GoinstallTool) Build() (string, *ToolError) {
	pushd := os.Getenv("GOPATH")
	os.Setenv("GOPATH", g.gopath)
	out, err := exec.Command("goinstall", g.target).CombinedOutput()
	if err != nil {
		return "", NewToolError(string(out))
	}
	os.Setenv("GOPATH", pushd)
	return string(out), nil
}

func (g GoinstallTool) Test() (string, *ToolError) {
	return "", nil
}

func (g GoinstallTool) Install() (string, *ToolError) {
	return "", nil
}

