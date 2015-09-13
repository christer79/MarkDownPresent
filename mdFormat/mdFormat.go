package mdFormat

import "regexp"

func init() {

}

// IsCommentedLine chercks if a string is a Markdown Comment or not
func IsCommentedLine(line string) bool {
	re := regexp.MustCompile("[[//]]")
	return re.MatchString(line)
}

// ExtractCommentDataFiled extracts labalse from a commented line given a key
func ExtractCommentDataFiled(line string, label string) string {
	// TODO: Default return value woudl be clever here
	commentRe := regexp.MustCompile(label + ": \"([^\"]*)\"")
	value := commentRe.FindStringSubmatch(line)
	if value == nil {
		return ""
	}
	return value[1]
}
