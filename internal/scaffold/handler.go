package scaffold

import (
	"context"
	"fmt"
)

const handlerTemplate = "./internal/scaffold/templates/handler.tmpl"

// createHandler will create the handler file for the model.
func (s *Scaffold) createHandler(ctx context.Context, m Model) error {
	m.ModelStructProperties = getStructProperties(m.Columns)
	m.ValidationRules = s.getValidationRules(m)
	path := fmt.Sprintf("%s/%s/%s", s.pwd, s.Paths.Handlers, m.LowerCase)
	filename := fmt.Sprintf("%s/%s.go", path, m.LowerCase)

	err := MkDir(path)
	if err != nil {
		s.l.Error(ctx, err.Error(), nil)

		return err
	}

	err = s.writeFile(handlerTemplate, filename, m)
	if err != nil {
		s.l.Error(ctx, err.Error(), nil)

		return err
	}

	s.l.Info(ctx, "handler has been written: "+filename, nil)

	return nil
}
