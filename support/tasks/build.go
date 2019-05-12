package main

import (
	"fmt"
	"github.com/lectio/graph/model"
	"github.com/lectio/graph/pipeline"
	"os"
	"path/filepath"
	"strings"
)

const defaultRepositoryName model.RepositoryName = "news.healthcareguys.com"

func configure() (*model.Configuration, error) {
	config, err := model.MakeConfiguration()
	if err != nil {
		return nil, err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var rootPath string
	switch {
	case strings.HasSuffix(currentDir, "tasks"):
		rootPath, err = filepath.Abs(filepath.Join(currentDir, "..", ".."))
	case strings.HasSuffix(currentDir, "support"):
		rootPath, err = filepath.Abs(filepath.Join(currentDir, ".."))
	default:
		rootPath, err = filepath.Abs(currentDir)
	}
	if err != nil {
		return nil, err
	}

	lls := config.LinkLifecyleSettings(model.DefaultSettingsPath)
	repos := config.Repositories(model.DefaultSettingsPath)

	repos.All = append(repos.All, model.FileRepository{
		Name:           defaultRepositoryName,
		URL:            model.URLText("file://" + rootPath),
		RootPath:       rootPath,
		CreateRootPath: false})

	lls.TraverseLinks = true
	lls.ScoreLinks.Score = false
	lls.ScoreLinks.Simulate = true

	return config, nil
}

func main() {
	config, err := configure()
	if err != nil {
		panic(err)
	}

	input := &model.BookmarksToMarkdownPipelineInput{
		BookmarksURL: "https://shah.dropmark.com/616548.json",
		Repository:   defaultRepositoryName,
		Settings:     model.DefaultSettingsPath,
		Strategy:     model.PipelineExecutionStrategyAsynchronous}

	task, err := pipeline.NewBookmarksToMarkdown(config, input)
	if err != nil {
		panic(err)
	}

	execution, err := task.Execute()
	if err != nil {
		panic(err)
	}
	result := execution.(*model.BookmarksToMarkdownPipelineExecution)

	fmt.Printf("[BookmarksToMarkdownPipelineExecution.Activities] %+v", result.Activities)
	fmt.Printf("[BookmarksToMarkdownPipelineExecution.Bookmarks.Activities] %+v", result.Bookmarks.Activities)
}
