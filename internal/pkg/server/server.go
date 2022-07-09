package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/go-chi/chi"
	chimid "github.com/go-chi/chi/middleware"

	"github.com/serjyuriev/shortener/internal/app/grpcsvc"
	"github.com/serjyuriev/shortener/internal/pkg/config"
	"github.com/serjyuriev/shortener/internal/pkg/handlers"
	"github.com/serjyuriev/shortener/internal/pkg/middleware"
	"github.com/serjyuriev/shortener/proto/grpchandlers"
	"google.golang.org/grpc/keepalive"
)

// Server provides method for application server management.
type Server interface {
	Start() error
}

type server struct {
	cfg         *config.Config
	handlers    *handlers.Handlers
	grpcservice *grpcsvc.Service
}

// NewServer initializes server.
func NewServer() (Server, error) {
	h, err := handlers.MakeHandlers()
	if err != nil {
		return nil, fmt.Errorf("unable to make handlers:\n%w", err)
	}

	cfg := config.GetConfig()
	if cfg.EnableHTTPS {
		if err = createCerfs(); err != nil {
			return nil, fmt.Errorf("unable to create certificate: %v", err)
		}
	}

	gsvc, err := grpcsvc.MakeService()
	if err != nil {
		return nil, fmt.Errorf("unable to make handlers:\n%w", err)
	}

	return &server{
		cfg:         cfg,
		handlers:    h,
		grpcservice: gsvc,
	}, nil
}

// Start creates new router, binds handlers and starts http server.
func (s *server) Start() error {
	go func() {
		listen, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Println("could not start grpc server")
			return
		}
		serv := grpc.NewServer(
			grpc.KeepaliveParams(
				keepalive.ServerParameters{
					MaxConnectionIdle: 5 * time.Minute,
				},
			),
		)
		grpchandlers.RegisterShortenerServer(
			serv,
			s.grpcservice,
		)

		log.Println("Сервер gRPC начал работу")
		if err := serv.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()

	r := chi.NewRouter()
	r.Use(chimid.Recoverer)
	r.Use(chimid.Compress(gzip.BestSpeed, zippableTypes...))
	r.Use(middleware.Gzipper)
	r.Use(middleware.Auth)
	r.Delete("/api/user/urls", s.handlers.DeleteURLsHandler)
	r.Get("/ping", s.handlers.PingHandler)
	r.Get("/{shortPath}", s.handlers.GetURLHandler)
	r.Get("/api/user/urls", s.handlers.GetUserURLsAPIHandler)
	r.Get("/api/internal/stats", s.handlers.GetStatsHandler)
	r.Post("/", s.handlers.PostURLHandler)
	r.Post("/api/shorten", s.handlers.PostURLApiHandler)
	r.Post("/api/shorten/batch", s.handlers.PostBatchHandler)

	server := &http.Server{
		Addr:    s.cfg.ServerAddress,
		Handler: r,
	}

	sigChan := make(chan os.Signal, 3)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigChan
		log.Printf("\r\nПолучен сигнал: %s", sig.String())

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	go func() {
		log.Println(http.ListenAndServe(":8082", nil))
	}()

	log.Printf("starting server on %s\n", s.cfg.ServerAddress)
	if s.cfg.EnableHTTPS {
		return server.ListenAndServeTLS("cert.pem", "key.pem")
	} else {
		return server.ListenAndServe()
	}
}

var zippableTypes = []string{
	"application/javascript",
	"application/json",
	"text/css",
	"text/html",
	"text/plain",
	"text/xml",
}

func createCerfs() error {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(164),
		Subject: pkix.Name{
			Organization: []string{"SVY"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certFile, err := os.OpenFile("cert.pem", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("unable to open cert file: %v", err)
		return err
	}
	defer certFile.Close()
	if _, err = certPEM.WriteTo(certFile); err != nil {
		log.Printf("unable to write data to cert file: %v", err)
		return err
	}

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	keyFile, err := os.OpenFile("key.pem", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("unable to open key file: %v", err)
		return err
	}
	defer keyFile.Close()
	if _, err = privateKeyPEM.WriteTo(keyFile); err != nil {
		log.Printf("unable to write data to key file: %v", err)
		return err
	}

	return nil
}
