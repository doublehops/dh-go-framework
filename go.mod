module github.com/doublehops/dh-go-framework

go 1.24.0

// @todo - this should be removed when dhapi is pushed to Github.
replace github.com/doublehops/dhapi => /home/b/workspace/dhapi-2

require (
	github.com/doublehops/go-common v0.0.0-20230910011642-8556bd635e3f
	github.com/go-sql-driver/mysql v1.8.1
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/jinzhu/copier v0.4.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/mythrnr/httprouter-group v0.9.1
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.45.0
	golang.org/x/text v0.31.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
