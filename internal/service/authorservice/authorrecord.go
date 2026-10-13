package authorservice

import (
	"context"
	"errors"

	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/model/author"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
	"github.com/doublehops/dh-go-framework/internal/service/userservice"
)

// AuthorRecord wraps an Author model with service access, so callers act on the record itself
// (record.Update(ctx, payload)) rather than passing the model back into the service.
type AuthorRecord struct {
	Record *author.Author
	svc    *AuthorService
}

// CanBeAccessedBy reports whether the requesting user may read or change this author: its
// creator, or a user with admin rights.
func (r *AuthorRecord) CanBeAccessedBy(u *userservice.UserRecord) bool {
	if u == nil {
		return false
	}

	return u.HasAdminRights() || r.Record.CreatedBy == u.Record.ID
}

// Create persists a new author record and reloads it from the database.
func (r *AuthorRecord) Create(ctx context.Context) (*AuthorRecord, error) {
	if err := r.Record.SetCreated(ctx); err != nil {
		r.svc.Log.Error(ctx, "error in SetCreated", logga.KVPs{"error": err.Error()})

		return nil, err
	}

	tx := r.svc.DB.MustBegin()
	defer tx.Rollback() //nolint:errcheck

	if err := r.svc.authorRepo.Create(ctx, tx, r.Record); err != nil {
		r.svc.Log.Error(ctx, service.UnableToSaveRecord+" "+err.Error(), nil)

		return nil, req.ErrCouldNotSaveRecord
	}

	if err := tx.Commit(); err != nil {
		r.svc.Log.Error(ctx, service.UnableToCommitTransaction+" "+err.Error(), nil)

		return nil, req.ErrCouldNotSaveRecord
	}

	if err := r.svc.authorRepo.GetByID(ctx, r.svc.DB, r.Record.ID, r.Record); err != nil {
		r.svc.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)

		return nil, errors.New(service.UnableToRetrieveRecord)
	}

	return r, nil
}

// Update copies the editable fields from payload onto the record and saves it. Only fields
// copied here can change — the ID and audit fields in the payload are ignored.
func (r *AuthorRecord) Update(ctx context.Context, payload *author.Author) error {
	r.Record.Name = payload.Name
	r.Record.SetUpdated(ctx)

	tx := r.svc.DB.MustBegin()
	defer tx.Rollback() //nolint:errcheck

	if err := r.svc.authorRepo.Update(ctx, tx, r.Record); err != nil {
		r.svc.Log.Error(ctx, service.UnableToUpdateRecord+" "+err.Error(), nil)

		return req.ErrCouldNotSaveRecord
	}

	if err := tx.Commit(); err != nil {
		r.svc.Log.Error(ctx, service.UnableToCommitTransaction+" "+err.Error(), nil)

		return req.ErrCouldNotSaveRecord
	}

	if err := r.svc.authorRepo.GetByID(ctx, r.svc.DB, r.Record.ID, r.Record); err != nil {
		r.svc.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)

		return errors.New(service.UnableToRetrieveRecord)
	}

	return nil
}

// Delete soft-deletes the author record.
func (r *AuthorRecord) Delete(ctx context.Context) error {
	r.Record.SetDeleted(ctx)

	tx := r.svc.DB.MustBegin()
	defer tx.Rollback() //nolint:errcheck

	if err := r.svc.authorRepo.Delete(ctx, tx, r.Record); err != nil {
		r.svc.Log.Error(ctx, "unable to delete record. "+err.Error(), nil)

		return req.ErrCouldNotSaveRecord
	}

	if err := tx.Commit(); err != nil {
		r.svc.Log.Error(ctx, service.UnableToCommitTransaction+" "+err.Error(), nil)

		return req.ErrCouldNotSaveRecord
	}

	return nil
}
