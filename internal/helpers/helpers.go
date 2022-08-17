package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/pierbusdev/conferenceRoomBookings/internal/config"
)

var appConfig *config.AppConfig

func NewHelpers(config *config.AppConfig) {
	appConfig = config
}

func ClientError(w http.ResponseWriter, status int) {
	appConfig.InfoLog.Printf("Client error with status of %d\n", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	appConfig.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
