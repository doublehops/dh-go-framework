# Validation

This validation package is designed to make adding validation to API endpoints as clear and
easy to use as possible. Although common validation rules exist for use, it's easy to create
custom ones as well as custom error messages.

# TODO
- See if it's possible to add required validation into `RunValidation` function.
- Add CI/CD to GitHub rules - `golint` and automated tests.

## Example Usage

You define the rules for each of the struct properties that need to be tested. Each
validation function should be placed comma-separated in the `Function` property of the
given rule.

```go
import "github.com/doublehops/dhapi/validator"

type Person struct {
    Name         string
    Age          string
    EmailAddress string
}

func main() {

    person := Person{
        Name:         "Jo",
        Age:          "Smith",
        EmailAddress: "jo.smith",
    }

    rules := []Rule{
        {"name", person.Name, false, []validator.ValidationFuncs{validator.MinLength(13, "")}},
        {"emailAddress", person.EmailAddress, false, []validator.ValidationFuncs{validator.EmailAddress("My custom error message")}},
    }
    
    errors := RunValidation(rules)
    j, _ := json.Marshal(errors)
    fmt.Println(string(j))
}
```

The response will contain the errors as an array per property as multiple rules for
each could fail. It should be easy for any frontend to consume. An example would look like this:
```json
{
  "emailAddress": [
    "My custom error message"
  ],
  "name": [
    "is not the minimum length"
  ]
}
```

## Normalising Before Validating

Models can define a `Normalise()` method that tidies input (trim whitespace, lower-case, de-duplicate)
before `Validate()` is called, so the rules check the value that will actually be stored. Handlers
call it straight after decoding the request:

```go
record.Normalise()

if errs := record.Validate(); len(errs) > 0 { ... }
```

## Checking a Referenced Record Exists

`RecordExistsByID(model, errorMessage)` checks that a non-soft-deleted row with the given ID exists.
The table name is the snake_case of the model's type name (`author.Author{}` → `author`,
`widgetpart.WidgetPart{}` → `widget_part`).

```go
{"authorId", b.AuthorID, true, []validator.ValidationFuncs{validator.RecordExistsByID(author.Author{}, "")}},
```

It needs a database connection, set once at startup with `validator.Init(db)` (see
`cmd/server/run.go`). Until then it fails with `ProcessingPropertyError`.

## Adding a Custom Validation Rule

You can create custom validation functions and use them just the same. The function
must confirm to the `ValidateFuncs` signature. Create the function and then assign
it to the variables in the `rules` slice like the pre-created functions.

An example would be:
```go
func CustomValidation(errorMessage string) ValidationFunctions {
	return func(required bool, value interface{}) (bool, string) {
		if errorMessage == "" {
			errorMessage = defaultErrorMessage
		}

		var v string
		var ok bool

		if v, ok = value.(string); !ok {
			return false, errorMessage
		}

		if v == "my-custom-value" {
			return true, ""
		}

		return false, errorMessage
	}
}
```