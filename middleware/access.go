package middleware

import (
	"net/http"
	"time"

	"fmt"

	"github.com/answer1991/daily-roll-logrus"
)

var (
	logger = drl.GetLogger("access")
)

func Access(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		wWrapper := newResponseWriterWrapper(w)
		handler.ServeHTTP(wWrapper, r)

		executeTime := time.Since(startTime)

		logger.
			WithField("url", r.RequestURI).
			WithField("method", r.Method).
			WithField("remoteAddr", r.RemoteAddr).
			WithField("code", wWrapper.code).
			WithField("cost", fmt.Sprintf("%v%s", float64(executeTime.Nanoseconds()/1e6), "ms")).
			Info("")
	})
}
