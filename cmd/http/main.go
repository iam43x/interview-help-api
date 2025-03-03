package main

import (
	"log"
	"net/http"
	"time"

	"github.com/iam43x/interview-help-4u/internal/api"
	"github.com/iam43x/interview-help-4u/internal/config"
	"github.com/iam43x/interview-help-4u/internal/encode"
	"github.com/iam43x/interview-help-4u/internal/gpt"
	"github.com/iam43x/interview-help-4u/internal/jwt"
	"github.com/iam43x/interview-help-4u/internal/util"

	"github.com/gorilla/mux"
)

func main() {
	cnf := config.LoadConfig()	
	encoder := encode.NewEncoder(cnf)
	gptClient := gpt.NewChatGptClient(cnf.ApiKey)
	
	/** JWT */
	rsa, err := jwt.ParseRSAPair(cnf)
	if err != nil {
		log.Fatalf(err.Error())
	}
	issuer := jwt.NewIssuer(rsa, 10 * time.Hour)
	jwtMiddle := jwt.NewJWTMiddleware(issuer)

	/** API */
	tr := api.NewTranscriptor(encoder, gptClient)
	qr := api.NewQueryAPI(gptClient)
	l := api.NewLoginAPI(issuer)

	r := mux.NewRouter()

	r.Handle("/v1/transcribe", util.EnableCORS(jwtMiddle.Handle(http.HandlerFunc(tr.TranscribeHttpHandler)))).Methods("POST")
	r.Handle("/v1/ask", util.EnableCORS(jwtMiddle.Handle(http.HandlerFunc(qr.QueryHttpHandler)))).Methods("POST")

	r.HandleFunc("/login", l.LoginHttpHandler).Methods("POST")

	server := &http.Server{
		Addr: "localhost:8080",
		Handler: r,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
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