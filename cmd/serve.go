package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/virtualroot/tofulicious/internal/server"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"server"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")

		s, _ := server.NewHTTPServer()
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			if err := s.Start(); err != nil {
				log.Fatal("error starting server", "err", err)
			}
		}()

		<-done

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer func() { cancel() }()

		if err := s.Shutdown(ctx); err != nil {
			log.Fatalf("server forced to shutdown: %v", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
