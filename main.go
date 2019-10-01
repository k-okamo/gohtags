package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: gohtags <filename>\n")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Printf("<html>\n")
	fmt.Printf("<head>\n")
	fmt.Printf("<title>%s</title>\n", os.Args[1])
	fmt.Printf("</head>\n")
	fmt.Printf("<body>\n")
	fmt.Printf("<pre>\n")
	scanner := bufio.NewScanner(f)
	for i := 1; scanner.Scan(); i++ {
		str := scanner.Text()
		out = ""
		scan(str)
		fmt.Println(makeOneLine(htmlEncode(out), i))
	}
	fmt.Printf("</pre>\n")
	fmt.Printf("</body>\n")
	fmt.Printf("</html>\n")

}

func makeOneLine(s string, i int) string {
	return fmt.Sprintf("<a id='L%d' name='L%d'></a>%4d %s", i, i, i, s)
}

type replacePair struct {
	from string
	to   string
}

var htmlSymbol = []replacePair{
	{from: "&", to: "&amp"}, // & should be replaced first.
	{from: "<", to: "&lt"},
	{from: ">", to: "&gt"},
	{from: "\"", to: "&quot"},
	{from: "'", to: "&#39"},
	{from: "\t", to: "    "},
}

var htmlTag = []replacePair{
	{from: "TK_KEYWORD_S", to: "<strong class='reserved'>"},
	{from: "TK_KEYWORD_E", to: "</strong>"},
	{from: "TK_STRING_S", to: "<em class='string'>"},
	{from: "TK_STRING_E", to: "</em>"},
	{from: "TK_COMMENT_S", to: "<em class='comment'>"},
	{from: "TK_COMMENT_E", to: "</em>"},
}

func htmlEncode(s string) string {
	ss := s
	for _, p := range htmlSymbol {
		r := strings.NewReplacer(p.from, p.to)
		ss = r.Replace(ss)
	}
	for _, p := range htmlTag {
		r := strings.NewReplacer(p.from, p.to)
		ss = r.Replace(ss)
	}
	return ss
}
