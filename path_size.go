package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory; supports -r (recursive), -H (human-readable), -a (include hidden)",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "recursive",
				Aliases: []string{"r"},
				Usage:   "recursive size of directories",
			},
			&cli.BoolFlag{
				Name:    "human",
				Aliases: []string{"H"},
				Usage:   "human-readable sizes (auto-select unit)",
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "include hidden files and directories",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return fmt.Errorf("не указан путь к файлу или директории")
			}

			path := c.Args().First()
			humanReadable := c.Bool("human")
			showAll := c.Bool("all")
			recursive := c.Bool("recursive")

			size, resolvedPath := GetSize(path, showAll, recursive)
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

func GetSize(path string, showAll bool, recursive bool) (int64, string) {
	fi, err := os.Lstat(path)
	if err != nil {
		return 0, path
	}

	if !fi.IsDir() {
		if fi.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				return 0, path
			}
			return GetSize(target, showAll, recursive)
		}
		return fi.Size(), path
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return 0, path
	}

	var total int64
	for _, file := range files {
		if !showAll && strings.HasPrefix(file.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(path, file.Name())

		if !file.IsDir() {
			newfi, err := os.Lstat(fullPath)
			if err != nil {
				fmt.Printf("Ошибка при обработке %s: %v\n", fullPath, err)
				continue
			}

			if newfi.Mode()&os.ModeSymlink != 0 {
				target, err := os.Readlink(fullPath)
				if err != nil {
					continue
				}
				targetFi, err := os.Stat(target)
				if err != nil {
					continue
				}
				if !targetFi.IsDir() {
					total += targetFi.Size()
				}
			} else {
				total += newfi.Size()
			}
		} else {
			if recursive {
				subSize, _ := GetSize(fullPath, showAll, recursive)
				total += subSize
			}
		}
	}

	return total, path
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
