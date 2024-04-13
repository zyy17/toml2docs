package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/zyy17/toml2docs/pkg/document"
)

func main() {
	var (
		inputFile    string
		outputFile   string
		templateFile string
		debug        bool
	)

	pflag.StringVarP(&inputFile, "input-file", "i", "", "The input toml file path.")
	pflag.StringVarP(&outputFile, "output-file", "o", "", "The output markdown file path. If not provided, it will be printed to the standard output.")
	pflag.StringVarP(&templateFile, "template-file", "t", "", "The template file path.")
	pflag.BoolVarP(&debug, "debug", "d", false, "Print debug information.")
	pflag.Parse()

	if len(inputFile) == 0 && len(templateFile) == 0 {
		pflag.Usage()
		os.Exit(1)
	}

	if len(inputFile) > 0 && len(templateFile) > 0 {
		fmt.Println("Please only provide one of the input file or the template file.")
		os.Exit(1)
	}

	var (
		docs string
		err  error
	)

	if len(inputFile) > 0 {
		docs, err = document.GenerateMarkdownFromFile(inputFile, &document.GenerateOptions{DebugMode: debug})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if len(templateFile) > 0 {
		docs, err = document.GenerateMarkdownFromTemplate(templateFile, &document.GenerateOptions{DebugMode: debug})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if len(outputFile) == 0 {
		// Print the markdown file to the standard output.
		fmt.Printf("%s", docs)
	} else {
		if err := os.WriteFile(outputFile, []byte(docs), 0644); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
