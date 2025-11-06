package employee

import(
	"context"

	"github.com/champion19/flighthours-api/core/domain"
)

func (r *repository) DeleteEmployee(ctx context.Context, id string) error {

tx, err := r.db.BeginTx(context.Background(), nil)
if err != nil {
	return err
}

	result,err:= tx.ExecContext(ctx,QueryDelete,id)
  if err!=nil{
		tx.Rollback()
		return domain.ErrUserCannotSave
	}

	rowsAffected,err:=result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return domain.ErrUserCannotSave
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return domain.ErrPersonNotFound
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
