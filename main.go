package main

import (
	"fmt"
	"os"

	"hfdl/cmd"

	"github.com/spf13/pflag"
)

func main() {
	var (
		modelID     string
		outputDir   string
		fileToFetch string
		token       string
	)

	pflag.StringVarP(&modelID, "model", "m", "", "Hugging Face model ID")
	pflag.StringVarP(&outputDir, "output", "o", "./downloads", "Base output directory")
	pflag.StringVarP(&fileToFetch, "file", "f", "", "Download a specific file")
	pflag.StringVarP(&token, "token", "t", "", "Hugging Face access token")
	pflag.Parse()

	if modelID == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s --model <model_id> [--output <dir>] [--file <filename>] [--token <token>]\n", os.Args[0])
		os.Exit(1)
	}

	config := cmd.NewDownloadConfig(modelID, outputDir, fileToFetch, token)
	if err := config.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
