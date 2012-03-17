package tool

import "testing"

const TEST_TOOL_PROJECT_GB = "testdata/gb_project"
const TEST_TOOL_PROJECT_MAKE = "testdata/make_project"
const TEST_TOOL_PROJECT_GOINSTALL = "testdata/goinstall_project"
const TEST_TOOL_PROJECT_GOINSTALL_TARGET = "testdata/goinstall_project"

func testTool(tool Tool) (string, *ToolError) {
	out, err := tool.Clean()
	if err != nil {
		return out, err
	}
	out, err = tool.Build()
	if err != nil {
		return out, err
	}
	out, err = tool.Clean()
	if err != nil {
		return out, err
	}
	return "", nil
}

func TestNewTool(t *testing.T) {
	tool := NewGBTool(TEST_TOOL_PROJECT_GB)
	if tool == nil {
		t.Fatal("Failed to create tool for", TEST_TOOL_PROJECT_GB)
	}
	_, err := testTool(tool)
	if err != nil {
		t.Error(err)
		t.Fatal("GB Test failed")
	}

	tool = NewMakeTool(TEST_TOOL_PROJECT_MAKE, "Makefile.gvm")
	if tool == nil {
		t.Fatal("Failed to create tool for", TEST_TOOL_PROJECT_MAKE)
	}
	_, err = testTool(tool)
	if err != nil {
		t.Error(err)
		t.Fatal("Make Test failed")
	}

	tool = NewGoinstallTool(TEST_TOOL_PROJECT_GOINSTALL, "project")
	if tool == nil {
		t.Fatal("Failed to create tool for", TEST_TOOL_PROJECT_MAKE)
	}
	_, err = testTool(tool)
	if err != nil {
		t.Error(err)
		t.Fatal("Make Test failed")
	}
}
