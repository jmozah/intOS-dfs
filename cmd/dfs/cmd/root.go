/*
Copyright © 2020 intOS Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	beeHost   string
	beePort   string
	httpPort  string
	verbosity string
	dataDir   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dfs",
	Short: "Decentralised file system over Swarm(https://ethswarm.org/)",
	Long: `dfs is the file system layer of internetOS. It is a thin layer over Swarm.  
It adds features to Swarm that is required by the internetOS to parallelize computation of data. 
It manages the metadata of directories and files created and expose them to higher layers.
It can also be used as a standalone personal, decentralised drive over the internet`,

	//Run: func(cmd *cobra.Command, args []string) {
	//},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	intOSdfs := `
	/$$             /$$      /$$$$$$   /$$$$$$                /$$  /$$$$$$
	|__/            | $$     /$$__  $$ /$$__  $$              | $$ /$$__  $$
	/$$ /$$$$$$$  /$$$$$$  | $$  \ $$| $$  \__/          /$$$$$$$| $$  \__//$$$$$$$
	| $$| $$__  $$|_  $$_/  | $$  | $$|  $$$$$$  /$$$$$$ /$$__  $$| $$$$   /$$_____/
	| $$| $$  \ $$  | $$    | $$  | $$ \____  $$|______/| $$  | $$| $$_/  |  $$$$$$
	| $$| $$  | $$  | $$ /$$| $$  | $$ /$$  \ $$        | $$  | $$| $$     \____  $$
	| $$| $$  | $$  |  $$$$/|  $$$$$$/|  $$$$$$/        |  $$$$$$$| $$     /$$$$$$$/
	|__/|__/  |__/   \___/   \______/  \______/          \_______/|__/    |_______/	
`
	fmt.Println(intOSdfs)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.intOS/dfs.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defaultDataDir := filepath.Join(home, ".intos/dfs")
	rootCmd.PersistentFlags().StringVar(&dataDir, "dataDir", defaultDataDir, "store data in this dir (default ~/.intos/dfs)")
	rootCmd.PersistentFlags().StringVar(&beeHost, "beeHost", "127.0.0.1", "bee host (default 127.0.0.1)")
	rootCmd.PersistentFlags().StringVar(&beePort, "beePort", "8080", "bee port (default 8080)")
	rootCmd.PersistentFlags().StringVar(&httpPort, "httpPort", "9090", "http port (default 9090)")
	rootCmd.PersistentFlags().StringVar(&verbosity, "verbosity", "5", "verbosity level (default 4)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home dir.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home dir with name ".dfs" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".intOS/dfs")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
