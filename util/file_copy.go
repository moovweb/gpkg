package util

import (
	"os"
	"io"
	"io/ioutil"
	"path/filepath"
)

func copyFile(dst, src string) (int64, error) { 
        sf, err := os.Open(src) 
        if err != nil { 
                return 0, err 
        } 
        defer sf.Close() 
        df, err := os.Create(dst) 
        if err != nil { 
                return 0, err 
        } 
        defer df.Close() 
        return io.Copy(df, sf) 
} 

func copyDir(dst, src string) error {
	dirs, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			err := os.MkdirAll(filepath.Join(dst, dir.Name()), 0755)
			if err != nil {
				return err
			}
			err = copyDir(filepath.Join(dst, dir.Name()), filepath.Join(src, dir.Name()))
			if err != nil {
				return err
			}
		} else {
			_, err := copyFile(filepath.Join(dst, dir.Name()), filepath.Join(src, dir.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func FileCopy(src string, dst string) error {
	err := copyDir(dst, src)
	if err != nil {
		return err
	}
	return nil
}
