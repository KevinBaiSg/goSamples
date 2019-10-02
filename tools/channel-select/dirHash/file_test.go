package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

func ExampleDirectory()  {
	root := "./"
	files := make([]string, 0)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}

	// Output:
	//
}

func ExampleTaskGenerator() {
	const max_thread = 10
	const root = "./"

	ctx, _ := context.WithCancel(context.Background()) // todo: cancel

	_, err := TaskGenerator(ctx, root)
	if err != nil {
		return
	}
	// Output:
	//
}