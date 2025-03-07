package mdutils

import "regexp"

var jsonMarkdownRegex = regexp.MustCompile("```(json)?((\n|.)*?)```")

// ExtractJSONFromMarkdown returns the contents of the first fenced code block in
// the markdown text md. If there is none, it returns md.
func ExtractJSONFromMarkdown(md string) string {
	matches := jsonMarkdownRegex.FindStringSubmatch(md)
	if matches == nil {
		return md
	}
	return matches[2]
}
