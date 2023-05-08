/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/anaregdesign/lantern-cli/query"
	"github.com/manifoldco/promptui"
	"log"
	"os"

	"github.com/anaregdesign/lantern/client"
	"github.com/spf13/cobra"
)

var (
	host string
	port int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lantern-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		cli, err := client.NewLantern(host, port)
		exec := query.Executor{Client: cli}
		if err != nil {
			return err
		}
		defer func() {
			err := cli.Close()
			log.Fatal(err)
		}()

		for {
			prompt := promptui.Prompt{
				Label: ">",
			}

			result, err := prompt.Run()
			if err != nil {
				return err
			}

			switch result {
			case "exit":
				return nil

			default:
				s, err := exec.Execute(ctx, result)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if s != "" {
					fmt.Println(s)
				}
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&host, "host", "H", "localhost", "host")
	rootCmd.Flags().IntVarP(&port, "port", "p", 6380, "port")
}
