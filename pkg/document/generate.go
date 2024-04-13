package document

import (
	"fmt"
	"strings"

	"github.com/pelletier/go-toml/v2/unstable"
)

// GenerateOptions represents the options for generating the markdown file.
type GenerateOptions struct {
	// DebugMode enables to print debug information.
	DebugMode bool
}

const (
	PlaceholderForEmpty = "--"
)

// GenerateMarkdown generates a markdown file from the input toml data.
func GenerateMarkdown(input []byte, opts *GenerateOptions) (string, error) {
	nodes, err := parse(input)
	if err != nil {
		return "", err
	}

	if opts != nil && opts.DebugMode {
		for _, node := range nodes {
			fmt.Printf("%s: %s\n", node.Kind, node.Data)
		}
		return "", nil
	}

	items, err := generateDocItems(nodes)
	if err != nil {
		return "", err
	}

	return doGenerateMarkdown(items)
}

type tomlNode struct {
	Kind unstable.Kind
	Data []byte
}

type docItem struct {
	key     string
	val     string
	typ     string
	comment string
}

func parse(data []byte) ([]*tomlNode, error) {
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

		// arrayTableIndex is used to keep track of the array table index.
		arrayTableIndex = make(map[string]int)
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
				key:     normalize(key, true),
				val:     normalize(string(n.Data), true),
				typ:     normalize(n.Kind.String(), false),
				comment: normalize(comment, false),
			})

			// Take the comment and reset it.
			comment = ""
			cursor += 3
		case unstable.Table:
			// Reset the parent key.
			parentKey = ""
			for i := cursor + 1; i < len(nodes); i++ {
				n := peek(nodes, i)
				if n == nil {
					cursor = i
					break
				}
				if n.Kind != unstable.Key {
					cursor = i
					break
				}

				// If the next node is a key, append it to the parent key.
				if parentKey != "" {
					parentKey = parentKey + "." + string(n.Data)
				} else {
					parentKey = string(n.Data)
				}
			}

			n := peek(nodes, cursor+1)
			if n == nil {
				return nil, fmt.Errorf("missing key for table")
			}
			items = append(items, &docItem{
				key:     normalize(parentKey, true),
				val:     normalize("", true),
				typ:     normalize("", false),
				comment: normalize(comment, false),
			})

			// Take the comment and reset it.
			comment = ""
		case unstable.ArrayTable:
			parentKey = ""
			n := peek(nodes, cursor+1)
			if n == nil {
				return nil, fmt.Errorf("missing array table key")
			}

			var arrayKey string
			if n.Kind == unstable.Key {
				arrayKey = string(n.Data)
			}

			if len(parentKey) > 0 {
				arrayKey = parentKey + "." + string(n.Data)
			}

			index := 0
			if _, ok := arrayTableIndex[arrayKey]; ok {
				index = arrayTableIndex[arrayKey] + 1
				arrayTableIndex[arrayKey] = index
			} else {
				arrayTableIndex[arrayKey] = 0
			}
			parentKey = fmt.Sprintf("%s[%d]", arrayKey, index)
			cursor += 2
		case unstable.Array:
			cursor += 2
		default:
			return nil, fmt.Errorf("unexpected node kind: %s", node.Kind)
		}
	}

	return items, nil
}

func doGenerateMarkdown(items []*docItem) (string, error) {
	buf := strings.Builder{}
	buf.WriteString("| Key | Type | Default | Descriptions |\n")
	buf.WriteString("| --- | -----| ------- | ----------- |\n")

	for _, item := range items {
		buf.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", item.key, item.typ, item.val, item.comment))
	}

	return buf.String(), nil
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

func normalize(input string, isCode bool) string {
	if len(input) == 0 {
		return PlaceholderForEmpty
	}

	if isCode {
		return "`" + input + "`"
	}

	return input
}
