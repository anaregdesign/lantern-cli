/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/anaregdesign/lantern-cli/service"
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
		srv := service.NewCLIService(cli)
		if err != nil {
			return err
		}
		defer func() {
			err := cli.Close()
			if err != nil {
				log.Fatal(err)
			}
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
				err := srv.Run(ctx, result)
				switch err {
				case nil:
				case service.ErrGetVertex:
					fmt.Println("Usage: get vertex <key: string>")
				case service.ErrGetEdge:
					fmt.Println("Usage: get edge <tail: string> <head: string>")

				case service.ErrPutVertex:
					fmt.Println("Usage: put vertex <key: string> <value: string> [<ttl_seconds: int>]")

				case service.ErrPutEdge:
					fmt.Println("Usage: put edge <tail: string> <head: string> <weight: float> [<ttl_seconds: int>]")

				case service.ErrDeleteVertex:
					fmt.Println("Usage: delete vertex <key: string>")

				case service.ErrDeleteEdge:
					fmt.Println("Usage: delete edge <tail: string> <head: string>")

				case service.ErrAddEdge:
					fmt.Println("Usage: add edge <tail: string> <head: string> <weight: float> [<ttl_seconds: int>]")

				case service.ErrIlluminate:
					fmt.Println("Usage: illuminate { neighbor | spt_relevance | spt_cost | msp_relevance | msp_cost } <seed: string> <step: int> <k: int> <tfidf: bool>")

				case service.ErrInvalidVerb:
					fmt.Println("Usage: { get | put | add | illuminate } ...")

				case service.ErrInvalidObjective:
					fmt.Println("{\n\tget { vertex | edge } | \n\tput { vertex | edge } | \n\tadd edge | \n\tilluminate { neighbor | spt_relevance | spt_cost | msp_relevance | msp_cost}\n} ...\nspt: shortest path tree\nmsp: minimum spanning tree")

				case service.ErrConnection:
					fmt.Println("server error")

				default:
					fmt.Println(err)
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
