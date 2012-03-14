package gpkg

import "exec"
import "os"

func FileCopy(src string, dst string) (err os.Error) {
	_, err = exec.Command("cp", "-r", src, dst).CombinedOutput()
	return
}

