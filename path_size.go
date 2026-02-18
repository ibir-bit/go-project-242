package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {

}

func GetSize(adress string) (int64, error) {
	var human bool
	flag.BoolVar(&human, "h", false, "human-readable sizes (auto-select unit)")

	flag.Parse()

	_, err := os.Lstat(adress)
	if err != nil {
		return 0, errors.New("Такого файла не существует")
	}

	targetInfo, err := os.Stat(adress)
	if err != nil {
		return 0, fmt.Errorf("не могу получить информацию: %w", err)
	}

	if !targetInfo.IsDir() {
		return targetInfo.Size(), nil
	}

	entries, err := os.ReadDir(adress)
	if err != nil {
		return 0, errors.New("Такого пути не существует")
	}

	var size int64

	for _, entry := range entries {
		fullPath := filepath.Join(adress, entry.Name())

		if entry.IsDir() {
			continue
		}

		info, err := os.Stat(fullPath)
		if err != nil {
			return 0, errors.New("Ошибка чтения")
		}

		size += info.Size()
	}

	return size, nil
}

func FormatSize(human bool, size int64) string {
	if human {
		fmt.Println(HumanizeSize(size))
	} else {
		fmt.Println(size)
	}
}

func HumanizeSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2fTB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2fGB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2fMB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2fKB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}
