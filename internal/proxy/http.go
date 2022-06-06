package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type HttpProxyStrategy struct {
	listenAddress string
	targetBaseUrl string
}

func (h *HttpProxyStrategy) String() string {
	return fmt.Sprintf("HttpProxyStrategy{listenAddress=%s,targetBaseUrl=%s}", h.listenAddress, h.targetBaseUrl)
}

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func (p *HttpProxyStrategy) copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func (p *HttpProxyStrategy) dropHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func (p *HttpProxyStrategy) appendHostToXForwardHeader(header http.Header, host string) {
	// If we aren't the first Proxy retain prior
	// X-Forwarded-For information as a comma+space
	// separated list and fold multiple headers into one.
	if prior, ok := header["X-Forwarded-For"]; ok {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}

func (p *HttpProxyStrategy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	log.Println(req.RemoteAddr, " ", req.Method, " ", req.URL)

	targetUrl, err := url.Parse(fmt.Sprintf("%s%s", p.targetBaseUrl, req.URL.Path))
	if err != nil {
		msg := "internal error"
		http.Error(wr, msg, http.StatusBadRequest)
		log.Printf("creating url failed: %v", err)
		return
	}

	client := &http.Client{}

	//http: Request.RequestURI can't be set in client requests.
	//http://golang.org/src/pkg/net/http/client.go
	req.RequestURI = ""

	p.dropHopHeaders(req.Header)

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		p.appendHostToXForwardHeader(req.Header, clientIP)
	}

	req.URL = targetUrl
	resp, err := client.Do(req)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		log.Println("forwarding failed:", err)
		return
	}
	defer resp.Body.Close()

	log.Println(req.RemoteAddr, " ", resp.Status)

	p.dropHopHeaders(resp.Header)

	p.copyHeader(wr.Header(), resp.Header)
	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
}

func (h *HttpProxyStrategy) Start(ctx context.Context, eventChan chan interface{}) {
	// TODO ctx unused
	log.Printf("Listening on %s...", h.listenAddress)
	if err := http.ListenAndServe(h.listenAddress, h); err != nil {
		eventChan <- err
	}
}
