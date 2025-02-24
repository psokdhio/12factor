/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"time"

	"github.com/psokdhio/12factor/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		h, err := server.NewPongHandler(server.PongHandlerOptions{
			BaseURL: "",
		})
		if err != nil {
			return err
		}
		err = server.Serve(ctx, h, server.ServeOptions{
			Addr:                viper.GetString("server.address"),
			ReadHeaderTimeout:   5 * time.Second,
			ShutdownGracePeriod: 30 * time.Second,
		})
		return err
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

	serveCmd.Flags().StringP("address", "a", "localhost:8080", "Serve at address")
	viper.BindPFlag("server.address", serveCmd.Flags().Lookup("address"))
}
