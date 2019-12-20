package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	//	if len(os.Args) < 2 {
	//		fmt.Fprintf(os.Stderr, "usage: gohtags <filename>\n")
	//		os.Exit(1)
	//	}
	//
	/*
		files, err := ioutil.ReadDir(".")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, file := range files {
			name := file.Name()
			if strings.HasSuffix(name, ".go") {
				f, err := os.Open(name)
				if err != nil {
					fmt.Println(err)
					continue
				}
				defer f.Close()
				w := os.Stdout
				WriteFile(f, w)
			}
		}
	*/
	files := readFiles(".")
	err := os.Mkdir("HTML", 0777)
	if err != nil {
		//log.Fatal(err)
	}
	i := 1
	for _, file := range files {
		rf, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//defer rf.Close()
		wf, err := os.Create(filepath.Join("HTML", strconv.Itoa(i)+".html"))
		if err != nil {
			fmt.Println(err)
			continue
		}
		//defer wf.Close()
		fmt.Printf("Processing %s\n", file)
		WriteFile(rf, wf)
		wf.Close()
		rf.Close()
		i++
	}
}

//readFiles(".")
func readFiles(search string) []string {

	var ret []string
	files, err := ioutil.ReadDir(search)

	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range files {
		subpath := filepath.Join(search, fi.Name())
		if fi.IsDir() {
			ret = append(ret, readFiles(subpath)...)
		} else {
			if !strings.HasSuffix(subpath, ".go") && !strings.HasSuffix(subpath, ".txt") {
				continue
			}
			ret = append(ret, subpath)
		}
	}
	return ret
}
