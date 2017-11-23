package main

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/tls"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/streamlist/streamlist/internal/archiver"
	"github.com/streamlist/streamlist/internal/logtailer"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/acme/autocert"
)

var (
	cli = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// flags
	backlink               string
	datadir                string
	debug                  bool
	httpAddr               string
	httpAdmins             arrayFlags
	httpAdminUsers         []string
	httpReadOnlys          arrayFlags
	httpHost               string
	httpPrefix             string
	letsencrypt            bool
	reverseProxyAuthHeader string
	reverseProxyAuthIP     string

	// set based on httpAddr
	httpIP   string
	httpPort string

	// logging
	logger  *zap.SugaredLogger
	logtail *logtailer.Logtailer

	// archiver
	archive *archiver.Archiver

	// config
	config *Config

	// version
	version string
)

func init() {
	dbInit()
	cli.StringVar(&backlink, "backlink", "", "backlink (optional)")
	cli.StringVar(&datadir, "data-dir", "/data", "data directory")
	cli.BoolVar(&debug, "debug", false, "debug mode")
	cli.StringVar(&httpAddr, "http-addr", ":80", "listen address")
	cli.Var(&httpAdmins, "http-admin", "HTTP basic auth user/password for admins.")
	cli.Var(&httpReadOnlys, "http-read-only", "HTTP basic auth user/password for read only users.")
	cli.StringVar(&httpHost, "http-host", "", "HTTP host")
	cli.StringVar(&httpPrefix, "http-prefix", "/streamlist", "HTTP URL prefix (not actually supported yet!)")
	cli.BoolVar(&letsencrypt, "letsencrypt", false, "enable TLS using Let's Encrypt")
	cli.StringVar(&reverseProxyAuthHeader, "reverse-proxy-header", "X-Authenticated-User", "reverse proxy auth header")
	cli.StringVar(&reverseProxyAuthIP, "reverse-proxy-ip", "", "reverse proxy auth IP")
}

func main() {
	var err error

	cli.Parse(os.Args[1:])

	// Create users in db if not exists, or set password and role if needed
	for _, httpUser := range httpAdmins {
		split := strings.Split(httpUser, ":")
		httpUsername := split[0]
		httpUserPassword := split[1]
		hasher := sha512.New()
		hasher.Write([]byte(httpUserPassword))
		httpAdminUsers = append(httpAdminUsers, httpUsername)
		var user User
		db.Where(User{Username: httpUsername}).Assign(User{Password: hex.EncodeToString(hasher.Sum(nil)), Role: "admin"}).FirstOrCreate(&user)
	}
	for _, httpUser := range httpReadOnlys {
		split := strings.Split(httpUser, ":")
		httpUsername := split[0]
		httpUserPassword := split[1]
		hasher := sha512.New()
		hasher.Write([]byte(httpUserPassword))
		var user User
		db.Where(User{Username: httpUsername}).Assign(User{Password: hex.EncodeToString(hasher.Sum(nil)), Role: "readonly"}).FirstOrCreate(&user)
	}

	// logtailer
	logtail, err = logtailer.NewLogtailer(200 * 1024)
	if err != nil {
		panic(err)
	}

	// logger
	atomlevel := zap.NewAtomicLevel()
	l := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
			zapcore.NewMultiWriteSyncer(zapcore.Lock(zapcore.AddSync(os.Stdout)), logtail),
			atomlevel,
		),
	)
	defer l.Sync()
	logger = l.Sugar()

	// debug logging
	if debug {
		atomlevel.SetLevel(zap.DebugLevel)
	}
	logger.Debugf("debug logging is enabled")

	// config
	config, err = NewConfig("config.json")
	if err != nil {
		logger.Fatal(err)
	}

	// archiver
	archive = archiver.NewArchiver(datadir, 2, logger)

	// datadir
	datadir = filepath.Clean(datadir)
	if _, err := os.Stat(datadir); err != nil {
		logger.Debugf("creating datadir %q", datadir)
		if err := os.MkdirAll(datadir, 0755); err != nil {
			logger.Fatal(err)
		}
	}

	// remove any temporary transcode files
	tmpfiles, _ := filepath.Glob(datadir + "/*.transcoding")
	for _, tmpfile := range tmpfiles {
		logger.Debugf("removing %q", tmpfile)
		if err := os.Remove(tmpfile); err != nil {
			logger.Errorf("removing %q failed: %s", tmpfile, err)
		}
	}

	// usage
	usage := func(msg string) {
		fmt.Fprintf(os.Stderr, "ERROR: "+msg+"\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s --http-host music.example.com --http-admin 'admin:$ecUrePas$0rd'\n\n", os.Args[0])
		cli.PrintDefaults()
		os.Exit(1)
	}

	// http admin
	if httpAdmins == nil && reverseProxyAuthIP == "" {
		usage("the --http-admin or the --reverseProxyAuthIP flag is required")
	}

	// http host
	if httpHost == "" {
		usage("the --http-host flag is required")
	}
	httpPrefix = strings.TrimRight(httpPrefix, "/")

	// http port
	httpIP, httpPort, err := net.SplitHostPort(httpAddr)
	if err != nil {
		usage("invalid --http-addr")
	}

	//
	// Routes
	//
	r := httprouter.New()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	r.HandleMethodNotAllowed = false

	// Handlers
	r.GET("/", log(auth(index, "readonly")))
	r.GET(prefix("/logs"), log(auth(logs, "admin")))
	r.GET(prefix("/"), log(auth(home, "readonly")))

	// User
	//r.GET(prefix("/user/create"), log(auth(createUser, "none")))

	// Library
	r.GET(prefix("/library"), log(auth(library, "readonly")))

	// Media
	r.GET(prefix("/media/thumbnail/:media"), log(auth(thumbnailMedia, "readonly")))
	r.GET(prefix("/media/view/:media"), log(auth(viewMedia, "readonly")))
	r.GET(prefix("/media/delete/:media"), log(auth(deleteMedia, "admin")))
	r.GET(prefix("/media/access/:filename"), auth(streamMedia, "readonly"))
	r.GET(prefix("/media/download/:filename"), auth(downloadMedia, "readonly"))

	// Publicly accessible streaming (using playlist id as "auth")
	r.GET(prefix("/stream/:list/:filename"), auth(streamMedia, "none"))

	// Import
	r.GET(prefix("/import"), log(auth(importHandler, "admin")))

	// Archiver
	r.GET(prefix("/archiver/jobs"), auth(archiverJobs, "admin"))
	r.POST(prefix("/archiver/save/:id"), log(auth(archiverSave, "admin")))
	r.GET(prefix("/archiver/cancel/:id"), log(auth(archiverCancel, "admin")))

	// List
	r.GET(prefix("/create"), log(auth(createList, "admin")))
	r.POST(prefix("/create"), log(auth(createList, "admin")))
	r.POST(prefix("/add/:list/:media"), log(auth(addMediaList, "admin")))
	r.POST(prefix("/remove/:list/:media"), log(auth(removeMediaList, "admin")))
	r.GET(prefix("/remove/:list/:media"), log(auth(removeMediaList, "admin")))

	r.GET(prefix("/edit/:id"), log(auth(editList, "admin")))
	r.POST(prefix("/edit/:id"), log(auth(editList, "admin")))
	r.GET(prefix("/shuffle/:id"), log(auth(shuffleList, "admin")))
	r.GET(prefix("/play/:id"), log(auth(playList, "none")))
	r.GET(prefix("/m3u/:id"), log(auth(m3uList, "none")))
	r.GET(prefix("/podcast/:id"), log(auth(podcastList, "none")))

	r.POST(prefix("/config"), log(auth(configHandler, "admin")))

	r.GET(prefix("/delete/:id"), log(auth(deleteList, "admin")))

	// API
	r.GET(prefix("/v1/status"), log(auth(v1status, "none")))

	// Assets
	r.GET(prefix("/static/*path"), auth(staticAsset, "none"))
	r.GET(prefix("/logo.png"), log(auth(logo, "none")))

	//
	// Server
	//
	httpTimeout := 48 * time.Hour
	maxHeaderBytes := 10 * (1024 * 1024) // 10 MB

	// Plain text web server.
	if !letsencrypt {
		plain := &http.Server{
			Handler:        r,
			Addr:           httpAddr,
			WriteTimeout:   httpTimeout,
			ReadTimeout:    httpTimeout,
			MaxHeaderBytes: maxHeaderBytes,
		}

		hostport := net.JoinHostPort(httpHost, httpPort)
		if httpPort == "80" {
			hostport = httpHost
		}
		logger.Infof("Streamlist (version: %s) %s", version, &url.URL{
			Scheme: "http",
			Host:   hostport,
			Path:   httpPrefix + "/",
		})

		logger.Fatal(plain.ListenAndServe())
	}

	// Let's Encrypt TLS mode

	// http redirect to https
	go func() {
		redir := httprouter.New()
		redir.GET("/*path", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			r.URL.Scheme = "https"
			r.URL.Host = net.JoinHostPort(httpHost, httpPort)
			http.Redirect(w, r, r.URL.String(), http.StatusFound)
		})

		plain := &http.Server{
			Handler:        redir,
			Addr:           net.JoinHostPort(httpIP, "80"),
			WriteTimeout:   httpTimeout,
			ReadTimeout:    httpTimeout,
			MaxHeaderBytes: maxHeaderBytes,
		}
		if err := plain.ListenAndServe(); err != nil {
			logger.Warnf("skipping redirect http port 80 to https port %s (%s)", httpPort, err)
		}
	}()

	// autocert
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(filepath.Join(datadir, ".autocert")),
		HostPolicy: autocert.HostWhitelist(httpHost, "www."+httpHost),
	}

	// TLS
	tlsConfig := tls.Config{
		GetCertificate: m.GetCertificate,
		NextProtos:     []string{"http/1.1"},
		Rand:           rand.Reader,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,

			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	// Override default for TLS.
	if httpPort == "80" {
		httpPort = "443"
		httpAddr = net.JoinHostPort(httpIP, httpPort)
	}

	secure := &http.Server{
		Handler:        r,
		Addr:           httpAddr,
		WriteTimeout:   httpTimeout,
		ReadTimeout:    httpTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	// Enable TCP keep alives on the TLS connection.
	tcpListener, err := net.Listen("tcp", httpAddr)
	if err != nil {
		logger.Fatalf("listen failed: %s", err)
		return
	}
	tlsListener := tls.NewListener(tcpKeepAliveListener{tcpListener.(*net.TCPListener)}, &tlsConfig)

	hostport := net.JoinHostPort(httpHost, httpPort)
	if httpPort == "443" {
		hostport = httpHost
	}
	logger.Infof("Streamlist (version: %s) %s", version, &url.URL{
		Scheme: "https",
		Host:   hostport,
		Path:   httpPrefix + "/",
	})
	logger.Fatal(secure.Serve(tlsListener))
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (l tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := l.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(10 * time.Minute)
	return tc, nil
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
