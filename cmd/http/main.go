package main

import (
	"log"
	"net/http"
	"time"

	"github.com/iam43x/interview-help-api/internal/api"
	"github.com/iam43x/interview-help-api/internal/config"
	"github.com/iam43x/interview-help-api/internal/encode"
	"github.com/iam43x/interview-help-api/internal/gpt"
	"github.com/iam43x/interview-help-api/internal/jwt"
	"github.com/iam43x/interview-help-api/internal/store"
	"github.com/iam43x/interview-help-api/internal/util"

	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cnf := config.LoadConfig()
	encoder := encode.NewEncoder(cnf)
	gptClient := gpt.NewChatGptClient(cnf.ApiKey)

	/** JWT */
	rsa, err := jwt.ParseRSAPair(cnf)
	if err != nil {
		log.Fatal(err.Error())
	}
	issuer := jwt.NewIssuer(rsa, 10*time.Hour)
	jwtMiddle := jwt.NewJWTMiddleware(issuer)

	/** API */
	qr := api.NewQueryAPI(gptClient)
	tr := api.NewTranscriptor(encoder, gptClient)
	s, err := store.NewStore("db/user.db")
	if err != nil {
		log.Fatalf("ошибки инициализации sqlite3 %s", err.Error())
	}
	l := api.NewLoginAPI(issuer, s)
	rg := api.NewRegistrationAPI(issuer, s)

	r := mux.NewRouter()

	r.Handle("/v1/transcribe", util.EnableCORS(jwtMiddle.Handle(http.HandlerFunc(tr.TranscribeHttpHandler)))).Methods("POST")
	r.Handle("/v1/ask", util.EnableCORS(jwtMiddle.Handle(http.HandlerFunc(qr.QueryHttpHandler)))).Methods("POST")

	r.Handle("/register", util.EnableCORS(http.HandlerFunc(rg.RegistrationHttpHandler))).Methods("POST")
	r.Handle("/login", util.EnableCORS(http.HandlerFunc(l.LoginHttpHandler))).Methods("POST")

	server := &http.Server{
		Addr:           "localhost:8080",
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	cancel := util.GracefulShutdownServer(server)
	defer cancel()
	log.Printf("Server start on [%s]!", server.Addr)

	err = server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %s\n", err)
	}
}
