package main

import (
	"log"
	"net/http"

	"cognito-example/config"
	"cognito-example/handlers"
	"cognito-example/middleware"
)

func main() {
	cfg := config.LoadConfig()
	authHandler, err := handlers.NewAuthHandler(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc(
		"/login",
		middleware.ChainMiddleWares(
			authHandler.Login,
			middleware.CreateAuthMiddleware(),
			middleware.CreateLoggerMiddleware(),
			middleware.CreateCORSMiddleware(),
		),
	)
	http.HandleFunc(
		"/callback",
		middleware.ChainMiddleWares(
			authHandler.Callback,
			middleware.CreateAuthMiddleware(),
			middleware.CreateLoggerMiddleware(),
			middleware.CreateCORSMiddleware(),
		),
	)
	http.HandleFunc(
		"/profile",
		middleware.ChainMiddleWares(
			authHandler.Profile,
			middleware.CreateAuthMiddleware(),
			middleware.CreateLoggerMiddleware(),
			middleware.CreateCORSMiddleware(),
		),
	)

	// Serve static home page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
				<body>
					<h1>Cognito OAuth2.0 Example</h1>
					<a href="/login">Login with Cognito</a>
				</body>
			</html>
		`))
	})

	log.Printf("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}
