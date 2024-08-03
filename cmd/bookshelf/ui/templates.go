package ui

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/r3d5un/Bookshelf/internal/books/types"
)

//go:embed "html" "static" "static"
var Files embed.FS

// Loads templates from the embedded filesystem containing templates assets.
func (m *Module) newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	m.logger.Info("listing HTML directory")
	tmplFiles, err := fs.Glob(Files, "html/*/*.tmpl")
	if err != nil {
		m.logger.Error("an error occurred while walking template directory", "error", err)
		return nil, err
	}

	m.logger.Info("adding template files to cache")
	for _, tmplFile := range tmplFiles {
		name := filepath.Base(tmplFile)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			tmplFile,
		}

		m.logger.Info("parsing template files", "templateFile", tmplFile, "patterns", patterns)
		template, err := template.New(name).Funcs(functions).ParseFS(Files, patterns...)
		if err != nil {
			m.logger.Error("unable to parse templates", "error", err)
			return nil, err
		}
		m.logger.Info("tempalte parsed")

		cache[name] = template
	}

	return cache, nil
}

// Contains functions that the templates can call internally
var functions = template.FuncMap{
	"humanDate": humanDate,
	"sub": func(a, b int) int {
		return a - b
	},
	"paragraphify": paragraphify,
}

type templateData struct {
	MyLibraryBooks       []*types.Book         `json:"myLibraryBooks,omitempty"`
	SelectedCategory     string                `json:"selectedCategory,omitempty"`
	BookSeriesAccordions []bookSeriesAccordion `json:"bookSeriesAccordions,omitempty"`
	BookData             types.Book            `json:"bookData,omitempty"`
}

type bookSeriesAccordion struct {
	ID          string     `json:"id"`
	Order       float32    `json:"order"`
	Title       string     `json:"title"`
	Published   *time.Time `json:"published,omitempty"`
	Description string     `json:"description"`
	Selected    bool       `json:"selected"`
}

type myLibraryBook struct {
	Title     string     `json:"title"`
	Series    []*string  `json:"series,omitempty"`
	Authors   []*string  `json:"authors,omitempty"`
	Published *time.Time `json:"published,omitempty"`
	Added     *time.Time `json:"added,omitempty"`
	Status    string     `json:"string,omitempty"`
}

var (
	ErrTemplateNotFound = errors.New("template not found")
)

func (m *Module) render(
	w http.ResponseWriter,
	status int,
	templateName string,
	data *templateData,
) {
	templates, ok := m.templateCache[templateName]
	if !ok {
		m.logger.Info(
			"error occurred when loading template from cache",
			"error", ErrTemplateNotFound,
		)
		m.serverError(w, ErrTemplateNotFound)
		return
	}

	buffer := new(bytes.Buffer)

	err := templates.ExecuteTemplate(buffer, "base", data)
	if err != nil {
		m.logger.Error(
			"error occurred while executing template",
			"error", err,
		)
		m.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buffer.WriteTo(w)
}

func (m *Module) renderPartial(
	w http.ResponseWriter,
	status int,
	templateName string,
	data *templateData,
) {
	templates, ok := m.templateCache[templateName]
	if !ok {
		m.logger.Info(
			"error occurred when loading template from cache",
			"error", ErrTemplateNotFound,
		)
		m.serverError(w, ErrTemplateNotFound)
		return
	}

	buffer := new(bytes.Buffer)

	err := templates.Execute(buffer, data)
	if err != nil {
		m.logger.Error(
			"error occurred while executing template",
			"error", err,
		)
		m.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buffer.WriteTo(w)
}

// Writes a response to a request with a status and the body as a string. Intended for use in smaller
// responses where a dedicated template is overkill.
func (m *Module) rawResponse(w http.ResponseWriter, status int, responseBody string) {
	w.WriteHeader(status)
	w.Write([]byte(responseBody))
}

func (m *Module) serverError(w http.ResponseWriter, err error) {
	m.logger.Error("a server error occurred", "error", err)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("2006-01-02")
}
func paragraphify(text string) template.HTML {
	paragraphs := strings.Split(text, "\n\n")
	var result string
	for _, para := range paragraphs {
		// Trim any surrounding whitespace and add paragraph tags
		trimmedPara := strings.TrimSpace(para)
		if trimmedPara != "" {
			result += "<p>" + template.HTMLEscapeString(trimmedPara) + "</p>\n"
		}
	}
	return template.HTML(result)
}
