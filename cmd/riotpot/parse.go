package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
	"github.com/riotpot/pkg/api"
	"github.com/riotpot/pkg/logger"
	"github.com/riotpot/pkg/plugins"
	"github.com/riotpot/pkg/proxy"
	"github.com/riotpot/ui"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	_ "github.com/riotpot/statik"
)

func setup(output string, pluginsPath string, services []string) {

	// Set the logger
	logger.Log = logger.New(zerolog.DebugLevel, output)

	// Load plugins
	px, err := plugins.LoadPlugins(pluginsPath)
	if err != nil {
		panic(err)
	}

	// Start services
	for _, p := range px {
		prx, err := proxy.Proxies.GetProxy(p.GetID())
		if err != nil {
			panic(err)
		}

		err = prx.Start()
		if err != nil {
			panic(err)
		}
	}
}

func createApiRouter(whitelist []string, startUi bool) *gin.Engine {
	router := gin.Default()
	router.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     whitelist,
				AllowMethods:     []string{"OPTIONS", "PUT", "PATCH", "GET", "DELETE"},
				AllowHeaders:     []string{"Content-Type", "Content-Length", "Origin"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
		),
	)

	group := router.Group("/api/")
	api.ProxiesRouter.AddToGroup(group)
	api.ServiceRouter.AddToGroup(group)

	if startUi {
		ui.AddRoutes(router)
	}

	statikFS, err := fs.New()
	if err != nil {
		logger.Log.Fatal().Err(err)
	}

	// Serve the Swagger UI files in the root of the api
	router.StaticFS("swagger", statikFS)
	return router
}

func parseSetupFlags(cmd *cobra.Command, args []string) {
	fgs := cmd.Flags()

	outFlag, err := fgs.GetString("output")
	if err != nil {
		panic(err)
	}

	srvFlag, err := fgs.GetStringArray("services")
	if err != nil {
		panic(err)
	}

	pluginsFlag, err := fgs.GetString("plugins")
	if err != nil {
		panic(err)
	}

	setup(outFlag, pluginsFlag, srvFlag)
}

func NewRootCommand() *cobra.Command {
	var cmds = &cobra.Command{
		Use: "riotpot",
		Run: func(cmd *cobra.Command, args []string) {
			parseSetupFlags(cmd, args)

			done := make(chan os.Signal, 1)
			signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
			fmt.Println("RIoTPot running, press ctrl+c to exit...")
			<-done
		},
	}

	rootFlags := cmds.Flags()
	rootFlags.StringArray("services", []string{}, "Starts a list of services")
	rootFlags.String("output", "", "Path to output file. E.g., 'path/to/riotpot.log'")
	rootFlags.String("plugins", "plugins/*.so", "Path to plugins folder")

	return cmds
}

func NewServerCommand() *cobra.Command {
	var cmdApi = &cobra.Command{
		Use:   "server",
		Short: "Starts RIoTPot as a server",
		Long:  "server starts RIoTPot as a server. It offers a REST API (and optionally a UI) to control the application while running",
		Run: func(cmd *cobra.Command, args []string) {
			parseSetupFlags(cmd, args)
			fgs := cmd.Flags()

			whitelistFlag, err := fgs.GetStringArray("whitelist")
			if err != nil {
				panic(err)
			}

			portFlag, err := fgs.GetInt("port")
			if err != nil {
				panic(err)
			}

			uiFlag, err := fgs.GetBool("with-ui")
			if err != nil {
				panic(err)
			}

			router := createApiRouter(whitelistFlag, uiFlag)
			addr := fmt.Sprintf(":%d", portFlag)
			err = router.Run(addr)
			if err != nil {
				panic(err)
			}
		},
	}

	apiFlags := cmdApi.Flags()
	apiFlags.StringArray("whitelist", []string{"http://localhost"}, "List of allowed hosts to interact with the API. Default: http://localhost")
	apiFlags.Int("port", 3000, "Server port port")
	apiFlags.Bool("with-ui", false, "Serve the UI as well")

	return cmdApi
}

func NewRiotpotCommand() *cobra.Command {

	cmds := NewRootCommand()
	rootFlags := cmds.Flags()

	cmdServer := NewServerCommand()
	serverFlags := cmdServer.Flags()
	serverFlags.AddFlagSet(rootFlags)

	cmds.AddCommand(cmdServer)
	return cmds
}
