package dbrepo

import (
	"context"
	"time"

	"github.com/pierbusdev/conferenceRoomBookings/internal/models"
)

func (rep *postgresDBRepo) AllUsers() bool {
	//TODO to update as soon as support for multiple users is ready
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

//SearchAvailabilityByDatesByOfficeId returns true if the availability exists for officeId and false if not
func (rep *postgresDBRepo) SearchAvailabilityByDatesByOfficeId(start, end time.Time, officeId int) (bool, error) {
	//checking to avoid hanging transactions
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var nRows int
	query := `
		select count(*) 
		from office_restrictions 
		where office_id = $1 and
		$2 < end_date and $3 > start_date;`
	row := rep.DB.QueryRowContext(ctx, query, officeId, start, end)

	if err := row.Scan(&nRows); err != nil {
		return false, err
	}

	if nRows == 0 {
		return true, nil
	}
	return false, nil
}

//SearchAvailabilityForAllOffices returns a slice of available offices if there are any in the date specified
func (rep *postgresDBRepo) SearchAvailabilityForAllOffices(start, end time.Time) ([]models.Office, error) {
	//checking to avoid hanging transactions
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		select id, office_name
		from offices 
		where id not in (
			select office_id 
			from office_restrictions 
			where $1 < end_date and $2 > start_date
		);`

	var offices []models.Office
	rows, err := rep.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return offices, err
	}
	defer rows.Close()

	for rows.Next() {
		var office models.Office
		if err := rows.Scan(&office.ID, &office.OfficeName); err != nil {
			return offices, err
		}
		offices = append(offices, office)
	}

	if err := rows.Err(); err != nil {
		return offices, err
	}

	return offices, nil
}

func (rep *postgresDBRepo) GetOfficeById(id int) (models.Office, error) {
	//checking to avoid hanging transactions
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `select id, office_name, created_at, updated_at from offices where id = $1`
	var office models.Office
	row := rep.DB.QueryRowContext(ctx, query, id)
	if err := row.Scan(&office.ID, &office.OfficeName, &office.CreatedAt, &office.UpdatedAt); err != nil {
		return office, err
	}
	return office, nil
}
