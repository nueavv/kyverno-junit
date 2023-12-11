/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/nueavv/kyverno-junit/utils/converter"
	"github.com/spf13/cobra"
)

var (
	filename        string
	output          string
	isclusterpolicy bool
)

const (
	cliName = "kyverno-junit"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     cliName,
	Args:    cobra.MaximumNArgs(1),
	Example: `hello`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed read file: %v", err)
		}

		switch isclusterpolicy {
		case true:
			report, err := converter.readClusterPolicyReport(data)
			if err != nil {
				return fmt.Errorf("failed cluster policy report file: %v", err)
			}
			err = converter.MakeClusterJunitReport(report, output)
			if err != nil {
				return fmt.Errorf("failed make report file: %v", err)
			}
		default:
			report, err := converter.readPolicyReport(data)
			if err != nil {
				return fmt.Errorf("failed policy report file: %v", err)
			}
			err = converter.MakeJunitReport(report, output)
			if err != nil {
				return fmt.Errorf("failed make report file: %v", err)
			}
		}
		fmt.Println("Success")
		return err
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
	rootCmd.Flags().BoolVarP(&isclusterpolicy, "cluster", "c", false, "Kyverno cluster policy")
	rootCmd.Flags().StringVarP(&filename, "filename", "f", "", "The kyverno report file(yaml)")
	rootCmd.Flags().StringVarP(&output, "output", "o", "report.xml", "report filename")
	if err := rootCmd.MarkFlagRequired("filename"); err != nil {
		fmt.Printf("error filename flag :%v", err)
	}
}
