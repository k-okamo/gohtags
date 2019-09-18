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
		fmt.Println(makeOneLine(htmlEncode(scanner.Text()), i))
	}
	fmt.Printf("</pre>\n")
	fmt.Printf("</body>\n")
	fmt.Printf("</html>\n")
}

func makeOneLine(s string, i int) string {
	return fmt.Sprintf("<a id='L%d' name='L%d'></a>%4d %s", i, i, i, s)
}

func htmlEncode(s string) string {
	var ss string
	// & should be replaced first.
	ss = strings.Replace(s, "&", "&amp;", -1)
	ss = strings.Replace(ss, "<", "&lt;", -1)
	ss = strings.Replace(ss, ">", "&gt;", -1)
	ss = strings.Replace(ss, "\"", "&quot;", -1)
	ss = strings.Replace(ss, "'", "&#39;", -1)
	ss = strings.Replace(ss, "\t", "    ", -1)
	return ss
}
