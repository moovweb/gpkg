package util

import "os"
import "exec"
import "strings"

func FileCopy(src string, dst string) os.Error {
	out, err := exec.Command("cp", "-r", src, dst).CombinedOutput()
	if err != nil {
		return os.NewError(strings.TrimSpace(string(out)))
	}
	return nil
}
