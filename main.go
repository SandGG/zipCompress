package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type file struct {
	name string
	body string
}

func main() {
	createFiles()
	var files = listFiles()
	addZipDir(files)
	for _, f := range files {
		fmt.Println(f)
	}
	fmt.Println("Zip Updated Successfully")
}

func createFiles() {
	var files = []file{
		{name: "text.txt", body: "Hello!\nInsert text here"},
		{name: "numbers.csv", body: "1,2,3\n4,5,6\n7,8,9"},
		{name: "names,txt", body: "Marcos\nAna\nRoman"},
	}

	for _, file := range files {
		var f, err = os.OpenFile("./files/"+file.name, os.O_WRONLY|os.O_CREATE, 0644)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		f.Write([]byte(file.body))
	}
}

func listFiles() []string {
	var files []string

	//type WalkFunc func(path string, info fs.FileInfo, err error) error
	filepath.Walk("./files", func(path string, info os.FileInfo, err error) error { //Walks the file tree rooted at root
		if !info.IsDir() { //It's a directory
			files = append(files, path) //Add files of directory
		}
		return nil
	})
	return files //return string whit files
}

func addZipDir(filepaths []string) {
	var file, err = os.OpenFile("./zip/files.zip", os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	//Create a new zip archive.
	var zipW = zip.NewWriter(file)
	defer zipW.Close()

	for _, filename := range filepaths {
		addFileToZip(filename, zipW)
	}
}

func addFileToZip(filename string, zipW *zip.Writer) {
	var file, err = os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	//Create file in zip
	var wr, errCreate = zipW.Create(filename)
	if err != nil {
		log.Fatal(errCreate)
	}

	//Copy info of file in file zip
	var _, errCopy = io.Copy(wr, file)
	if err != nil {
		log.Fatal(errCopy)
	}
}
