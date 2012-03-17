package tool

import "os"
import "exec"

type GBTool struct {
	sandbox string
}

func NewGBTool(sandbox string) Tool {
	return Tool(GBTool{sandbox:sandbox})
}

func (gb *GBTool) runCommand(cmd string) (string, *ToolError) {
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

func (gb GBTool) Clean() (string, *ToolError) {
	return gb.runCommand("-c")
}

func (gb GBTool) Build() (string, *ToolError) {
	return gb.runCommand("-b")
}

func (gb GBTool) Test() (string, *ToolError) {
	return gb.runCommand("-t")
}

func (gb GBTool) Install() (string, *ToolError) {
	return gb.runCommand("-i")
}

