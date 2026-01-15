package daily_logbook_detail

import (
	"context"
	"strings"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
	"github.com/champion19/flighthours-api/platform/logger"
)

// UpdateDailyLogbookDetail updates an existing daily logbook detail
func (r *repository) UpdateDailyLogbookDetail(ctx context.Context, tx output.Tx, detail domain.DailyLogbookDetail) error {
	log.Info(logger.LogDailyLogbookDetailUpdate, "data", detail.ToLogger())

	entity := FromDomain(&detail)

	sqlTx, ok := tx.(*common.SQLTX)
	if !ok {
		log.Error(logger.LogDailyLogbookDetailUpdateError, "error", "invalid transaction type")
		return domain.ErrInvalidTransaction
	}

	stmt := sqlTx.Tx.StmtContext(ctx, r.stmtUpdate)

	_, err := stmt.ExecContext(ctx,
		entity.FlightRealDate,
		entity.FlightNumber,
		entity.AirlineRouteID,
		entity.ActualAircraftRegistrationID,
		entity.Passengers,
		entity.OutTime,
		entity.TakeoffTime,
		entity.LandingTime,
		entity.InTime,
		entity.PilotRole,
		entity.CompanionName,
		entity.AirTime,
		entity.BlockTime,
		entity.DutyTime,
		entity.ApproachType,
		entity.FlightType,
		entity.ID, // WHERE clause
	)

	if err != nil {
		log.Error(logger.LogDailyLogbookDetailUpdateError, "error", err.Error())

		// Parse FK constraint errors for better error messages
		errStr := err.Error()
		if strings.Contains(errStr, "foreign key constraint") || strings.Contains(errStr, "1452") {
			if strings.Contains(errStr, "airline_route") {
				return domain.ErrFlightInvalidRoute
			}
			if strings.Contains(errStr, "aircraft_registration") {
				return domain.ErrFlightInvalidAircraft
			}
		}
		return domain.ErrFlightCannotUpdate
	}

	log.Info(logger.LogDailyLogbookDetailUpdateOK, "id", detail.ID)
	return nil
}
