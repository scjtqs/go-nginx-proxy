package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type Pxy struct{}

func (p *Pxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)
	transport := http.DefaultTransport

	var scheme=req.URL.Scheme
	// step 1
	outReq := new(http.Request)
	*outReq = *req // this only does shallow copies of maps
	client_host := outReq.Host
	outReq.Host = "nginx.scjtqs.com"
	outReq.URL.Host="nginx.scjtqs.com"
	outReq.URL.Scheme="https"
	//client_port:=outReq.URL.Port()
	outReq.URL.Port()
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}
	cookie := "host=" + client_host + ";scheme="+scheme
	outReq.Header.Set("Cookie", cookie)

	// step 2
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	// step 3
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}

	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
	res.Body.Close()
}

func main() {
	fmt.Println("Serve on :8889")
	http.Handle("/", &Pxy{})
	http.ListenAndServe("127.0.0.1:8889", nil)
}
