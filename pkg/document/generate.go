package document

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/pelletier/go-toml/v2/unstable"
)

// GenerateOptions represents the options for generating the markdown file.
type GenerateOptions struct {
	// DebugMode enables to print debug information.
	DebugMode bool

	// DocsCommentPrefix is the prefix of the comments that will be used to generate the documentation.
	DocsCommentPrefix string
}

const (
	PlaceholderForEmpty = "--"
	BRTag               = "<br/>"
	NoneValue           = "None"
)

const (
	NoneDefault = "+toml2docs:none-default"
)

// globalOpts is the global options for generating the markdown file.
var globalOpts = &GenerateOptions{DebugMode: false, DocsCommentPrefix: "#"}

// GenerateMarkdown generates a markdown file from the input toml data.
func GenerateMarkdown(input []byte, opts *GenerateOptions) (string, error) {
	nodes, err := parse(input)
	if err != nil {
		return "", err
	}

	if opts != nil {
		globalOpts = opts
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

// GenerateMarkdownFromFile generates a markdown file from the input toml file.
func GenerateMarkdownFromFile(filename string, opts *GenerateOptions) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return GenerateMarkdown(data, opts)
}

// GenerateMarkdownFromTemplate generates a markdown file from a Go template with `toml2docs` function.
func GenerateMarkdownFromTemplate(templateFileName string, opts *GenerateOptions) (string, error) {
	data, err := os.ReadFile(templateFileName)
	if err != nil {
		return "", err
	}

	if opts != nil {
		globalOpts = opts
	}

	funcMap := template.FuncMap{
		"toml2docs": toml2docs,
	}

	tmpl, err := template.New("toml2docs").Funcs(funcMap).Parse(string(data))
	if err != nil {
		return "", nil
	}

	w := &strings.Builder{}
	if err := tmpl.Execute(w, nil); err != nil {
		return "", err
	}

	return reduceRedundantNewline(w.String()), nil
}

func toml2docs(filename string) (string, error) {
	return GenerateMarkdownFromFile(filename, nil)
}

type tomlNode struct {
	Kind unstable.Kind
	Data []byte
}

type docItem struct {
	key     string
	val     string
	typ     string
	comment *tomlComment
}

type tomlComment struct {
	rawComments []string
	noneDefault bool
}

func (c *tomlComment) String() string {
	return normalize(strings.Join(c.rawComments, BRTag), false)
}

func (c *tomlComment) IsNoneDefault() bool {
	return c.noneDefault
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
		comment   = new(tomlComment)
		items     []*docItem

		// arrayTableIndex is used to keep track of the array table index.
		arrayTableIndex = make(map[string]int)
	)

	// Process the nodes to generate the doc items.
	for cursor < len(nodes) {
		node := nodes[cursor]
		switch node.Kind {
		case unstable.Comment:
			if !strings.HasPrefix(string(node.Data), globalOpts.DocsCommentPrefix) {
				// Skip the comment.
				cursor++
				continue
			}
			rawComment := processComment(string(node.Data))
			if strings.Contains(rawComment, NoneDefault) {
				comment.noneDefault = true
			} else {
				comment.rawComments = append(comment.rawComments, rawComment)
			}

			// Move to the next node.
			cursor++
		case unstable.KeyValue:
			n := peek(nodes, cursor+1)
			if n == nil {
				return nil, fmt.Errorf("missing key value")
			}

			if n.Kind == unstable.Array {
				// Skip the array node. The actual value is the next node.
				cursor++
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
				comment: comment,
			})

			// Take the tomlComment and reset it.
			comment = new(tomlComment)
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
				comment: comment,
			})

			// Take the tomlComment and reset it.
			comment = new(tomlComment)
		case unstable.ArrayTable:
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

			index := 0
			if _, ok := arrayTableIndex[parentKey]; ok {
				index = arrayTableIndex[parentKey] + 1
				arrayTableIndex[parentKey] = index
			} else {
				arrayTableIndex[parentKey] = 0
			}

			if index == 0 {
				items = append(items, &docItem{
					key:     normalize(fmt.Sprintf("[[%s]]", parentKey), true),
					val:     normalize("", true),
					typ:     normalize("", false),
					comment: comment,
				})
				comment = new(tomlComment)
			}
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
		if item.comment.IsNoneDefault() {
			item.val = normalize(NoneValue, true)
		}
		if item.typ == "String" && item.val == PlaceholderForEmpty {
			item.val = "`\"\"`" // indicate empty string.
		}

		buf.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", item.key, item.typ, item.val, item.comment.String()))
	}

	return reduceRedundantNewline(buf.String()), nil
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
	return strings.TrimPrefix(strings.TrimLeft(input2, "#"), " ")
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

func reduceRedundantNewline(s string) string {
	re := regexp.MustCompile(`\n+$`)
	return re.ReplaceAllString(s, "\n")
}
