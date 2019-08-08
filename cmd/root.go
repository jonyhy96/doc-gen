package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jonyhy96/doc-gen/generator"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var filename string

func init() {
	rootCmd.Flags().StringP("file", "f", "", "the relatepath of file base on your current path")
	viper.SetDefault("author", "jonyhy <github/jonyhy96>")
}

var rootCmd = &cobra.Command{
	Use:     "doc-gen",
	Short:   "Author jonyhy96 github.com/jonyhy96\ndoc-gen helps you genarate apidoc for specific file\n",
	Version: "1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		filename, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatalln(err.Error())
		}
		if filename == "" {
			log.Fatalln("must support filename")
		}
		currpath, _ := os.Getwd()
		generator.ParsePackagesFromDir(currpath)
		apidoc, err := generator.Scan(currpath, filename)
		if err != nil {
			log.Fatalln(err.Error())
		}
		generator.Gen(*apidoc, filepath.Join(currpath, filename)+".doc")
	},
}

// Execute command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
