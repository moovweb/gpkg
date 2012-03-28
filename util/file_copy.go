package util

import (
	"errors"
	"os/exec"
)
import "strings"

func FileCopy(src string, dst string) error {
	out, err := exec.Command("cp", "-rf", src, dst).CombinedOutput()
	if err != nil {
		return errors.New(strings.TrimSpace(string(out)))
	}
	return nil
}
