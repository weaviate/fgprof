package fgprof

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Handler returns an http handler that takes an optional "seconds" query
// argument that defaults to "30" and produces a profile over this duration.
// The optional "format" parameter controls if the output is written in
// Google's "pprof" format (default) or Brendan Gregg's "folded" stack format.
func Handler(ignoreFunctions ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var seconds int
		var err error
		if s := r.URL.Query().Get("seconds"); s == "" {
			seconds = 30
		} else if seconds, err = strconv.Atoi(s); err != nil || seconds <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "bad seconds: %d: %s\n", seconds, err)
			return
		}

		var ignoreFunctionsQuery []string
		if i := r.URL.Query()["ignore"]; len(i) == 0 {
			// do nothing, empty array is fine
		} else {
			ignoreFunctionsQuery = i

		}

		format := Format(r.URL.Query().Get("format"))
		if format == "" {
			format = FormatPprof
		}

		stop := Start(w, format, append(ignoreFunctions, ignoreFunctionsQuery...))
		defer stop()
		time.Sleep(time.Duration(seconds) * time.Second)
	})
}
