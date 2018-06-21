package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var Base = filepath.Base("")
var Template *template.Template
var pathToTepmlates = filepath.Join(Base, "/templates/**/*")

func init() {
	var mainDir = "./templates"
	// dirs, err := ioutil.ReadDir(mainDir)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, dir := range dirs {
	// 	var dirName = mainDir + "/" + dir.Name()
	// 	files, _ := ioutil.ReadDir(dirName)

	// 	for _, file := range files {
	// 		Template = template.Must(template.New(mainDir + "/" + dir.Name() + "/" + file.Name()).ParseGlob(mainDir + "/" + dir.Name() + "/" + file.Name()))
	// 	}
	// 	// fmt.Println(file.Name())
	// }
	Template = findAndParseTemplates(mainDir, nil)
}

func findAndParseTemplates(rootDir string, funcMap template.FuncMap) (*template.Template, error) {
	cleanRoot := filepath.Clean(rootDir)
	pfx := len(cleanRoot) + 1
	root := template.New("")

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if e1 != nil {
				return e1
			}

			b, e2 := ioutil.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			t := root.New(name).Funcs(funcMap)
			t, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	return root, err
}
