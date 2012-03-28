package tool

import "os"
import "os/exec"

type MakeTool struct {
	sandbox  string
	filename string
}

func NewMakeTool(sandbox string, filename string) Tool {
	return Tool(MakeTool{sandbox: sandbox, filename: filename})
}

func (m MakeTool) runCommand(cmd string) (string, error) {
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
		return "", NewToolError(err.Error() + string(out))
	}
	err = os.Chdir(pushd)
	if err != nil {
		return "", NewToolError("Failed to chdir " + pushd)
	}
	return string(out), nil
}

func (m MakeTool) Clean() (string, error) {
	return m.runCommand("clean")
}

func (m MakeTool) Build() (string, error) {
	return m.runCommand("build")
}

func (m MakeTool) Test() (string, error) {
	return m.runCommand("test")
}

func (m MakeTool) Install() (string, error) {
	return m.runCommand("install")
}
