package employee

import(
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)

func (r *repository) DeleteEmployee(id string) error {

tx, err := r.db.BeginTx(context.Background(), nil)
if err != nil {
	return err
}

	_,err = tx.ExecContext(context.Background(),QueryDelete,id)
  if err!=nil{
		tx.Rollback()
		return domain.ErrUserCannotSave
	}


	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
