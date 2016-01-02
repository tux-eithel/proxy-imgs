package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	flagOrigin      string
	flagRemote      string
	flagCheckStatus bool
	flagBindPort    int
	flagRegexp      Rgxs
)

func init() {

	flag.StringVar(&flagOrigin, "o", "", `Origin site: -o "http://site1.dev:80"`)
	flag.StringVar(&flagRemote, "r", "", `Remote site: -r "http://site2.dev:80"`)
	flag.BoolVar(&flagCheckStatus, "s", false, "if remote site gives error code > 400, request will be proxy to Origin. Default FALSE")
	flag.IntVar(&flagBindPort, "p", 80, "Port where bind the service: -p 8081 bind service to port")

	flag.Var(&flagRegexp, "f", `Patter to proxy to Remote site: -f ".jpg?" -f "wp-content/uploads/*"`)

}

func main() {

	var err error

	flag.Parse()
	if flagOrigin == "" || flagRemote == "" {
		log.Fatalln("insert at lease the Origin and the Remote site")
	}

	if flagBindPort < 1 {
		log.Fatalln("port must be a value > 1")
	}

	bs1, err := url.Parse(flagOrigin)
	if err != nil {
		log.Fatalln(err)
	}

	b2, err := url.Parse(flagRemote)
	if err != nil {
		log.Fatalln(err)
	}

	rp := prepareDoubleProxy(bs1, b2, []string(flagRegexp), flagCheckStatus)

	err = http.ListenAndServe(":"+strconv.Itoa(flagBindPort), rp)
	if err != nil {
		log.Fatalln(err)
	}

}

func prepareDoubleProxy(origin *url.URL, remote *url.URL, toMatch []string, checkStatus404 bool) *httputil.ReverseProxy {

	rg := make([]*regexp.Regexp, len(toMatch))

	var r *regexp.Regexp
	var err error

	i := 0
	for _, val := range toMatch {
		r, err = regexp.Compile(val)
		if err == nil {
			rg[i] = r
			i++
		} else {
			log.Printf("ignored pattern '%s'. why? %s", val, err)
		}
	}

	director := func(req *http.Request) {

		match := false
		for _, val := range rg {
			if val.Match([]byte(req.URL.Path)) {
				match = true
				break
			}
		}

		if !match {
			// didn't match, proxy the request to the Origin
			editRequest(origin, origin.RawQuery, req)
		} else {
			// match! so proxy to the remote server
			editRequest(remote, remote.RawQuery, req)

			// if checkStatus404, before proxy the request to Remote, check remote server response. If StatuCode > 400, proxy to Origin
			if checkStatus404 {
				reqTmp := new(http.Request)
				*reqTmp = *req
				reqTmp.RequestURI = ""
				reqTmp.URL = req.URL
				editRequest(remote, remote.RawQuery, reqTmp)
				resp, err := http.DefaultClient.Do(reqTmp)
				if err != nil {
					log.Printf("Error sub-request to check the status code: %s", err)
				}

				if err != nil || resp.StatusCode > 400 {
					editRequest(origin, origin.RawQuery, req)
				}
			}
		}

	}

	return &httputil.ReverseProxy{Director: director}
}

// editRequest is copy-paste of https://golang.org/src/net/http/httputil/reverseproxy.go#L69
func editRequest(target *url.URL, targetQuery string, req *http.Request) {
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}
}

// singleJoiningSlash is copy-paste of https://golang.org/src/net/http/httputil/reverseproxy.go#L51
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
