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

func htmlEncode(s string) string {
	// & should be replaced first.
	rSymbol := strings.NewReplacer("&", "&amp", "<", "&lt", ">", "&gt", "\"", "&quot", "'", "&#39", "\t", "    ")
	rDquote := strings.NewReplacer("TK_STRING_S", "<em class='string'>", "TK_STRING_E", "</em>")
	rComment := strings.NewReplacer("TK_COMMENT_S", "<em class='comment'>", "TK_COMMENT_E", "</em>")

	return rComment.Replace(rDquote.Replace(rSymbol.Replace(s)))
}
