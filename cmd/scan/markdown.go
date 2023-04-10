package scan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/pandodao/botastic-go"
	"gopkg.in/yaml.v2"
)

const (
	BlockTypeParagraph = iota
	BlockTypeList
	BlockTypeCode
	BlockTypeQuote
	BlockTypeUnknown
)

func extractMardownFileByParagraph(file string) ([]*botastic.CreateIndexesItem, error) {
	var sections []string

	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	frontmatter, content, err := extractFrontmatterAndContent(bytes)
	if err != nil {
		return nil, err
	}

	properties := "{}"
	if frontmatter != nil {
		properties, err = yamlToJSON(frontmatter)
		if err != nil {
			return nil, err
		}
	}

	lines := strings.Split(string(content), "\n")
	ix := 0
	section := ""
	for ix < len(lines) {
		line := strings.TrimSpace(lines[ix])
		if line == "" {
			// ignore empty line
			ix += 1
			continue
		}

		var currentBlockType int
		// check header
		headerRegex := regexp.MustCompile("^#{1,6} (.*)")
		headerMatch := headerRegex.FindStringSubmatch(line)
		if headerMatch != nil {
			// reach a title
			header := removeMarkdownSyntax(headerMatch[0])
			// save current section and start a new one
			if section != "" {
				sections = append(sections, section)
			}
			section = header
			ix += 1
			continue
		}

		// other lines
		currentBlockType = recognizeBlockType(line)
		if currentBlockType == BlockTypeUnknown {
			ix += 1
			continue
		}

		line = removeMarkdownSyntax(line)
		if line == "" {
			ix += 1
			continue
		}

		section = fmt.Sprintf("%s\n%s", section, line)
		ix += 1
	}

	if section != "" {
		sections = append(sections, section)
	}

	items := make([]*botastic.CreateIndexesItem, 0)
	for ix, sec := range sections {
		if isTrivialSentence(sec) {
			continue
		}
		items = append(items, &botastic.CreateIndexesItem{
			ObjectID:   fmt.Sprintf("%s/%s-%d", filepath.Dir(file), filepath.Base(file), ix),
			Category:   "plain-text",
			Data:       sec,
			Properties: properties,
		})
	}

	return items, nil
}

func extractMardownFileByLine(file string) ([]*botastic.CreateIndexesItem, error) {
	var sections []string

	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	frontmatter, content, err := extractFrontmatterAndContent(bytes)
	if err != nil {
		return nil, err
	}

	properties := "{}"
	if frontmatter != nil {
		properties, err = yamlToJSON(frontmatter)
		if err != nil {
			return nil, err
		}
	}

	lines := strings.Split(string(content), "\n")
	ix := 0
	header := ""
	for ix < len(lines) {
		line := strings.TrimSpace(lines[ix])
		if line == "" {
			// ignore empty line
			ix += 1
			continue
		}

		var currentBlockType int
		// check header
		headerRegex := regexp.MustCompile("^#{1,6} (.*)")
		headerMatch := headerRegex.FindStringSubmatch(line)
		if headerMatch != nil {
			// reach a title, move to a new section
			header = removeMarkdownSyntax(headerMatch[0])
			ix += 1
			if ix >= len(lines) {
				break
			}
			line = lines[ix]
		}

		// other lines
		currentBlockType = recognizeBlockType(line)
		if currentBlockType == BlockTypeUnknown {
			ix += 1
			continue
		}

		line = removeMarkdownSyntax(line)
		if line == "" {
			ix += 1
			continue
		}

		if header != "" {
			line = header + "\n" + line
			header = ""
		}

		sections = append(sections, line)
		ix += 1
	}

	items := make([]*botastic.CreateIndexesItem, 0)
	for ix, sec := range sections {
		if isTrivialSentence(sec) {
			continue
		}
		items = append(items, &botastic.CreateIndexesItem{
			ObjectID:   fmt.Sprintf("%s/%s-%d", filepath.Dir(file), filepath.Base(file), ix),
			Category:   "plain-text",
			Data:       sec,
			Properties: properties,
		})
	}

	return items, nil
}

func removeMarkdownSyntax(text string) string {
	text = strings.TrimSpace(text)
	// Regular expressions to match Markdown syntax
	boldRegex := regexp.MustCompile("\\*\\*(.*?)\\*\\*")
	italicRegex := regexp.MustCompile("_(.*?)_")
	strikethroughRegex := regexp.MustCompile("~~(.*?)~~")
	codeRegex := regexp.MustCompile("`(.*?)`")
	// linkRegex := regexp.MustCompile("\\[(.*?)\\]\\((.*?)\\)")
	// codeBlockRegex := regexp.MustCompile("```(.*?)```")
	quoteBlockRegex := regexp.MustCompile("> (.*?)")
	htmlRegex := regexp.MustCompile("<.*?>")

	// Replace bold syntax with plain text
	text = boldRegex.ReplaceAllString(text, "$1")

	// Replace italic syntax with plain text
	text = italicRegex.ReplaceAllString(text, "$1")

	// Replace strikethrough syntax with plain text
	text = strikethroughRegex.ReplaceAllString(text, "$1")

	// Replace code syntax with plain text
	text = codeRegex.ReplaceAllString(text, "$1")

	// Replace link syntax with plain text
	// text = linkRegex.ReplaceAllString(text, "$1")

	// Replace code block syntax with plain text
	// text = codeBlockRegex.ReplaceAllString(text, "$1")

	// Replace quote block syntax with plain text
	text = quoteBlockRegex.ReplaceAllString(text, "$1")

	// Remove HTML tags
	text = htmlRegex.ReplaceAllString(text, "")

	return strings.TrimSpace(text)
}

func recognizeBlockType(line string) int {
	if len(line) == 0 {
		return BlockTypeUnknown
	}
	if line[0] == '*' || line[0] == '-' || line[0] == '+' {
		return BlockTypeList
	}
	if line[0] == '>' {
		return BlockTypeQuote
	}
	if strings.HasPrefix(line, "```") {
		return BlockTypeCode
	}
	if line[0] != '[' && line[0] != '!' && line[0] != '|' {
		return BlockTypeParagraph
	}
	return BlockTypeUnknown
}

func extractFrontmatterAndContent(content []byte) ([]byte, []byte, error) {
	frontMatterPattern := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n`)
	match := frontMatterPattern.FindSubmatchIndex(content)

	if len(match) < 4 {
		return nil, content, nil
	}

	frontMatter := content[match[2]:match[3]]
	remainingContent := content[match[1]:]

	return frontMatter, bytes.TrimSpace(remainingContent), nil
}

func yamlToJSON(yamlContent []byte) (string, error) {
	var yamlObj map[string]interface{}
	err := yaml.Unmarshal(yamlContent, &yamlObj)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling YAML: %v", err)
	}

	jsonContent, err := json.Marshal(yamlObj)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	return string(jsonContent), nil
}

func isTrivialSentence(sentence string) bool {
	if len(sentence) < 16 {
		return true
	}

	// if sentence contains only symbols and numbers, it is trivial
	for _, r := range sentence {
		if !unicode.IsSymbol(r) && !unicode.IsNumber(r) && !unicode.IsSpace(r) && !unicode.IsPunct(r) {
			return false
		}
	}

	return true
}
