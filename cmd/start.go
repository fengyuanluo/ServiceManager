package cmd

import (
	"fmt"
	"net/http"

	"github.com/fengyuanluo/go-backend-template/docs"
	"github.com/fengyuanluo/go-backend-template/internal/conf"
	"github.com/fengyuanluo/go-backend-template/internal/db"
	"github.com/fengyuanluo/go-backend-template/internal/op"
	_ "github.com/fengyuanluo/go-backend-template/internal/server/handlers"
	"github.com/fengyuanluo/go-backend-template/internal/server/assets"
	"github.com/fengyuanluo/go-backend-template/internal/server/middleware"
	"github.com/fengyuanluo/go-backend-template/internal/server/router"
	"github.com/fengyuanluo/go-backend-template/internal/utils/log"
	"github.com/fengyuanluo/go-backend-template/internal/utils/shutdown"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var cfgFile string
var registerRoutes = router.RegisterAll
var initDB = db.InitDB
var closeDB = db.Close
var initUser = op.UserInit

func initializeApp() error {
	if err := initDB(); err != nil {
		return err
	}
	if err := initUser(); err != nil {
		return err
	}
	return nil
}

func setupRouter() (*gin.Engine, error) {
	r := gin.New()
	r.Use(middleware.Cors())
	r.Use(middleware.Logger())
	r.Use(middleware.StaticEmbed("/", assets.StaticFS))
	if err := registerRoutes(r); err != nil {
		return nil, err
	}
	return r, nil
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start " + conf.APP_NAME,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		conf.PrintBanner()
		if err := conf.Load(cfgFile); err != nil {
			return err
		}
		log.SetLevel(conf.AppConfig.Log.Level)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		sd := shutdown.New(log.Logger)

		if conf.IsDebug() {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}

		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", conf.AppConfig.Server.Host, conf.AppConfig.Server.Port)
		docs.SwaggerInfo.BasePath = "/api"

		if err := initializeApp(); err != nil {
			log.Errorf("app init error: %v", err)
			return
		}
		sd.Register(closeDB)

		r, err := setupRouter()
		if err != nil {
			log.Errorf("router init error: %v", err)
			return
		}

		httpSrv := &http.Server{Addr: fmt.Sprintf("%s:%d", conf.AppConfig.Server.Host, conf.AppConfig.Server.Port), Handler: r}
		go func() {
			if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Errorf("http server listen and serve error: %v", err)
			}
		}()
		sd.Register(httpSrv.Close)
		sd.Listen()
	},
}

func init() {
	startCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./data/config.json)")
	rootCmd.AddCommand(startCmd)
}
