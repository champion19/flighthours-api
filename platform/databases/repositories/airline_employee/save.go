package airline_employee

import (
	"context"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/go-sql-driver/mysql"
)

// SaveAirlineEmployee creates a new airline employee in the database
func (r *repository) SaveAirlineEmployee(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
	employeeToSave := FromDomain(&employee)

	// Cast the transaction to the concrete type
	dbTx, ok := tx.(*common.SQLTX)
	if !ok {
		log.Error(logger.LogDatabaseUnavailable, "error", "invalid transaction type")
		return domain.ErrInvalidTransaction
	}

	_, err := dbTx.ExecContext(ctx, QueryInsert,
		employeeToSave.ID,
		employeeToSave.Name,
		employeeToSave.AirlineID,
		employeeToSave.Email,
		employeeToSave.IdentificationNumber,
		employeeToSave.Bp,
		employeeToSave.StartDate,
		employeeToSave.EndDate,
		employeeToSave.Active,
		employeeToSave.Role,
		employeeToSave.KeycloakUserID,
	)

	if err != nil {
		// Check for specific MySQL errors
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1452:
				// Foreign key constraint fails (e.g., invalid airline)
				log.Error(logger.LogDatabaseUnavailable,
					"employee_id", employee.ID,
					"error", "invalid foreign key reference",
					"mysql_error", mysqlErr.Message)
				return domain.ErrInvalidForeignKey
			case 1062:
				// Duplicate entry
				log.Error(logger.LogDatabaseUnavailable,
					"employee_id", employee.ID,
					"error", "duplicate entry",
					"mysql_error", mysqlErr.Message)
				return domain.ErrDuplicateUser
			}
		}
		log.Error(logger.LogDatabaseUnavailable, "employee_id", employee.ID, "error", err)
		return err
	}

	return nil
}
