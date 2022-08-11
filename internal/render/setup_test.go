package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pierbusdev/conferenceRoomBookings/internal/config"
	"github.com/pierbusdev/conferenceRoomBookings/internal/models"
)

var session *scs.SessionManager
var testAppConfig config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	//creating app config
	//TODO change this to true when in production
	testAppConfig.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true //cookie must persist after user closes the browser
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testAppConfig.Session = session

	appConfig = &testAppConfig //appConfig is in render.go
	os.Exit(m.Run())
}

// myWriter is a custom writer that implements the http.ResponseWriter interface.
type myWriter struct{}

func (w myWriter) Header() http.Header {
	return http.Header{}
}

func (w *myWriter) WriteHeader(int) {}

func (w *myWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
