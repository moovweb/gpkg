package tool

import "os"
import "path/filepath"
import . "errors"

type ToolError struct{ msg string }

func NewToolError(msg string) *ToolError { return &ToolError{msg: msg} }
func (e *ToolError) String() string      { return "Tool Error: " + e.msg }

type Tool interface {
	Clean() (string, Error)
	Build() (string, Error)
	Test() (string, Error)
	Install() (string, Error)
}

func NewTool(path string) (Tool, Error) {
	_, err := os.Open(filepath.Join(path, "Makefile.gvm"))
	if err == nil {
		return Tool(NewMakeTool(path, "Makefile.gvm")), nil
	}

	return Tool(NewGbTool(path)), nil
}
