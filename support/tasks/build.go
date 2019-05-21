package main

import (
	"context"
	"fmt"
	"github.com/lectio/link"
	"github.com/lectio/markdown"
	"github.com/lectio/publish"
	"path/filepath"
	"runtime"
	"strings"
)

func basePath(executablePath string) (string, error) {
	var basePath string
	var err error
	switch {
	case strings.HasSuffix(executablePath, "tasks"):
		basePath, err = filepath.Abs(filepath.Join(executablePath, "..", ".."))
	case strings.HasSuffix(executablePath, "support"):
		basePath, err = filepath.Abs(filepath.Join(executablePath, ".."))
	default:
		basePath, err = filepath.Abs(executablePath)
	}
	if err != nil {
		return "", err
	}
	return basePath, nil
}

func main() {
	// ex, err := os.Executable()
	// if err != nil {
	// 	panic(err)
	// }
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Unable to get name of executable go file")
		return
	}
	executablePath := filepath.Dir(filename)

	basePath, err := basePath(executablePath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Configured basePath: %s.\n", basePath)

	ctx := context.Background()

	bpc := markdown.NewBasePathConfigurator(basePath)
	linkFactory := link.NewFactory()

	publisher, err := publish.NewMarkdownPublisher(ctx, true, linkFactory, bpc, &publish.CommandLineProgressReporter{})
	if err != nil {
		panic(err)
	}

	if err = publisher.Publish(ctx, "https://shah.dropmark.com/616548.json"); err != nil {
		panic(err)
	}
}
