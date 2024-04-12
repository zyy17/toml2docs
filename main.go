package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2/unstable"
	"github.com/spf13/pflag"
)

type tomlNode struct {
	Kind unstable.Kind
	Data []byte
}

type docItem struct {
	key     string
	val     string
	comment string
}

func main() {
	var (
		inputFile  string
		outputFile string
		debug      bool
	)

	pflag.StringVarP(&inputFile, "input-file", "i", "", "The input toml file path.")
	pflag.StringVarP(&outputFile, "output-file", "o", "output.md", "The output markdown file path.")
	pflag.BoolVarP(&debug, "debug", "d", false, "Print debug information.")
	pflag.Parse()

	if len(inputFile) == 0 {
		pflag.PrintDefaults()
		os.Exit(1)
	}

	nodes, err := parse(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if debug {
		for _, node := range nodes {
			fmt.Printf("%s: %s\n", node.Kind, node.Data)
		}
		os.Exit(0)
	}

	items, err := generateDocItems(nodes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := generateMarkdown(items, outputFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parse(inputFile string) ([]*tomlNode, error) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, err
	}

	// Gather all the nodes in the TOML file.
	var nodes []*tomlNode

	var parseNode func(*unstable.Parser, int, *unstable.Node)
	parseNode = func(p *unstable.Parser, indent int, e *unstable.Node) {
		if e == nil {
			return
		}
		nodes = append(nodes, &tomlNode{Kind: e.Kind, Data: e.Data})
		parseNode(p, indent+1, e.Child())
		parseNode(p, indent, e.Next())
	}

	// Should keep the comments to be able to generate the markdown file.
	parser := &unstable.Parser{KeepComments: true}
	parser.Reset(data)

	for parser.NextExpression() {
		parseNode(parser, 0, parser.Expression())
	}

	if err := parser.Error(); err != nil {
		return nil, err
	}

	return nodes, nil
}

func generateDocItems(nodes []*tomlNode) ([]*docItem, error) {
	var (
		cursor    = 0
		parentKey = ""
		comment   = ""
		items     []*docItem
	)

	// Process the nodes to generate the doc items.
	for cursor < len(nodes) {
		node := nodes[cursor]
		switch node.Kind {
		case unstable.Comment:
			comment = processComment(string(node.Data))

			// If the previous node is a comment, append the current comment to the previous one.
			p := peek(nodes, cursor-1)
			if p != nil && p.Kind == unstable.Comment {
				comment = processComment(string(p.Data) + " " + comment)
			}

			// Move to the next node.
			cursor = cursor + 1
		case unstable.KeyValue:
			n := peek(nodes, cursor+1)
			if n == nil {
				return nil, fmt.Errorf("missing key value")
			}

			if n.Kind == unstable.Array {
				// Skip the array node. The actual value is the next node.
				cursor = cursor + 1
			}

			nn := peek(nodes, cursor+2)
			if nn == nil {
				return nil, fmt.Errorf("missing key")
			}

			key := string(nn.Data)
			if len(parentKey) > 0 {
				key = parentKey + "." + key
			}

			items = append(items, &docItem{
				comment: comment,
				val:     string(n.Data),
				key:     key,
			})

			// Take the comment and reset it.
			comment = ""
			cursor += 3
		case unstable.Table:
			n := peek(nodes, cursor+1)
			if n == nil {
				return nil, fmt.Errorf("missing key for table")
			}
			parentKey = string(n.Data)

			nn := peek(nodes, cursor+2)
			// If the next node is a key, append it to the parent key.
			if nn != nil && nn.Kind == unstable.Key {
				parentKey += "." + string(nn.Data)
				cursor = cursor + 1
			}
			items = append(items, &docItem{
				comment: comment,
				val:     "",
				key:     parentKey,
			})

			// Take the comment and reset it.
			comment = ""
			cursor += 2
		case unstable.Array:
			cursor += 2
		default:
			// Stop the loop.
			cursor = len(nodes)
		}
	}

	return items, nil
}

func generateMarkdown(items []*docItem, outputFile string) error {
	buf := strings.Builder{}
	buf.WriteString("| Key | Default | Descriptions |\n")
	buf.WriteString("| --- | ------- | ----------- |\n")

	for _, item := range items {
		var (
			key     = item.key
			val     = item.val
			comment = item.comment
		)

		if len(val) == 0 {
			val = "--"
		} else {
			val = "`" + val + "`"
		}

		if len(comment) == 0 {
			comment = "--"
		}

		if len(key) == 0 {
			key = "--"
		} else {
			key = "`" + key + "`"
		}

		buf.WriteString(fmt.Sprintf("| %s | %s | %s |\n", key, val, comment))
	}

	return os.WriteFile(outputFile, []byte(buf.String()), 0644)
}

// peek returns the node at the given index.
func peek(nodes []*tomlNode, i int) *tomlNode {
	if i < len(nodes) && i >= 0 {
		return nodes[i]
	}
	return nil
}

// processComment removes the comment prefix and trims the spaces.
func processComment(input string) string {
	input2 := strings.TrimPrefix(input, " ")
	return strings.TrimPrefix(strings.TrimPrefix(input2, "#"), " ")
}
