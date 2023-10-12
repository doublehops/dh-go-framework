package service

import (
	"github.com/doublehops/dhapi-example/internal/handlers"
	"github.com/doublehops/dhapi-example/internal/model"
	"github.com/doublehops/dhapi-example/internal/repository/repositoryauthor"
)

type AuthorService struct {
	app        *handlers.App
	authorRepo *repositoryauthor.RepositoryAuthor
}

func New(app *handlers.App, authorRepo *repositoryauthor.RepositoryAuthor) *AuthorService {
	return &AuthorService{
		app:        app,
		authorRepo: authorRepo,
	}
}

func (s AuthorService) Create(author *model.Author) (model.Author, error) {
	return *author, nil
}
