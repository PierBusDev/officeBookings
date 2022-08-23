package repository

import (
	"time"

	"github.com/pierbusdev/conferenceRoomBookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertOfficeRestriction(restr models.OfficeRestriction) error
	SearchAvailabilityByDatesByOfficeId(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllOffices(start, end time.Time) ([]models.Office, error)
}
