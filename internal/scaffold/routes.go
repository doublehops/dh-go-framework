package scaffold

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"text/template"
)

const routeTemplate = "./internal/scaffold/templates/routes.tmpl"

// printRoutes will print the routes for the user to manually add to the routes table.
func (s *Scaffold) printRoutes(ctx context.Context, m Model) error {
	m.ModelStructProperties = getStructProperties(m.Columns)

	f, err := os.Open(routeTemplate)
	if err != nil {
		e := "unable to open template. " + err.Error()
		s.l.Error(ctx, e, nil)

		return errors.New(e)
	}
	defer func() { _ = f.Close() }()

	source, err := io.ReadAll(f)
	if err != nil {
		e := "unable to read file." + err.Error()
		s.l.Error(ctx, e, nil)

		return errors.New(e)
	}

	buf := bytes.Buffer{}

	t, err := template.New("model").Parse(string(source))
	if err != nil {
		e := "unable to parse template." + err.Error()
		s.l.Error(ctx, e, nil)

		return errors.New(e)
	}
	err = t.Execute(&buf, m)
	if err != nil {
		e := "unable to write new file. " + err.Error()
		s.l.Error(ctx, e, nil)

		return errors.New(e)
	}

	if _, err := os.Stdout.Write(buf.Bytes()); err != nil {
		e := "unable to write to stdout. " + err.Error()
		s.l.Error(ctx, e, nil)

		return errors.New(e)
	}

	return nil
}
