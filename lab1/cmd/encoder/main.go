package main

import (
	"fmt"
	"lab1/pkg/crypt"
	"log"
	"os"
	"path"
)

func must[T any](val T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func readFile(filename string) string {
	bytes := must(os.ReadFile(filename))
	return string(bytes)
}

func main() {
	keyWord := "bar"
	base := "./texts"
	files := must(os.ReadDir(base))
	for _, fileInfo := range files {
		filename := path.Join(base, fileInfo.Name())
		text := readFile(filename)
		matrix := crypt.Matrix(len(keyWord)+1, keyWord)
		keyWordLine := crypt.KeyWordLine(text, keyWord)
		encodedText := crypt.EncodeLine(text, keyWordLine, matrix, keyWord)
		fmt.Println(encodedText)
	}
}
