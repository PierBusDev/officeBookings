package dbrepo

import (
	"context"
	"time"

	"github.com/pierbusdev/conferenceRoomBookings/internal/models"
)

func (rep *postgresDBRepo) AllUsers() bool {
	return true
}

func (rep *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	//checking to avoid hanging transactions
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int
	query := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, office_id, created_at, updated_at) 
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	err := rep.DB.QueryRowContext(ctx, query,
		res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.OfficeID, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (rep *postgresDBRepo) InsertOfficeRestriction(restr models.OfficeRestriction) error {
	//checking to avoid hanging transactions
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `insert into office_restrictions (start_date, end_date, office_id, reservation_id, created_at, updated_at, restriction_id) 
				values ($1, $2, $3, $4, $5, $6, $7)`
	_, err := rep.DB.ExecContext(ctx, query,
		restr.StartDate, restr.EndDate, restr.OfficeID, restr.ReservationID, time.Now(), time.Now(), restr.RestrictionID)
	if err != nil {
		return err
	}
	return nil
}
