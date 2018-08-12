package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

// Templates is the map of pointers to the templates as per directory structure
var Templates = make(map[string]*template.Template)

// var Template *template.Template

func init() {
	var temDir = "templates"
	tems, _ := ioutil.ReadDir(temDir)
	if len(tems) == 0 {
		temDir = "../templates"
	}
	setTemp(temDir, "")
	fmt.Println("\n==============================================================")
	fmt.Println("                       Templates parsed                         ")
	fmt.Println("==============================================================")
	// fmt.Println("templates", Templates["templates"].DefinedTemplates())
	// fmt.Println(Templates)
}

func setTemp(dirName string, subDir string) {
	newDirName := dirName + subDir
	fileInfo, _ := os.Lstat(newDirName)
	if !fileInfo.IsDir() {
		// fmt.Println(newDirName)
		// if Template == nil {
		// 	Template = template.Must(template.New(subDir).ParseFiles(newDirName))
		// } else {
		// 	Template = template.Must(Template.New(subDir).Parse(newDirName))
		// }
		value, ok := Templates[dirName]
		if ok {
			value = template.Must(value.ParseFiles(newDirName))
		} else {
			value = template.Must(template.ParseFiles(newDirName))
		}
		Templates[dirName] = value
	} else {
		filePtr, _ := os.Open(newDirName)
		defer filePtr.Close()
		fileNames, _ := filePtr.Readdirnames(-1)
		for _, fileName := range fileNames {
			setTemp(newDirName, "/"+fileName)
		}
	}
}

// func returnFileName(path string) []string {
// 	file, _ := os.Open(mainDir)
// 	strs, _ := file.Readdir(-1)
// 	var names []string
// 	for _, str := range strs {
// 		names = append(names, str.Name())
// 	}
// 	fmt.Println(names)
// }

// func findAndParseTemplates(rootDir string) (*template.Template, error) {
// 	cleanRoot := filepath.Clean(rootDir)
// 	fmt.Println("cleanRoot", cleanRoot)
// 	pfx := len(cleanRoot) + 1
// 	root := template.New("")

// 	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
// 		if !info.IsDir() && strings.HasSuffix(path, ".html") {
// 			if e1 != nil {
// 				return e1
// 			}

// 			b, e2 := ioutil.ReadFile(path)
// 			if e2 != nil {
// 				return e2
// 			}

// 			name := path[pfx:]
// 			t := root.New(name)
// 			t, e2 = t.Parse(string(b))
// 			if e2 != nil {
// 				return e2
// 			}
// 		}

// 		return nil
// 	})

// 	return root, err
// }
