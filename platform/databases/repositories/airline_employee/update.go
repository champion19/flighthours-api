package airline_employee

import (
	"context"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/go-sql-driver/mysql"
)

// UpdateAirlineEmployee updates an airline employee in the database
func (r *repository) UpdateAirlineEmployee(ctx context.Context, tx output.Tx, employee domain.AirlineEmployee) error {
	employeeToUpdate := FromDomain(&employee)

	log.Debug("UpdateAirlineEmployee: Starting update",
		"employee_id", employeeToUpdate.ID,
		"name", employeeToUpdate.Name,
		"airline_id", employeeToUpdate.AirlineID,
		"email", employeeToUpdate.Email,
		"active", employeeToUpdate.Active)

	// Cast the transaction to the concrete type
	dbTx, ok := tx.(*common.SQLTX)
	if !ok {
		log.Error(logger.LogDatabaseUnavailable, "error", "invalid transaction type")
		return domain.ErrInvalidTransaction
	}

	result, err := dbTx.ExecContext(ctx, QueryUpdate,
		employeeToUpdate.Name,
		employeeToUpdate.AirlineID,
		employeeToUpdate.Email,
		employeeToUpdate.IdentificationNumber,
		employeeToUpdate.Bp,
		employeeToUpdate.StartDate,
		employeeToUpdate.EndDate,
		employeeToUpdate.Active,
		employeeToUpdate.Role,
		employeeToUpdate.KeycloakUserID,
		employeeToUpdate.ID,
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

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrAirlineEmployeeNotFound
	}

	log.Debug("UpdateAirlineEmployee: Query executed successfully", "rows_affected", rowsAffected)

	return nil
}
