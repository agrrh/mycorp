package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ledongthuc/goterators"
	"github.com/spf13/cobra"

	"github.com/agrrh/mycorp/internal/domain/scenario"
	"github.com/agrrh/mycorp/internal/domain/scenario_store"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mycorp",
		Short: "MyCorp CLI Tool",
	}

	baseURL := os.Getenv("MYCORP_CLI_URL")
	if baseURL == "" {
		fmt.Fprintf(os.Stderr, "Environment variable MYCORP_CLI_URL is required\n")
		os.Exit(1)
	}

	scenariosURL := fmt.Sprintf("%s/scenarios", baseURL)

	scStore := scenario_store.NewCLI(scenariosURL)

	if err := scStore.Fetch(); err != nil {
		log.Fatal(err)
	}

	commandGroups := goterators.Group(scStore.List(), func(sc *scenario.ScenarioCLI) string {
		return sc.Metadata.Namespace
	})

	// Register commands dynamically
	for _, group := range commandGroups {
		namespaceCmd := &cobra.Command{
			Use: group[0].Metadata.Namespace,
		}

		for _, sc := range group {
			// TODO: Generate usage string based on inputs

			cmd := &cobra.Command{
				Use: sc.Metadata.Name,
			}

			// Add input parameters as flags
			for _, input := range sc.Spec.Inputs {
				switch input.Type {
				case "string":
					cmd.Flags().String(input.Name, fmt.Sprintf("%v", input.Default), input.GetCLIDescription())
				case "bool":
					cmd.Flags().Bool(input.Name, false, input.GetCLIDescription())
				case "int":
					cmd.Flags().Int(input.Name, 0, input.GetCLIDescription())
				default:
					cmd.Flags().String(input.Name, "", input.GetCLIDescription())
				}
			}

			// Set command execution logic
			cmd.Run = func(cmd *cobra.Command, args []string) {
				fmt.Printf("Executing %s ...\n\n", sc.Metadata.GetFullName())

				output, err := sc.Run(fmt.Sprintf("%s/%s", scenariosURL, sc.Metadata.GetFullName()))
				if err != nil {
					fmt.Printf("Error running command %s: %s\n", sc.Metadata.GetFullName(), err)
					os.Exit(1)
				}

				fmt.Println(output)
				os.Exit(0)
			}

			namespaceCmd.AddCommand(cmd)
		}

		rootCmd.AddCommand(namespaceCmd)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
