package dbrepo

import (
	"time"

	"github.com/pierbusdev/conferenceRoomBookings/internal/models"
)

func (rep *testDBRepo) AllUsers() bool {
	return true
}

func (rep *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

func (rep *testDBRepo) InsertOfficeRestriction(restr models.OfficeRestriction) error {
	return nil
}

//SearchAvailabilityByDatesByOfficeId returns true if the availability exists for officeId and false if not
func (rep *testDBRepo) SearchAvailabilityByDatesByOfficeId(start, end time.Time, officeId int) (bool, error) {
	return false, nil
}

//SearchAvailabilityForAllOffices returns a slice of available offices if there are any in the date specified
func (rep *testDBRepo) SearchAvailabilityForAllOffices(start, end time.Time) ([]models.Office, error) {
	var offices []models.Office
	return offices, nil
}

func (rep *testDBRepo) GetOfficeById(id int) (models.Office, error) {
	var office models.Office
	return office, nil
}
