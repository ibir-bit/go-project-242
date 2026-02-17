package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetSize_File(t *testing.T) {
	tmpDir := t.TempDir()
	
	content := []byte("Hello, World!")
	tmpFile := filepath.Join(tmpDir, "test_file.txt")
	err := os.WriteFile(tmpFile, content, 0644)
	if err != nil {
		t.Fatalf("Не удалось создать тестовый файл: %v", err)
	}

	size, err := GetSize(tmpFile)
	if err != nil {
		t.Errorf("GetSize вернул ошибку для файла: %v", err)
	}

	expected := int64(len(content))
	if size != expected {
		t.Errorf("GetSize = %d, ожидалось %d", size, expected)
	}
}

func TestGetSize_Directory_FirstLevel(t *testing.T) {
	tmpDir := t.TempDir()

	// Создаем файлы первого уровня
	file1 := filepath.Join(tmpDir, "file1.txt")
	os.WriteFile(file1, []byte("12345"), 0644)

	file2 := filepath.Join(tmpDir, "file2.txt")
	os.WriteFile(file2, []byte("12345"), 0644)

	// Создаем поддиректорию с файлом (не должен учитываться)
	subDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDir, 0755)
	file3 := filepath.Join(subDir, "file3.txt")
	os.WriteFile(file3, []byte("12345"), 0644)

	// Ожидаемый размер: только file1.txt + file2.txt = 10 байт
	expected := int64(10)

	size, err := GetSize(tmpDir)
	if err != nil {
		t.Errorf("GetSize вернул ошибку для директории: %v", err)
	}

	if size != expected {
		t.Errorf("GetSize = %d, ожидалось %d", size, expected)
	}
}

func TestGetSize_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	size, err := GetSize(tmpDir)
	if err != nil {
		t.Errorf("GetSize вернул ошибку для пустой директории: %v", err)
	}

	if size != 0 {
		t.Errorf("GetSize = %d, ожидалось 0 для пустой директории", size)
	}
}

func TestGetSize_NonExistentPath(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistent := filepath.Join(tmpDir, "does_not_exist.txt")

	size, err := GetSize(nonExistent)
	if err == nil {
		t.Error("GetSize должен вернуть ошибку для несуществующего пути")
	}

	if size != 0 {
		t.Errorf("GetSize = %d, ожидалось 0 при ошибке", size)
	}
}

func TestGetSize_WithSymlink(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Создаем файл
	targetFile := filepath.Join(tmpDir, "target.txt")
	err := os.WriteFile(targetFile, []byte("Target"), 0644)
	if err != nil {
		t.Fatalf("Не удалось создать целевой файл: %v", err)
	}
	
	// Создаем симлинк
	symlink := filepath.Join(tmpDir, "symlink.txt")
	err = os.Symlink(targetFile, symlink)
	if err != nil {
		t.Skip("Пропускаем тест симлинков: не поддерживается на этой платформе")
	}
	
	// Размер симлинка должен быть размером файла, на который он указывает
	size, err := GetSize(symlink)
	if err != nil {
		t.Errorf("GetSize вернул ошибку для симлинка: %v", err)
	}
	
	expected := int64(6) // "Target" = 6 байт
	if size != expected {
		t.Errorf("GetSize = %d, ожидалось %d", size, expected)
	}
}

func TestGetSize_TableDriven(t *testing.T) {
	tmpDir := t.TempDir()

	// Создаём тестовые файлы
	os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("123"), 0644)   // 3 байта
	os.WriteFile(filepath.Join(tmpDir, "b.txt"), []byte("12345"), 0644) // 5 байт
	os.Mkdir(filepath.Join(tmpDir, "sub"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "sub", "c.txt"), []byte("1234567"), 0644) // 7 байт (не учитывается)

	tests := []struct {
		name     string
		path     string
		expected int64
		wantErr  bool
	}{
		{
			name:     "файл a.txt",
			path:     filepath.Join(tmpDir, "a.txt"),
			expected: 3,
			wantErr:  false,
		},
		{
			name:     "файл b.txt",
			path:     filepath.Join(tmpDir, "b.txt"),
			expected: 5,
			wantErr:  false,
		},
		{
			name:     "директория (только файлы первого уровня)",
			path:     tmpDir,
			expected: 8, // 3 + 5 = 8 (c.txt в sub не учитывается)
			wantErr:  false,
		},
		{
			name:     "несуществующий файл",
			path:     filepath.Join(tmpDir, "none.txt"),
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size, err := GetSize(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if size != tt.expected {
				t.Errorf("GetSize() = %d, expected %d", size, tt.expected)
			}
		})
	}
}

func BenchmarkGetSize(b *testing.B) {
	tmpDir := b.TempDir()

	for i := 0; i < 10; i++ {
		file := filepath.Join(tmpDir, fmt.Sprintf("file%d.txt", i))
		os.WriteFile(file, []byte("Benchmark test"), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetSize(tmpDir)
	}
}
