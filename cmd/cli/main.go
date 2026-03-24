package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/bigwhite/issue2md/internal/converter"
	"github.com/bigwhite/issue2md/internal/fetcher"
	"github.com/bigwhite/issue2md/internal/parser"
	"github.com/bigwhite/issue2md/internal/writer"
)

// Config holds CLI configuration
type Config struct {
	Token     string
	UserLink  bool
	OutputFile string
	URL       string
}

// getTokenFromEnv reads GITHUB_TOKEN from environment
func getTokenFromEnv() (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return "", errors.New("GITHUB_TOKEN environment variable is not set")
	}
	return token, nil
}

// parseFlags parses command-line flags
func parseFlags(args []string) (*Config, error) {
	if len(args) == 0 {
		return nil, errors.New("requires a URL argument")
	}

	cfg := &Config{}

	fs := flag.NewFlagSet("issue2md", flag.ContinueOnError)
	userLink := fs.Bool("user-link", false, "render usernames as links")
	fs.BoolVar(userLink, "u", false, "short for --user-link")

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	remaining := fs.Args()
	if len(remaining) == 0 {
		return nil, errors.New("requires a URL argument")
	}

	cfg.URL = remaining[0]
	cfg.UserLink = *userLink

	if len(remaining) > 1 {
		cfg.OutputFile = remaining[1]
	}

	return cfg, nil
}

// run is the main execution logic
// baseURL 用于指定 GitHub API 地址，测试时可以使用 mock server
func run(ctx context.Context, cfg *Config, baseURL string) error {
	// 1. Parse URL
	parsedURL, err := parser.Parse(cfg.URL)
	if err != nil {
		return err
	}

	// 2. Create Fetcher with GitHub client
	if cfg.Token == "" {
		return errors.New("GITHUB_TOKEN environment variable is not set")
	}
	client := fetcher.NewGitHubClient(baseURL, cfg.Token)
	fet := fetcher.NewGitHubFetcher(client)

	// 3. Fetch resource from GitHub
	resource, err := fet.Fetch(ctx, parsedURL)
	if err != nil {
		return err
	}

	// 4. Convert to Markdown
	conv := converter.NewConverter()
	output, err := conv.Convert(ctx, resource, cfg.UserLink)
	if err != nil {
		return err
	}

	// 5. Write output
	w := writer.New()
	return w.Write(output.FullContent, cfg.OutputFile)
}

func main() {
	// 读取 GITHUB_TOKEN
	token, err := getTokenFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// 解析命令行 flags
	cfg, err := parseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	cfg.Token = token

	// 执行主流程
	err = run(context.Background(), cfg, "https://api.github.com")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
