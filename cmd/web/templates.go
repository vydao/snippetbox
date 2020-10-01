package main

import (
    "html/template"
    "net/url"
    "path/filepath"
    "time"

    "vy.dao/snippetbox/pkg/forms"
    "vy.dao/snippetbox/pkg/models"
)

type templateData struct {
    CurrentYear   int
    Form          *forms.Form
    Snippet       *models.Snippet
    Snippets      []*models.Snippet
}

func humanDate(t time.Time) string {
    if t.IsZero() {
        return ""
    }

    return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
    "humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
    cache := map[string]*template.Template{}

    pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))

    if err != nil {
      return nil, err
    }

    for _, page := range pages {
      name := filepath.Base(page)

      template, err := template.New(name).Funcs(functions).ParseFiles(page)

      if err != nil {
        return nil, err
      }

      template, err = template.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))

      if err != nil {
        return nil, err
      }

      template, err = template.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
      if err != nil {
        return nil, err

      }
      cache[name] = template
    }
    return cache, nil
}