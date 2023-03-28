package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("pong"))
		if err != nil {
			return
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		targetUrl := "https://api.openai.com"
		target, err := url.Parse(targetUrl)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		r.Header.Set("Authorization", r.Header.Get("Authorization"))
		r.Header.Set("Content-Type", "application/json")
		r.Host = r.URL.Host

		proxy := httputil.NewSingleHostReverseProxy(target)

		if r.Header.Get("Accept") == "text/event-stream" {
			handleSSE(w, r, proxy)
		} else {
			proxy.ServeHTTP(w, r)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func handleSSE(w http.ResponseWriter, r *http.Request, proxy *httputil.ReverseProxy) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	// Start proxying SSE request to OpenAI API
	done := make(chan bool)
	go func() {
		proxy.ServeHTTP(w, r)
		done <- true
	}()

	// Wait for SSE request to finish
	<-done
}
