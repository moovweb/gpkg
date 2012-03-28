package tool

import "os"
import "os/exec"

type GbTool struct {
	sandbox string
}

func NewGbTool(sandbox string) Tool {
	return Tool(GbTool{sandbox: sandbox})
}

func (gb GbTool) runCommand(cmd string) (string, error) {
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

func (gb GbTool) Clean() (string, error) {
	return gb.runCommand("-c")
}

func (gb GbTool) Build() (string, error) {
	return gb.runCommand("-b")
}

func (gb GbTool) Test() (string, error) {
	return gb.runCommand("-t")
}

func (gb GbTool) Install() (string, error) {
	return gb.runCommand("-i")
}
