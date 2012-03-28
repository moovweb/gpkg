package tool

import "os"
import "path/filepath"

type ToolError struct{ msg string }

func NewToolError(msg string) *ToolError { return &ToolError{msg: msg} }
func (e *ToolError) Error() string      { return "Tool Error: " + e.msg }

type Tool interface {
	Clean() (string, error)
	Build() (string, error)
	Test() (string, error)
	Install() (string, error)
}

func NewTool(path string) (Tool, error) {
	_, err := os.Open(filepath.Join(path, "Makefile.gvm"))
	if err == nil {
		return Tool(NewMakeTool(path, "Makefile.gvm")), nil
	}

	return Tool(NewGbTool(path)), nil
}
