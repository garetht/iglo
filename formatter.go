package iglo

import (
	"bitbucket.org/pkg/inflect"
	"bytes"
	"fmt"
	bf "github.com/russross/blackfriday"
	"html/template"
	"io"
	"strings"
)

func labelize(method string) string {
	switch method {
	case "GET":
		return "info"
	case "POST":
		return "success"
	case "PUT":
		return "primary"
	case "PATCH":
		return "warning"
	case "DELETE":
		return "danger"
	}

	return "default"
}

func iconize(method string) string {
	switch method {
	case "GET":
		return "arrow-circle-down"
	case "POST":
		return "plus-square"
	case "PUT":
		return "pencil"
	case "PATCH":
		return "warning"
	case "DELETE":
		return "dot-circle-o"
	}

	return "default"
}

func markdownize(str string) template.HTML {
	b := bf.MarkdownCommon([]byte(str))
	return template.HTML(string(b))
}

func HTML(w io.Writer, api *API, templatePaths []string) error {
	// template functions
	funcMap := template.FuncMap{
		"dasherize":   inflect.Dasherize,
		"trim":        strings.Trim,
		"labelize":    labelize,
		"markdownize": markdownize,
		"iconize":     iconize,
	}
	tmpl, err := template.New("html").Funcs(funcMap).ParseFiles(templatePaths...)

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = tmpl.ExecuteTemplate(w, "Blueprint", api)

	if err != nil {
		return err
	}

	return nil
}

func MarkdownToHTMLWithTemplate(w io.Writer, r io.Reader, templatePaths []string) error {
	data, err := ParseMarkdown(r)
	if err != nil {
		return err
	}

	err = JSONToHTML(w, bytes.NewBuffer(data), templatePaths)
	if err != nil {
		return err
	}

	return nil
}

func MarkdownToHTML(w io.Writer, r io.Reader) error {
	filepath := []string{"../../templates/flatly.html"}
	return MarkdownToHTMLWithTemplate(w, r, filepath)
}

func JSONToHTML(w io.Writer, r io.Reader, templatePaths []string) error {
	api, err := ParseJSON(r)
	if err != nil {
		return err
	}

	err = HTML(w, api, templatePaths)
	if err != nil {
		return err
	}

	return nil
}
