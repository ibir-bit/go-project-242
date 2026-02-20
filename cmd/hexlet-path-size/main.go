package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:    "hexlet-path-size",
		Usage:   "print size of a file or directory",
		Version: "v1.0.0",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "help",
				Aliases: []string{"h"},
				Usage:   "show help",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Основная логика программы
			fmt.Println("Программа для подсчета размера файлов и директорий")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}
