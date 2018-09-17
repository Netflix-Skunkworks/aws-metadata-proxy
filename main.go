package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
)

var (
	whitelist [5]string
)

func init() {
	whitelist[0] = "aws-sdk-"
	whitelist[1] = "Botocore/"
	whitelist[2] = "Boto3/"
	whitelist[3] = "aws-cli/"
	whitelist[4] = "aws-chalice/"
}

func checkUserAgent(ua string) bool {
	for i := range whitelist {
		if strings.HasPrefix(ua, whitelist[i]) {
			return true
		}
	}
	return false
}

func ProxyHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ua := string(r.Header.Get("User-Agent"))

		if !checkUserAgent(ua) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized"))
			return
		}
		p.ServeHTTP(w, r)
	}
}

func main() {

	log.Println("Starting proxy...")

	target := "http://169.254.169.254"
	remote, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	router := mux.NewRouter()

	router.HandleFunc("/{path:.*}", ProxyHandler(proxy))

	go func() {
		log.Fatal(http.ListenAndServe("127.0.0.1:9090", router))
	}()

	// Check for interrupt signal and exit cleanly
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received, exiting proxy...")

}
