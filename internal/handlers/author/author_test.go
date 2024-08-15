package author

import (
	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/httprequest"
	"github.com/doublehops/dh-go-framework/internal/logga"
	"testing"
)

func TestMain(m *testing.M) {

}

func TestCRUD(t *testing.T) {
	logg := &config.Logging{
		Writer:       "",
		LogLevel:     "debug",
		OutputFormat: "JSON",
	}
	l, err := logga.New(logg)

	req := httprequest.Requester{}
}
