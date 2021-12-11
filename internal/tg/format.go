package tg

import (
	"fmt"
	"strings"

	"github.com/xaxes/vikop-gorace/pkg/wykop"
)

var escChars = []byte{'_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!'}

const plusSign = "\xE2\x9E\x95"

func FormatEntry(entry wykop.Entry) string {
	var escPairs []string
	for _, ch := range escChars {
		str := string(ch)
		escPairs = append(escPairs, str, `\`+str)
	}
	repl := strings.NewReplacer(escPairs...)

	msg := fmt.Sprintf(`
%s %d \(%s\)
%s
*URL:* %s

%s
`,
		plusSign, entry.VoteCount, repl.Replace(entry.Author.Login),
		repl.Replace(entry.Embed.URL),
		repl.Replace(entry.URL), repl.Replace(entry.Body))

	if len(msg) > 4096 {
		return msg[:4096]
	}

	return msg

}
