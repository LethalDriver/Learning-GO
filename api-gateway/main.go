package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

func main() {
	publicKey, err := getRsaPublicKey()
	if err != nil {
		log.Fatalf("Failed to load RSA public key: %v", err)
	}

	port := os.Getenv("PORT")

	authService := &AuthService{publicKey: publicKey}

	mux := http.NewServeMux()

	mux.Handle("/user-service/", http.StripPrefix("/user-service", proxyHandler("http://user-service:8081")))
	mux.Handle("/chat-service/", JWTMiddleware(authService, http.StripPrefix("/chat-service", proxyHandler("http://chat-service:8082"))))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Println("API Gateway listening on port 8080...")
	log.Fatal(server.ListenAndServe())
}

func proxyHandler(target string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url, err := url.Parse(target)
		if err != nil {
			http.Error(w, "Invalid target URL", http.StatusInternalServerError)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, fmt.Sprintf("Proxy error: %v", err), http.StatusBadGateway)
		}
		proxy.ServeHTTP(w, r)
	})
}
