package dbrepo

import (
	"database/sql"

	"github.com/pierbusdev/conferenceRoomBookings/internal/config"
	"github.com/pierbusdev/conferenceRoomBookings/internal/repository"
)

type postgresDBRepo struct {
	AppConfig *config.AppConfig
	DB        *sql.DB
}

type testDBRepo struct {
	AppConfig *config.AppConfig
	DB        *sql.DB
}

func NewPostgresDBRepo(conn *sql.DB, appConfig *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		AppConfig: appConfig,
		DB:        conn,
	}
}

func NewTestingRepo(appConfig *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		AppConfig: appConfig,
	}
}
