package notifications

import (
	"errors"
	html "html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	text "text/template"

	"github.com/go-kit/kit/log"
	"github.com/markbates/pkger"
)

type TemplateRepository interface {
	Text(values Template) (string, error)
	HTML(values Template) (string, error)
}

type templateRepository struct {
	textTemplates text.Template
	htmlTemplates html.Template
}

func NewTemplateRepository(logger log.Logger) (TemplateRepository, error) {
	ht := html.New("notifications")
	tt := text.New("notifications")

	logger.Log("level", "info", "msg", "Loading templates")

	err := pkger.Walk("/configs/notifications/", func(path string, info os.FileInfo, err error) error {
		logger.Log("level", "info", "msg", "Walking", "name", info.Name(), "path", path, "ext", filepath.Ext(info.Name()))

		ext := strings.ToLower(filepath.Ext(info.Name()))

		switch ext {
		case ".txt":
			f, err := pkger.Open(path)
			if err != nil {
				return err
			}

			content, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}

			_, err = tt.New(info.Name()).Parse(string(content))
			if err != nil {
				return err
			}

			logger.Log("level", "info", "msg", "Loaded template - "+info.Name())
		case ".html":
			f, err := pkger.Open(path)
			if err != nil {
				return err
			}

			content, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}

			_, err = ht.New(info.Name()).Parse(string(content))
			if err != nil {
				return err
			}

			logger.Log("level", "info", "msg", "Loaded template - "+info.Name())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logger.Log("level", "info", "msg", "Loaded templates")

	return &templateRepository{
		textTemplates: *tt,
		htmlTemplates: *ht,
	}, nil
}

func (t *templateRepository) Text(values Template) (string, error) {
	template := t.textTemplates.Lookup(values.TemplateName() + ".txt")
	if template == nil {
		return "", errors.New("Unable to find template: " + values.TemplateName())
	}

	generated := strings.Builder{}

	err := template.Execute(&generated, values)
	if err != nil {
		return "", err
	}

	return generated.String(), nil
}

func (t *templateRepository) HTML(values Template) (string, error) {
	template := t.htmlTemplates.Lookup(values.TemplateName() + ".html")
	if template == nil {
		return "", errors.New("Unable to find template: " + values.TemplateName())
	}

	generated := strings.Builder{}

	err := template.Execute(&generated, values)
	if err != nil {
		return "", err
	}

	return generated.String(), nil
}
