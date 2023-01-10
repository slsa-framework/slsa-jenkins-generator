package main

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"slsa-jenkins-generator/provenance"
)

var Command = &cobra.Command{
	Use:   "Provenance Generator",
	Short: "Generate Provenance",
	Long:  "Generate SLSA Provenance",
	Run:   generator,
}

func init() {

	Command.Flags().StringP(
		"artifact_path",
		"a",
		".",
		"path to artifact")

	Command.Flags().StringP(
		"output_path",
		"o",
		".",
		"path to provenance")

	Command.MarkFlagRequired("artifact_path")

}

func generator(cmd *cobra.Command, args []string) {

	Opt := func(opt string) string {
		cmdOpt, err := cmd.Flags().GetString(opt)
		if err != nil {
			print("Failed to read command option %v", err)
		}
		return cmdOpt
	}

	artifact := Opt("artifact_path")

	outpath := Opt("output_path")

	provenance.GenerateSLSA(artifact, outpath)

}

func main() {
	if err := Command.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
