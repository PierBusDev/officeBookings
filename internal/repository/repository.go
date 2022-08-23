package repository

import "github.com/pierbusdev/conferenceRoomBookings/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertOfficeRestriction(restr models.OfficeRestriction) error
}
