package scaffold

import (
	"context"
	"fmt"
)

const serviceTemplate = "./internal/scaffold/templates/service.tmpl"

// createService will create the service.
func (s *Scaffold) createService(ctx context.Context, m Model) error {
	m.ModelStructProperties = getStructProperties(m.Columns)
	path := fmt.Sprintf("%s/%s/%s", s.pwd, s.Paths.Service, m.LowerCase+"service")
	filename := fmt.Sprintf("%s/%sservice.go", path, m.LowerCase)

	err := MkDir(path)
	if err != nil {
		s.l.Error(ctx, err.Error(), nil)

		return err
	}

	err = s.writeFile(serviceTemplate, filename, m)
	if err != nil {
		s.l.Error(ctx, err.Error(), nil)

		return err
	}

	s.l.Info(ctx, "service has been written: "+filename, nil)

	return nil
}
