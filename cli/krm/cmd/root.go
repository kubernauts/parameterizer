package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	releaseVersion string
	cfgFile        string
	rootCmd        = &cobra.Command{
		Use:   "krm",
		Short: "A tool for Kubernetes app definitions parameterization",
		Long: `A tool that takes Kubernetes app definitions as an input along with
user-defined parameters to create a Kubernetes YAML manifest that you
can feed into an installer.`,
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init initializes the tool
func init() {
	cobra.OnInitialize(initConfig)
	// Global flags and configuration settings:
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.krm.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// use config file from the flag:
		viper.SetConfigFile(cfgFile)
	} else {
		// find home directory:
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// search config in home dir with name ".krm" (without extension):
		viper.AddConfigPath(home)
		viper.SetConfigName(".krm")
	}
	// read in environment variables that match:
	viper.AutomaticEnv()
	// if a krm config file is found, read it in:
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
