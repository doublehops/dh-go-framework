module github.com/doublehops/dh-go-framework

go 1.21.0

// @todo - this should be removed when dhapi is pushed to Github.
replace github.com/doublehops/dhapi => /home/b/workspace/dhapi-2

require (
	github.com/doublehops/go-common v0.0.0-20230910011642-8556bd635e3f
	github.com/doublehops/go-migration v0.0.3
	github.com/go-sql-driver/mysql v1.8.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/mythrnr/httprouter-group v0.9.1
	github.com/stretchr/testify v1.9.0
	golang.org/x/text v0.16.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
