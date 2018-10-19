// Copyright (c) 2018 soren yang
//
// Licensed under the MIT License
// you may not use this file except in complicance with the License.
// You may obtain a copy of the License at
//
//     https://opensource.org/licenses/MIT
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package svs

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/lsytj0413/ena/logger"
	"github.com/lsytj0413/tyche/pkg/conf"
)

// Server is svs proj server
type Server interface {
	Start() (chan struct{}, error)
}

type server struct {
	c  *conf.Config
	fs *flag.FlagSet

	stop chan struct{}
}

var (
	usageline = ``
)

// New will construct a Server instance
func New() (Server, error) {
	s := &server{
		c:  conf.New(),
		fs: flag.NewFlagSet("svs", flag.ContinueOnError),
	}
	s.fs.Usage = func() {
		fmt.Fprintf(os.Stderr, usageline)
	}

	c := s.c
	s.fs.StringVar(&c.DefaultListenClientURL, "listen-client-url", "http://localhost:80", "URL to listen on for client traffic.")
	s.fs.StringVar(&c.Name, "name", c.Name, "Human-readable name for this member.")
	s.fs.StringVar(&c.ClientTLSInfo.CertFile, "client-cert-file", c.ClientTLSInfo.CertFile, "Path to the client server TLS cert file.")
	s.fs.StringVar(&c.ClientTLSInfo.KeyFile, "client-key-file", c.ClientTLSInfo.KeyFile, "Path to the client server TLS key file.")
	s.fs.StringVar(&c.ClientTLSInfo.TrustedCAFile, "client-trusted-ca-file", c.ClientTLSInfo.TrustedCAFile, "Path to the client server TLS trusted CA cert file.")
	s.fs.BoolVar(&c.ClientTLSInfo.ClientCertAuth, "client-cert-auth", false, "Enable client cert authentication.")
	s.fs.BoolVar(&c.ClientTLSInfo.InsecureSkipVerify, "client-auto-tls", false, "Client TLS using generated certificate.")
	s.fs.StringVar(&c.ClientTLSInfo.CRLFile, "client-crl-file", "", "Path to the client certificate revocation list file.")
	s.fs.BoolVar(&c.IsDebug, "debug", false, "enable debug log output")
	s.fs.BoolVar(&c.IsPprof, "pprof", false, "enable pprof")

	// wechat config
	s.fs.StringVar(&c.WxAppID, "wx-appid", "", "wechat appid")
	s.fs.StringVar(&c.WxAppSecret, "wx-appsecret", "", "wechat appsecret")
	s.fs.StringVar(&c.WxToken, "wx-token", "", "wechat token")
	s.fs.StringVar(&c.WxEncodingAESKey, "wx-aeskey", "", "wechat encoding aes key")

	s.stop = make(chan struct{}, 1)
	return s, nil
}

func parseArgs(s *server) error {
	err := s.fs.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	if len(s.fs.Args()) != 0 {
		return fmt.Errorf("'%s' is not a valid flag", s.fs.Arg(0))
	}

	return nil
}

func validateConfig(c *conf.Config) error {
	listenURL, err := url.Parse(c.DefaultListenClientURL)
	if err != nil {
		return fmt.Errorf("invalid listen-client-url: %s", err.Error())
	}

	if listenURL.Scheme == "https" {
		if c.ClientTLSInfo.CertFile == "" || c.ClientTLSInfo.KeyFile == "" {
			return fmt.Errorf("listen on https without keyfile or certfile exists")
		}
		c.IsTLSEnable = true
	}

	if c.ClientTLSInfo.ClientCertAuth {
		if c.ClientTLSInfo.TrustedCAFile == "" && !c.ClientTLSInfo.InsecureSkipVerify {
			return fmt.Errorf("client auth enable without client-trusted-ca-file or client-auto-tls set")
		}
	}

	return nil
}

func (s *server) Start() (chan struct{}, error) {
	if err := parseArgs(s); err != nil {
		return nil, err
	}

	if s.c.IsDebug {
		logger.SetLogLevel(logger.DebugLevel)
	}

	if err := validateConfig(s.c); err != nil {
		return nil, err
	}

	return s.start()
}

func (s *server) newTLSConfig() (*tls.Config, error) {
	tlsConfig := &tls.Config{}

	if s.c.ClientTLSInfo.ClientCertAuth {
		if s.c.ClientTLSInfo.InsecureSkipVerify {
			tlsConfig.InsecureSkipVerify = true
			tlsConfig.ClientAuth = tls.RequireAnyClientCert
		} else {
			pool := x509.NewCertPool()
			caCrt, err := ioutil.ReadFile(s.c.ClientTLSInfo.TrustedCAFile)
			if err != nil {
				return nil, fmt.Errorf("client auth cafile read error: %s", err.Error())
			}

			pool.AppendCertsFromPEM(caCrt)
			tlsConfig.ClientCAs = pool
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
	}

	if s.c.ClientTLSInfo.CRLFile != "" {
		logger.Infof("ignore client crlfile: %s", s.c.ClientTLSInfo.CRLFile)
	}

	return tlsConfig, nil
}

func (s *server) start() (chan struct{}, error) {
	tlsConfig, err := s.newTLSConfig()
	if err != nil {
		return nil, err
	}

	listenURL, _ := url.Parse(s.c.DefaultListenClientURL)
	srv := &http.Server{
		Addr:      listenURL.Hostname() + ":" + listenURL.Port(),
		TLSConfig: tlsConfig,
	}

	r := gin.New()
	r.Use(logMiddleware())
	r.GET("/version", wrapperHandler(s.Version))
	r.GET("/", wrapperHandler(s.Index))
	r.GET("/api/wx/mainEntry", s.WxVerify)
	r.POST("/api/wx/mainEntry", s.WxEntry)

	if s.c.IsPprof {
		pprof.Register(r, nil)
	}

	srv.Handler = r

	ch := make(chan error, 1)
	go func() {
		var err error
		if s.c.IsTLSEnable {
			logger.Infof("Listening and serving HTTPS on %s", srv.Addr)
			err = srv.ListenAndServeTLS(s.c.ClientTLSInfo.CertFile, s.c.ClientTLSInfo.KeyFile)
		} else {
			logger.Infof("Listening and serving HTTP on %s", srv.Addr)
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			logger.Errorf("Server Start Failed: %s", err)
			ch <- err
		}
	}()

	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)

		var err error
		select {
		case <-quit:
		case err = <-ch:
		}

		// Close on signal
		if err == nil {
			logger.Infof("Shutdown Server ...")

			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				logger.Errorf("Server Shutdown: ", err)
			}
		}

		s.stop <- struct{}{}
	}()

	return s.stop, nil
}
