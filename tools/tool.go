package tool

import "os"
import "path/filepath"
import . "github.com/moovweb/gpkg/errors"

type ToolError struct{ msg string }

func NewToolError(msg string) *ToolError { return &ToolError{msg: msg} }
func (e *ToolError) String() string      { return "Tool Error: " + e.msg }

type Tool interface {
	Clean() (string, Error)
	Build() (string, Error)
	Test() (string, Error)
	Install() (string, Error)
}

func NewTool(path string) Tool {
	_, err := os.Open(filepath.Join(path, "Makefile.gvm"))
	if err == nil {
		return Tool(NewMakeTool(path, "Makefile.gvm"))
	}

	return Tool(NewGbTool(path))
}
