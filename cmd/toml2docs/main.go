package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/zyy17/toml2docs/pkg/document"
)

func main() {
	var (
		inputFile  string
		outputFile string
		debug      bool
	)

	pflag.StringVarP(&inputFile, "input-file", "i", "", "The input toml file path.")
	pflag.StringVarP(&outputFile, "output-file", "o", "", "The output markdown file path. If not provided, it will be printed to the standard output.")
	pflag.BoolVarP(&debug, "debug", "d", false, "Print debug information.")
	pflag.Parse()

	if len(inputFile) == 0 {
		pflag.PrintDefaults()
		os.Exit(1)
	}

	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	docs, err := document.GenerateMarkdown(data, &document.GenerateOptions{DebugMode: debug})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(outputFile) == 0 {
		// Print the markdown file to the standard output.
		fmt.Println(docs)
	} else {
		if err := os.WriteFile(outputFile, []byte(docs), 0644); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
