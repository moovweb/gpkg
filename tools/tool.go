package tool

type ToolError struct{ msg string }

func NewToolError(msg string) *ToolError { return &ToolError{msg: msg} }
func (e *ToolError) String() string      { return "Tool Error: " + e.msg }

type Tool interface {
	Clean() (string, *ToolError)
	Build() (string, *ToolError)
	Test() (string, *ToolError)
	Install() (string, *ToolError)
}
