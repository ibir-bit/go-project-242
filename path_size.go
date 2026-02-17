package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func main() {

}

func GetSize(adress string) (int64, error) {
	// Проверяем существование пути
	_, err := os.Lstat(adress)
	if err != nil {
		return 0, errors.New("Такого файла не существует")
	}

	// Получаем информацию о целевом файле (переходя по ссылкам)
	targetInfo, err := os.Stat(adress)
	if err != nil {
		return 0, fmt.Errorf("не могу получить информацию: %w", err)
	}

	// Если это файл - возвращаем размер
	if !targetInfo.IsDir() {
		return targetInfo.Size(), nil
	}

	// Читаем содержимое директории (os.ReadDir автоматически работает с симлинками)
	entries, err := os.ReadDir(adress)
	if err != nil {
		return 0, errors.New("Такого пути не существует")
	}

	var size int64

	// Проходим по всем файлам в директории
	for _, entry := range entries {
		fullPath := filepath.Join(adress, entry.Name())

		// Пропускаем поддиректории
		if entry.IsDir() {
			continue
		}

		// Получаем информацию о файле (переходя по ссылкам)
		info, err := os.Stat(fullPath)
		if err != nil {
			return 0, errors.New("Ошибка чтения")
		}

		size += info.Size()
	}

	return size, nil
}
