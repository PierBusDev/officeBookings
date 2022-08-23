package models

import "time"

//Below all the models for the database
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Office struct {
	ID         int
	OfficeName string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	OfficeID  int
	CreatedAt time.Time
	UpdatedAt time.Time
	Office    Office
}

type OfficeRestriction struct {
	ID            int
	OfficeID      int
	ReservationID int
	RestrictionID int
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Office        Office
	Reservation   Reservation
	Restriction   Restriction
}
