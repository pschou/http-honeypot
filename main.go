package main

import (
	"fmt"
	"github.com/Maldris/mathparse"
	"github.com/pschou/go-params"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func honeypot(w http.ResponseWriter, req *http.Request) {
	ua, ua_ok := req.Header["User-Agent"]
	if ua_ok && len(ua) > 0 {
		str := ua[0]
		bash_math := strings.Index(str, "$((")
		if bash_math >= 0 {
			tmp := str[bash_math+3:]
			en := strings.Index(tmp, "))")
			tr := mathparse.NewParser(tmp[:en])
			tr.Resolve()
			fmt.Printf("%s %v %s", str[:bash_math], tr.GetValueResult(), tmp[2+en:])
			str = fmt.Sprintf("%s%v%s", str[:bash_math], tr.GetValueResult(), tmp[2+en:])
		}

		//m1 := regexp.MustCompile(`\$\(\[0-9+1 ]*\)\)`)
		//str := m1.ReplaceAllString(ua[0], "%$1")

		m1 := regexp.MustCompile(`\\x([0-9]*)`)
		str = m1.ReplaceAllString(str, "%$1")
		str, _ = url.QueryUnescape(str)
		echo := strings.Index(str, " echo ")
		if echo >= 0 {
			shell_cmd := strings.Split(str[echo:], ";")
			for _, s := range shell_cmd {
				fmt.Println("index of echo ", s)
				s = strings.TrimSpace(s)
				if strings.HasPrefix(s, "echo Content-Type:") {
				} else if strings.HasPrefix(s, "/usr/bin/id") {
					fmt.Fprintln(w, "uid=0(root) uid=0(root) groups=0(root)")
				} else if strings.HasPrefix(s, "echo ") {
					fmt.Fprintf(w, "%s\n", strings.TrimPrefix(s, "echo "))
				}
			}
			fmt.Fprintf(w, "Anything else?\n\n")
			//return
		}
	}

	str, _ := url.QueryUnescape(req.URL.String())
	if strings.HasPrefix(str, "/DB4Web/") {
		parts := strings.Split(str, "/")
		_, port, _ := net.SplitHostPort(parts[2])
		if port != "" {
			fmt.Fprintf(w, "oooh, do you have something you want to show me at %s?  I'll go take a look.  Nice!\n\n", parts[2])
			conn, err := net.Dial("tcp", parts[2])
			if err == nil {
				defer conn.Close()
				fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
			}
		}
	}

	fmt.Fprintf(w, "%s\n\nRequest: %+v\nEscaped: %s\n", title, req.URL, str)
	if show_headers {
		for name, headers := range req.Header {
			for _, h := range headers {
				fmt.Fprintf(w, "%v: %v\n", name, h)
			}
		}
	}
}

var title = "This is my server."
var show_headers = true
var version = ""
var listen = ":8090"

func init() {
	params.CommandLine.Title = "jqURL - URL and JSON parser tool, Written by Paul Schou (github.com/pschou/jqURL), Version: " + version
	params.StringVar(&title, "header", title, "What header to print for blank requests", "STRING")
	params.StringVar(&listen, "listen", listen, "Listen address", "PORT")
	params.PresVar(&show_headers, "show_header", "Should we show headers?")
	params.Parse()
}

func main() {
	http.HandleFunc("/", honeypot)

	fmt.Println("Version", version, "Listen", listen)
	http.ListenAndServe(listen, nil)
}
