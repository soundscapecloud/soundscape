package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/julienschmidt/httprouter"
)

/*
	"hasmedia": func(lid, mid string) bool {
		list, err := FindList(lid)
		if err != nil {
			return false
		}
		media, err := FindMedia(mid)
		if err != nil {
			return false
		}
		return list.HasMedia(media)
	},
*/

var (
	funcMap = template.FuncMap{
		"sub": func(a, b int64) int64 {
			return a - b
		},
		"add": func(a, b int64) int64 {
			return a + b
		},
		"nums": func(max int) (nums []int) {
			for i := 0; i < max; i++ {
				nums = append(nums, i)
			}
			return nums
		},
		"mediaexists": func(id string) bool {
			_, err := FindMedia(id)
			return err == nil
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"bytes": func(n int64) string {
			return humanize.Bytes(uint64(n))
		},
		"time": humanize.Time,
		"duration": func(seconds int64) string {
			hours := seconds / 3600
			seconds -= hours * 3600

			minutes := seconds / 60
			seconds -= minutes * 60

			if hours > 0 {
				return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
			}
			return fmt.Sprintf("%d:%02d", minutes, seconds)
		},
	}
	errorPageHTML = `
        <html>
            <head>
                <title>Error</title>
            </head>
            <body>
                <h2 style="color: orangered;">An error has occurred. <a href="/streamlist/logs">Check the logs</a></h2>
            </body>
        </html>
    `
)

func Redirect(w http.ResponseWriter, r *http.Request, format string, a ...interface{}) {
	location := httpPrefix
	location += fmt.Sprintf(format, a...)
	http.Redirect(w, r, location, http.StatusFound)
}

func Error(w http.ResponseWriter, err error) {
	logger.Error(err)

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, errorPageHTML)
}

func Prefix(path string) string {
	return httpPrefix + path
}

func Log(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Request info
		addr := r.RemoteAddr
		xff := r.Header.Get("X-Forwarded-For")
		realip := r.Header.Get("X-Real-IP")
		method := r.Method
		rang := r.Header.Get("Range")
		path := r.RequestURI

		// Run the handler
		start := time.Now()
		h(w, r, ps)
		elapsed := int64(time.Since(start) / time.Millisecond)

		// Response info
		mime := w.Header().Get("Content-Type")
		logger.Infof("%q %q %q %q %q %q %q %d ms", addr, xff, realip, method, path, rang, mime, elapsed)
	}
}

func auth(h httprouter.Handle, role string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := ""

		// none role is no auth required
		if role == "none" {
			h(w, r, ps)
			return
		}

		// Method: Basic Auth (if we're not behind a reverse proxy, use basic auth)
		if httpAdmins != nil {
			var userList []string
			// Admin are always OK
			userList = append(userList, httpAdmins...)
			// If role readonly, we add readonly users
			if role == "readonly" {
				userList = append(userList, httpReadOnlys...)
			}
			for _, httpUser := range userList {
				split := strings.Split(httpUser, ":")
				httpUsername := split[0]
				httpPassword := split[1]
				user, password, _ := r.BasicAuth()
				if user == httpUsername && password == httpPassword {
					ps = append(ps, httprouter.Param{Key: "user", Value: user})
					h(w, r, ps)
					return
				}
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="Sign-in Required"`)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Method: Reverse Proxy (if we're behind a reverse proxy, trust it.)

		clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if clientIP == reverseProxyAuthIP {
			user = r.Header.Get(reverseProxyAuthHeader)
		}

		//if user == "" && !optional {
		if user == "" {
			logger.Errorf("auth failed: client %q", clientIP)
			if backlink != "" {
				http.Redirect(w, r, backlink, http.StatusFound)
				return
			}
			http.NotFound(w, r)
			return
		}

		// Add "user" to params.
		if user != "" {
			ps = append(ps, httprouter.Param{Key: "user", Value: user})
		}
		h(w, r, ps)
	}
}

func XML(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	enc := xml.NewEncoder(w)
	enc.Indent("", "    ")
	if err := enc.Encode(data); err != nil {
		logger.Error(err)
	}
}

func JSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(data); err != nil {
		logger.Error(err)
	}
}

func HTML(w http.ResponseWriter, target string, data interface{}) {
	t := template.New(target)
	t.Funcs(funcMap)
	for _, filename := range AssetNames() {
		if !strings.HasPrefix(filename, "templates/") {
			continue
		}
		name := strings.TrimPrefix(filename, "templates/")
		b, err := Asset(filename)
		if err != nil {
			Error(w, err)
			return
		}

		var tmpl *template.Template
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		if _, err := tmpl.Parse(string(b)); err != nil {
			Error(w, err)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		Error(w, err)
		return
	}
}
