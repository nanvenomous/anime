package main

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

var (
	allowedOrigins = []string{
		"http://localhost",
	}
	devMode = true
)

func serveFiles(mux *http.ServeMux, filesPath string) {
	var enforceDevMode = func(w http.ResponseWriter) {
		if devMode {
			w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Add("Pragma", "no-cache")
			w.Header().Add("Expires", "0")
		}
	}

	var fileServer = http.FileServer(http.Dir("bundle"))
	mux.HandleFunc("/"+filesPath+"/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/"+filesPath)
		enforceDevMode(w)
		fileServer.ServeHTTP(w, r)
	})
}

func withErrorLogging(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func setHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	for _, alOrg := range allowedOrigins {
		if strings.Contains(origin, alOrg) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			setHeaders(w, r)
			return
		}
		setHeaders(w, r)
		next.ServeHTTP(w, r)
	})
}

func serve() error {
	var (
		mux     = http.NewServeMux()
		handler http.Handler
	)

	handler = mux
	CorsMiddleware(handler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var (
			tmplts = template.Must(template.ParseFiles(
				"templates/index.html",
			))
			tmpl = Index{
				SlidingText: &SlidingText{
					Words:  []string{"You", "Your child", "Your spouse", "Your family"},
					LeadIn: "At law for",
				},
				StyleSheets: styleSheets,
			}
		)

		withErrorLogging(tmplts.Execute(w, &tmpl))
	})

	serveFiles(mux, "bundle")

	srvr := &http.Server{
		Addr:    ":7676",
		Handler: handler,
	}

	return srvr.ListenAndServe()
}

func main() {
	if err := serve(); err != nil {
		panic(err)
	}
}
