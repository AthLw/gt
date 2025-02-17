package web

import (
	"context"
	"crypto/tls"
	"embed"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/isrc-cas/gt/logger"
	"github.com/isrc-cas/gt/predef"
	"github.com/isrc-cas/gt/server"
	"github.com/isrc-cas/gt/server/web/api"
	"github.com/isrc-cas/gt/util"
	wServer "github.com/isrc-cas/gt/web/server"
	"github.com/isrc-cas/gt/web/server/middleware"
	"github.com/isrc-cas/gt/web/server/model/request"
	webUtil "github.com/isrc-cas/gt/web/server/util"
	"github.com/libp2p/go-reuseport"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var FS embed.FS

type Server struct {
	server       *http.Server
	logger       logger.Logger // have no right to close logger
	tokenManager *wServer.TokenManager
	enableTLS    bool
}

func NewWebServer(s *server.Server) (*Server, error) {

	//make sure the web Addr is valid
	err := checkConfig(s)
	if err != nil {
		return nil, err
	}

	tm := wServer.NewTokenManager(wServer.DefaultTokenManagerConfig())

	//set the log file
	fullPath := filepath.Join(util.GetAppDir(), "web_server.log")
	f, _ := os.Create(fullPath)
	gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()
	err = setRoutes(s, tm, r)
	if err != nil {
		return nil, err
	}

	ws, err := getServer(s, tm, r)
	if err != nil {
		return nil, err
	}

	go ws.start(func() {
		validURL, _ := webUtil.SwitchToValidWebURL(ws.server.Addr, ws.enableTLS)
		if isFirstStart(s) {
			tempKey, err := getTempKey(ws, s)
			if err == nil {
				//add tempKey to url
				validURL = webUtil.CreateUrlWithTempKey(validURL, tempKey)
				s.Logger.Info().Str("url", validURL).Msg("first start, browser url")
				s.Logger.Info().Msg("You have 3 chances to use this within the next 3 minutes. Please use it as soon as possible.")
			} else {
				s.Logger.Info().Err(err).Msg("failed to CreateUrlWithTempKey")
			}
		}
		if err := webUtil.OpenBrowser(validURL); err != nil {
			s.Logger.Info().Err(err).Msg("failed to open browser, please open it manually")
		}
	})
	return ws, err
}

// checkConfig checks the config of webAddr and set signingKey if it is not set
func checkConfig(s *server.Server) error {
	if len(s.Config().WebAddr) == 0 {
		return errors.New("option web_addr must be set")
	}
	if strings.IndexByte(s.Config().WebAddr, ':') == -1 {
		s.Config().WebAddr = ":" + s.Config().WebAddr
	}
	if len(s.Config().Config) == 0 {
		s.Config().Config = util.GetDefaultServerConfigPath()
		s.Logger.Info().Str("config", s.Config().Config).Msg("use default config path")
	}
	if len(s.Config().SigningKey) == 0 {
		s.Config().SigningKey = util.RandomString(predef.DefaultSigningKeySize)
	}
	if len(s.Config().WebCertFile) > 0 && len(s.Config().WebKeyFile) > 0 {
		if s.Config().WebCertFile != "auto" && s.Config().WebKeyFile != "auto" {
			_, certErr := os.Stat(s.Config().WebCertFile)
			_, keyErr := os.Stat(s.Config().WebKeyFile)
			if os.IsNotExist(certErr) || os.IsNotExist(keyErr) {
				return errors.New("provided webCertFile or webKeyFile does not exist")
			}
		}
	}
	return nil
}
func getServer(s *server.Server, tokenManager *wServer.TokenManager, r *gin.Engine) (*Server, error) {
	webServer := &Server{
		server: &http.Server{
			Addr:    s.Config().WebAddr,
			Handler: r,
		},
		logger:       s.Logger,
		tokenManager: tokenManager,
		enableTLS:    false,
	}
	certFile := s.Config().WebCertFile
	keyFile := s.Config().WebKeyFile

	if len(certFile) == 0 && len(keyFile) == 0 {
		// if both certFile and keyFile are empty, use HTTP, no need to do anything else
		return webServer, nil
	} else if certFile == "auto" && keyFile == "auto" {
		// if both certFile and keyFile are "auto", generate a self-signed certificate and enable HTTPS
		webServer.server.TLSConfig = webUtil.GenerateCertification()
		webServer.enableTLS = true
	} else {
		// if user provides certFile and keyFile, load them and enable HTTPS
		tlsCert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			err = errors.New("failed to load certFile or keyFile for web server")
			return nil, err
		}
		webServer.server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
		}
		webServer.enableTLS = true
	}
	return webServer, nil
}

func setRoutes(s *server.Server, tm *wServer.TokenManager, r *gin.Engine) error {
	staticFiles, err := fs.Sub(FS, "dist")
	if err != nil {
		return err
	}
	r.StaticFS("/static", http.FS(staticFiles))

	PublicGroup := r.Group("/")
	{
		PublicGroup.POST("/api/login", api.Login(s))
		PublicGroup.GET("/api/health", api.HealthCheck)
		PublicGroup.GET("/api/verify", api.VerifyTempKey(tm))
	}
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.JWTAuthMiddleware(s.Config().SigningKey, predef.DefaultTokenDuration))
	{
		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/change", api.ChangeUserInfo(s))
			userGroup.GET("/info", api.GetUserInfo(s))
		}
		configGroup := apiGroup.Group("/config")
		{
			configGroup.GET("/running", api.GetRunningConfig(s))
			configGroup.GET("/file", api.GetConfigFromFile(s))
			configGroup.POST("/save", api.SaveConfigToFile(s))
		}

		serverGroup := apiGroup.Group("/server")
		{
			serverGroup.GET("/info", api.GetServerInfo)
			serverGroup.PUT("/restart", api.Restart)
			serverGroup.PUT("/stop", api.Stop)
			serverGroup.PUT("/kill", api.Kill)
		}

		connectionGroup := apiGroup.Group("/connection")
		{
			connectionGroup.GET("/list", api.GetConnectionInfo(s))
		}

		permissionGroup := apiGroup.Group("/permission")
		{
			permissionGroup.GET("/menu", api.GetMenu(s))
		}
	}

	if s.Config().EnablePprof {
		pprofGroup := r.Group("/debug/pprof")
		{
			pprofGroup.GET("/", gin.WrapF(pprof.Index))
			pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
			pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
			pprofGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
			pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
			pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
			pprofGroup.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
			pprofGroup.GET("/block", gin.WrapH(pprof.Handler("block")))
			pprofGroup.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
			pprofGroup.GET("/heap", gin.WrapH(pprof.Handler("heap")))
			pprofGroup.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
			pprofGroup.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
		}

	}

	r.NoRoute(func(ctx *gin.Context) {
		data, err := fs.ReadFile(FS, "dist/index.html")
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	return nil
}

func (s *Server) start(readyCallback func()) {
	var err error
	defer s.logger.Info().Err(err).Msg("web server stopped")
	var ln net.Listener
	if util.IsNoArgs() {
		addr := s.server.Addr
		if addr == "" {
			addr = ":http"
		}
		for {
			s.logger.Info().Str("addr", addr).Msg("web server started")
			ln, err = reuseport.Listen("tcp", addr)
			if err == nil {
				break
			}
			var opErr *net.OpError
			if errors.As(err, &opErr) {
				var sysErr *os.SyscallError
				if errors.As(opErr.Err, &sysErr) {
					if errors.Is(sysErr.Err, syscall.EADDRINUSE) {
						s.logger.Warn().Err(err).Msg("web server address already in use, retrying...")
						var host string
						host, _, err = net.SplitHostPort(addr)
						if err != nil {
							return
						}
						addr = host + ":0"
						continue
					}
				}
			}
			return
		}
		if readyCallback != nil {
			readyCallback()
		}
		err = s.server.Serve(ln)
		return
	} else {
		s.logger.Info().Str("addr", s.server.Addr).Msg("web server started")
		if s.enableTLS {
			s.logger.Info().Msg("start web server with TLS")
			addr := s.server.Addr
			if addr == "" {
				addr = ":https"
			}

			ln, err = reuseport.Listen("tcp", addr)
			if err != nil {
				return
			}

			defer func() {
				_ = ln.Close()
			}()
			if readyCallback != nil {
				readyCallback()
			}
			err = s.server.ServeTLS(ln, "", "")
			return
		} else {
			s.logger.Info().Msg("start web server without TLS")
			addr := s.server.Addr
			if addr == "" {
				addr = ":http"
			}
			ln, err = reuseport.Listen("tcp", addr)
			if err != nil {
				return
			}
			if readyCallback != nil {
				readyCallback()
			}
			err = s.server.Serve(ln)
			return
		}
	}
}

// isFirstStart checks if the server is started for the first time, and set a random user
func isFirstStart(s *server.Server) bool {
	var isFirst bool
	// if the Admin and Password is not set,
	// treat it as the first start
	if len(s.Config().Admin) == 0 {
		s.Config().Admin = util.RandomString(predef.DefaultAdminSize)
		isFirst = true
	}
	if len(s.Config().Password) == 0 {
		s.Config().Password = util.RandomString(predef.DefaultPasswordSize)
		isFirst = true
	}
	return isFirst
}

func getTempKey(ws *Server, s *server.Server) (string, error) {
	tempUser := request.User{
		Username: s.Config().Admin,
		Password: s.Config().Password,
	}
	token, err := webUtil.GenerateToken(s.Config().SigningKey, predef.DefaultTokenDuration, "gt-server", tempUser)
	if err != nil {
		return "", err
	}
	tempKey := ws.tokenManager.GenerateTempKey(token)

	return tempKey, nil
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
