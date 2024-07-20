/*
Copyright © 2024 Sean Macdonald sean.r.macdonald@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var OptJson bool
var OptBase64 bool
var OptQuoted bool
var OptDelimiter string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gorph",
	Short: "gorph is a recursive file watcher",
	Long:  `Gorph is a recursuve file watcher that allows you to specify complex file matching pattersn using doublestar '**' and other globbing techniques.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		fmt.Println("json", OptJson)
		fmt.Println("base64", OptBase64)
		fmt.Println("quoted", OptQuoted)
		fmt.Println("delim", OptDelimiter)
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
	rootCmd.Flags().BoolVar(&OptJson, "json", false, "format output as JSON")
	rootCmd.Flags().BoolVar(&OptBase64, "base64", false, "encode output as base64")
	rootCmd.Flags().BoolVar(&OptQuoted, "quoted", false, "wrap strings in single quotes")
	rootCmd.Flags().StringVarP(&OptDelimiter, "delimiter", "d", ", ", "a character or sequence of characters used to delimit values in output")
	//rootCmd.MarkFlagsMutuallyExclusive("json", "quoted")
}
