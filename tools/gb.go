package tool

import "os"
import "exec"
import . "errors"

type GbTool struct {
	sandbox string
}

func NewGbTool(sandbox string) Tool {
	return Tool(GbTool{sandbox: sandbox})
}

func (gb GbTool) runCommand(cmd string) (string, Error) {
	pushd, err := os.Getwd()
	if err != nil {
		return "", NewToolError("Failed to get working directory")
	}
	err = os.Chdir(gb.sandbox)
	if err != nil {
		return "", NewToolError("Failed to chdir " + gb.sandbox)
	}
	out, err := exec.Command("gb", cmd).CombinedOutput()
	if err != nil {
		return "", NewToolError(string(out))
	}
	err = os.Chdir(pushd)
	if err != nil {
		return "", NewToolError("Failed to chdir " + pushd)
	}
	return string(out), nil
}

func (gb GbTool) Clean() (string, Error) {
	return gb.runCommand("-c")
}

func (gb GbTool) Build() (string, Error) {
	return gb.runCommand("-b")
}

func (gb GbTool) Test() (string, Error) {
	return gb.runCommand("-t")
}

func (gb GbTool) Install() (string, Error) {
	return gb.runCommand("-i")
}
