/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
    "os"

    "github.com/spf13/cobra"
    "k8s.io/cli-runtime/pkg/genericclioptions"
)

var KubernetesConfigFlags *genericclioptions.ConfigFlags

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "kubectl-x",
    Short: "",
    Long:  "",
}

func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}

func init() {

    KubernetesConfigFlags = genericclioptions.NewConfigFlags(true)
    rootCmd.Flags().BoolP("x", "x", false, "Help message for toggle")
    KubernetesConfigFlags.AddFlags(rootCmd.PersistentFlags())
}
