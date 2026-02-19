package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {

}

func GetSize(path string) (int64, string, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return 0, path, err
	}
	if fi.IsDir() == false {
		return fi.Size(), path, nil
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		files, err := os.ReadDir(path)
		if err != nil {
			return 0, path, err
		}
		var value int64
		for _, file := range files {
			fullPath := filepath.Join(path, file.Name())
			newfi, err := os.Lstat(fullPath)
			if err != nil {
				fmt.Printf("Ошибка при обработке %s: %v", fullPath, err)
				continue
			}
			if newfi.Mode()&os.ModeSymlink == 0 {
				if newfi.IsDir() == false {
					value += newfi.Size()
				}

			}
		}
		return value, path, nil
	}
	target, err := os.Readlink(path)
	if err != nil {
		return 0, path, err
	}
	targetFi, err := os.Stat(target)
	if err != nil {
		return 0, path, err
	}
	if targetFi.IsDir() == false {
		return targetFi.Size(), os.Readlink(path), nil
	}
	tfile, err := os.ReadDir(target)
	if err != nil {
		return 0, path, err
	}
	var tvalue int64
	for _, file2 := range tfile {
		fullPath2 := filepath.Join(path, file2.Name())
		newtfi, err := os.Lstat(fullPath2)
		if err != nil {
			fmt.Printf("Ошибка при обработке %s: %v", fullPath, err)
			continue
		}
		if newtfi.IsDir() == false {
			tvalue += newtfi.Size()
		}

	}
	return tvalue, target, nil
}
