package notifications

import (
	"errors"
	html "html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	text "text/template"

	"github.com/markbates/pkger"
	log "github.com/moov-io/identity/pkg/logging"
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

	logger.Info().Log("Loading templates")

	err := pkger.Walk("/configs/notifications/", func(path string, info os.FileInfo, err error) error {

		ext := strings.ToLower(filepath.Ext(info.Name()))

		logCtx := logger.Info().WithMap(map[string]string{
			"name": info.Name(),
			"path": path,
			"ext":  ext,
		})

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

			logCtx.Log("Loaded template - " + info.Name())
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

			logCtx.Log("Loaded template - " + info.Name())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logger.Info().Log("Loaded templates")

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
