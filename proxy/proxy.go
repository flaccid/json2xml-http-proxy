package proxy

import (
	"bytes"
	"github.com/clbanning/mxj"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type transport struct {
	http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))

	// dump the egress response
	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Error(err)
	}
	log.Info(string(responseDump))

	return resp, nil
}

var _ http.RoundTripper = &transport{}
var xml []byte

func convJ2X(json []byte) []byte {
	m, err := mxj.NewMapJson(json)
	if err != nil {
		log.Error("error mapping json: ", err)
	}

	xml, err := m.Xml()
	if err != nil {
		log.Error("error converting xml: ", err)
	}

	return xml
}

// Serve : initialises a reverse proxy accepting JSON, forwarding as XML
func Serve(scheme string, host string, listenPort string) {
	// create a reverse proxy to the desired backend host
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: scheme,
		Host:   host,
	})

	// take control of the proxy transport and director
	proxy.Transport = &transport{http.DefaultTransport}
	realDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		// modify headers to taste
		req.Header.Set("X-Reverse-Proxy", "j2xrp")
		req.Header.Set("Content-Type", "application/xml")

		// request body JSON->XML happens here
		// we only care about body content if this is a PUT or POST
		if req.Method == "PUT" || req.Method == "POST" {
			buf := new(bytes.Buffer)
			if req.ContentLength > 0 {
				// stream the body and convert the expected JSON to XML
				buf.ReadFrom(req.Body)
				s := buf.String()
				xml = convJ2X([]byte(s))

			} else {
				xml = []byte(`<?xml version="1.0" ?>`)
			}
			req.Body = ioutil.NopCloser(strings.NewReader(string(xml)))
			req.ContentLength = int64(len(string(xml)))
		}

		// dump the ingress request
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Error(err)
		}
		log.Info(string(requestDump))

		realDirector(req)
	}

	http.Handle("/", proxy)
	log.Info("starting http server")
	log.Info("proxy requests to ", scheme, "://", host)
	log.Info("listening for requests on :" + listenPort)
	log.Fatal(http.ListenAndServe(":"+listenPort, proxy))
}
