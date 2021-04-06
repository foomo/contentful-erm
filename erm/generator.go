package erm

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/foomo/gocontentful/erm/templates"
	"golang.org/x/tools/imports"
)

func formatAndFixImports(filename string) error {
	sourceBytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	formattedSource, err := format.Source(sourceBytes)
	if err != nil {
		return err
	}
	finalSource, err := imports.Process(filename, formattedSource, nil)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, finalSource, 0644)
}

func generate(filename string, tpl []byte, conf spaceConf) error {
	fmt.Println("Processing", filename)
	tmpl, err := template.New("generate-" + filename).Funcs(conf.FuncMap).Parse(string(tpl))
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(f, conf); err != nil {
		return err
	}
	if err := formatAndFixImports(filename); err != nil {
		return err
	}

	return nil
}

// generateCode generates API to and value objects for the space
func generateCode(conf spaceConf) error {
	for file, tpl := range map[string][]byte{
		filepath.Join(conf.PackageDir, "gocontentfulvobase"+goExt): templates.TemplateVoBase,
		filepath.Join(conf.PackageDir, "gocontentfulvo"+goExt):     templates.TemplateVo,
		filepath.Join(conf.PackageDir, "gocontentfulvolib"+goExt):  templates.TemplateVoLib,
	} {
		if err := generate(file, tpl, conf); err != nil {
			return err
		}
	}
	for _, contentType := range conf.ContentTypes {
		conf.ContentType = contentType
		if err := generate(
			filepath.Join(conf.PackageDir, "gocontentfulvolib"+strings.ToLower(contentType.Sys.ID)+goExt),
			templates.TemplateVoLibContentType,
			conf,
		); err != nil {
			return err
		}
	}
	return nil
}
