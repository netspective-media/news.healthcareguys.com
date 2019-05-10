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
	result, err := model.MakeConfiguration()
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

	settings := result.DefaultBundle()
	settings.Repositories.All = append(settings.Repositories.All, model.FileRepository{
		Name:           defaultRepositoryName,
		URL:            model.URLText("file://" + rootPath),
		RootPath:       rootPath,
		CreateRootPath: false})

	settings.Links.TraverseLinks = true
	settings.Links.ScoreLinks.Score = true
	settings.Links.ScoreLinks.Simulate = true

	return result, nil
}

func main() {
	config, err := configure()
	if err != nil {
		panic(err)
	}

	input := &model.BookmarksToMarkdownPipelineInput{
		BookmarksURL:       "https://shah.dropmark.com/616548.json",
		Flavor:             model.MarkdownFlavorHugoContent,
		Repository:         defaultRepositoryName,
		SettingsBundle:     model.DefaultSettingsBundleName,
		Strategy:           model.PipelineExecutionStrategyAsynchronous,
		ContentPathRel:     "content/post",
		ImagesCachePathRel: "static/img/content/post",
		ImagesCacheRootURL: "/img/content/post"}

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
