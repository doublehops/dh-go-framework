package scaffold

import (
	"context"
	"fmt"
)

const (
	serviceTemplate           = "./internal/scaffold/templates/service.tmpl"
	serviceRecordTemplate     = "./internal/scaffold/templates/servicerecord.tmpl"
	serviceRepositoryTemplate = "./internal/scaffold/templates/servicerepository.tmpl"
)

// createService will create the service, its record wrapper and the repository interface it
// depends on.
func (s *Scaffold) createService(ctx context.Context, m Model) error {
	m.ModelStructProperties = getStructProperties(m.Columns)
	path := fmt.Sprintf("%s/%s/%s", s.pwd, s.Paths.Service, m.LowerCase+"service")

	err := MkDir(path)
	if err != nil {
		s.l.Error(ctx, err.Error(), nil)

		return err
	}

	files := map[string]string{
		serviceTemplate:           fmt.Sprintf("%s/%sservice.go", path, m.LowerCase),
		serviceRecordTemplate:     fmt.Sprintf("%s/%srecord.go", path, m.LowerCase),
		serviceRepositoryTemplate: path + "/repository.go",
	}

	for src, filename := range files {
		if err = s.writeFile(src, filename, m); err != nil {
			s.l.Error(ctx, err.Error(), nil)

			return err
		}

		s.l.Info(ctx, "service file has been written: "+filename, nil)
	}

	return nil
}
