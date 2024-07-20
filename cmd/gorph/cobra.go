package main

import (
	"fmt"

	"github.com/sean9999/go-flargs"
	"github.com/spf13/cobra"
)

var optJson bool
var optBase64 bool
var optQuoted bool
var optDelimiter string

var env = flargs.NewCLIEnvironment("/")

var rootCmd = &cobra.Command{
	Use:   "gorph",
	Short: "gorph is a recursive file watcher",
	Long:  `Gorph is a recursuve file watcher that allows you to specify complex file matching patterns using doublestar '**' and other globbing techniques.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			panic("two arguments are required: root and pattern")
		}
		if len(args) > 2 {
			fmt.Fprintln(env.ErrorStream, "ignoring arguments beyond the 2nd")
		}
		flags := outputOptions{
			json:      optJson,
			base64:    optBase64,
			quoted:    optQuoted,
			delimiter: optDelimiter,
		}
		rootPath := args[0]
		pattern := args[1]
		watch(env, rootPath, pattern, flags)
	},
}

func init() {
	rootCmd.Flags().BoolVar(&optJson, "json", false, "format output as JSON")
	rootCmd.Flags().BoolVar(&optBase64, "base64", false, "encode output as base64")
	rootCmd.Flags().BoolVar(&optQuoted, "quoted", false, "wrap strings in single quotes")
	rootCmd.Flags().StringVarP(&optDelimiter, "delimiter", "d", "\t", "a character or sequence of characters used to delimit values in output")
}
