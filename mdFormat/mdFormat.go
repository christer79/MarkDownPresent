package mdFormat

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func init() {

}

// IsCommentedLine chercks if a string is a Markdown Comment or not
func IsCommentedLine(line string) bool {
	re := regexp.MustCompile("[[//]]")
	return re.MatchString(line)
}

// ExtractCommentDataFiled extracts labalse from a commented line given a key
func ExtractCommentDataFiled(line string, label string, def string) string {
	// TODO: Default return value woudl be clever here
	commentRe := regexp.MustCompile(label + ": \"([^\"]*)\"")
	value := commentRe.FindStringSubmatch(line)
	if value == nil {
		return def
	}
	return value[1]
}

// MarkDownToHTML takes markdown and returns html
func MarkDownToHTML(markdown []byte) []byte {
	content := blackfriday.MarkdownCommon(markdown)
	html := bluemonday.UGCPolicy().SanitizeBytes(content)
	prefix := []byte("<div>")
	suffix := []byte("</div>")
	html = append(prefix, html...)
	html = append(html, suffix...)

	return html
}
