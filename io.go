package main

import (
	"bufio"
	"fmt"
	"io"
	"ioutil"
	"log"
	"os"
)

var fileNo map[string]int

func MakeFileListPage(path string, w io.Writer) {

	w.Write([]byte("<html>\n"))
	w.Write([]byte("<head>\n"))
	buf := fmt.Sprintf("<title>%s</title>\n", path)
	w.Write([]byte(buf))
	printCSS(w)
	w.Write([]byte("</head>\n"))
	w.Write([]byte("<body>\n"))
	buf := fmt.Sprintf("<h2>%s</h2>\n", path)
	w.Write([]byte(buf))

	w.Write([]byte("<ul>"))
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for f := range files {
		buf := fmt.Sprintf("<li><a href=''>%s</a>\n", f.Name())
	}
	w.Write([]byte("</ul>"))

	w.Write([]byte("</body>\n"))
	w.Write([]byte("</html>\n"))

}

func WriteFile(f *os.File, w io.Writer) {

	w.Write([]byte("<html>\n"))
	w.Write([]byte("<head>\n"))
	buf := fmt.Sprintf("<title>%s</title>\n", f.Name())
	w.Write([]byte(buf))
	printCSS(w)
	w.Write([]byte("</head>\n"))
	w.Write([]byte("<body>\n"))
	w.Write([]byte("<pre>\n"))
	scanner := bufio.NewScanner(f)
	var s scanStr
	for i := 1; scanner.Scan(); i++ {
		s.in = scanner.Text()
		s.out = ""
		s.isSeparator = true
		scan(&s)
		buf = makeOneLine(s.out, i)
		w.Write([]byte(buf))
	}
	w.Write([]byte("</pre>\n"))
	w.Write([]byte("</body>\n"))
	w.Write([]byte("</html>\n"))
}

func printCSS(w io.Writer) {
	// embed CSS into html(temporary).
	w.Write([]byte("<style type='text/css'>"))
	w.Write([]byte("body{color: #B8BCC7; background-color: #222628; font-family: 'MyricaM M';font-size: 100%; line-height: 0.95em;}"))
	w.Write([]byte("pre {font-family: 'MyricaM M', Courier, sans-serif;}"))
	w.Write([]byte("a {color: #B9BCC7;}"))
	w.Write([]byte("em {font-style: normal;}"))
	w.Write([]byte("em.comment {color: #505C77;}"))
	w.Write([]byte("em.string {color: #89B3C2;}"))
	w.Write([]byte("strong.reserved {color: #708EBB;}"))
	w.Write([]byte("strong.constants {color: #CCA6BF;}"))
	w.Write([]byte("</style>"))
}

func makeOneLine(s string, i int) string {
	return fmt.Sprintf("<a id='L%d' name='L%d'></a>%4d %s\n", i, i, i, s)
}
