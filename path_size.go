package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "filesize",
		Usage: "Memory",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "human",
				Aliases: []string{"h"},
				Usage:   "human-readable sizes (auto-select unit)",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return fmt.Errorf("не указан путь к файлу или директории")
			}

			path := c.Args().First()
			humanReadable := c.Bool("human")

			size, resolvedPath := GetSize(path)
			fmt.Printf("%s\t%s\n", FormatSize(size, humanReadable), resolvedPath)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}

func GetSize(path string) (int64, string) {
	fi, err := os.Lstat(path)
	if err != nil {
		return 0, path
	}
	if !fi.IsDir() {
		return fi.Size(), path
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		files, err := os.ReadDir(path)
		if err != nil {
			return 0, path
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
				if !newfi.IsDir() {
					value += newfi.Size()
				}

			}
		}
		return value, path
	}
	target, err := os.Readlink(path)
	if err != nil {
		return 0, path
	}
	return GetSize(target)

}

func FormatSize(size int64, human bool) string {
	if !human {
		return strconv.FormatInt(size, 10) + "B"
	}
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.1fTB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.1fGB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.1fMB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.1fKB", float64(size)/KB)
	default:
		return strconv.FormatInt(size, 10) + "B"
	}
}
