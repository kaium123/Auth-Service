package command

import (
	"auth/common/logger"
	"auth/route"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	var serverPort string
	defaultServerPort := viper.GetString("SERVER_PORT")
	serverCmd.PersistentFlags().StringVar(&serverPort, "port", defaultServerPort, "Server port")
	viper.BindPFlag("SERVER_PORT", serverCmd.PersistentFlags().Lookup("port"))

	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server",
	Run: func(cmd *cobra.Command, args []string) {
		router := route.Setup()
		port := viper.GetString("SERVER_PORT")

		logger.LogInfo("Running server on port: " + port)

		server := &http.Server{
			Addr:    ":" + port,
			Handler: router,
		}

		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.LogInfo("Shutting down server...")

		//shutdown gin
		serverCloseContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		fmt.Println("getting ready for auth server shutdown")
		if err := server.Shutdown(serverCloseContext); err != nil {
			log.Fatal("Server Shutdown err:", err)
		}

		time.Sleep(time.Millisecond * 100)
		fmt.Println("bye bye")
	},
}
