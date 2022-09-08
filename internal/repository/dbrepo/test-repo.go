package dbrepo

import (
	"errors"
	"time"

	"github.com/pierbusdev/conferenceRoomBookings/internal/models"
)

func (rep *testDBRepo) AllUsers() bool {
	return true
}

func (rep *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.OfficeID != 1 { //for tests
		return 0, errors.New("some error for testing")
	}
	return 1, nil
}

func (rep *testDBRepo) InsertOfficeRestriction(restr models.OfficeRestriction) error {
	if restr.OfficeID != 1 {
		return errors.New("some error for testing")
	}
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
	if id > 2 {
		return office, errors.New("some error for testing")
	}
	return office, nil
}
