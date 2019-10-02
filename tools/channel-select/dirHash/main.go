package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	const max_thread = 10
	const root = "/Users/kevin/working/myleetcode"
	var wg sync.WaitGroup

	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	fileStream, err := TaskGenerator(ctx, root)
	if err != nil {
		log.Fatalln(err)
		return
	}

	wg.Add(max_thread)
	for i := 0; i < max_thread; i++ {
		go func() error {
			defer wg.Done()
			for {
				select {
				case <- ctx.Done():
					return ctx.Err()
				case file, ok := <- fileStream:
					if ok == false {
						return errors.New("fileStream is closed")
					}
					handleHash(file)
				}
			}
		}()
	}
	wg.Wait()
}

func handleHash(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(file, err)
		return
	}
	hash := sha256.Sum256(data)
	log.Printf("%s: %x", file, hash)
}

func TaskDFS(ctx context.Context, root string, fileChan chan string) error {
	// dfs
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			TaskDFS(ctx, filepath.Join(root, file.Name()), fileChan)
		} else {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case fileChan <- filepath.Join(root, file.Name()):
			}
		}
	}

	return nil
}

func TaskGenerator(ctx context.Context, root string) (chan string, error) {
	filePathStream := make(chan string, 10)

	go func() {
		defer close(filePathStream)
		TaskDFS(ctx, root, filePathStream)
	}()

	return filePathStream, nil
}
