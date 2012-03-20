package tool

import "os"
import "exec"
import . "github.com/moovweb/gpkg/errors"

type MakeTool struct {
	sandbox  string
	filename string
}

func NewMakeTool(sandbox string, filename string) Tool {
	return Tool(MakeTool{sandbox: sandbox, filename: filename})
}

func (m MakeTool) runCommand(cmd string) (string, Error) {
	pushd, err := os.Getwd()
	if err != nil {
		return "", NewToolError("Failed to get working directory")
	}
	err = os.Chdir(m.sandbox)
	if err != nil {
		return "", NewToolError("Failed to chdir " + m.sandbox)
	}
	out, err := exec.Command("make", "-f", m.filename, cmd).CombinedOutput()
	if err != nil {
		return "", NewToolError(err.String() + string(out))
	}
	err = os.Chdir(pushd)
	if err != nil {
		return "", NewToolError("Failed to chdir " + pushd)
	}
	return string(out), nil
}

func (m MakeTool) Clean() (string, Error) {
	return m.runCommand("clean")
}

func (m MakeTool) Build() (string, Error) {
	return m.runCommand("build")
}

func (m MakeTool) Test() (string, Error) {
	return m.runCommand("test")
}

func (m MakeTool) Install() (string, Error) {
	return m.runCommand("install")
}
