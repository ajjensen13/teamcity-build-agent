/*
Copyright © 2020 A. Jensen <jensen.aaro@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package buildtools

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "buildtools",
	Short: "Provides build tools for TeamCity build agents",
	Long:  `Provides build tools for TeamCity build agents`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if out != "" {
			log.Printf("buildtools: output is being directed to %s", out)

			d := filepath.Dir(out)
			fi, err := os.Stat(d)
			switch {
			case os.IsNotExist(err):
				log.Printf("buildtools: directory %s does not exist. It will be created", d)
				err = os.MkdirAll(d, 0755)
				if err != nil {
					return fmt.Errorf("buildtools: error creating directory: %w", err)
				}
			case err != nil:
				return fmt.Errorf("buildtools: error gathering info for output directory: %w", err)
			case !fi.IsDir():
				return fmt.Errorf("%s exists but is not a directory", fi.Name())
			}

			f, err := os.Create(out)
			if err != nil {
				return fmt.Errorf("buildtools: error creating output file %s: %w", out, err)
			}
			os.Stdout = f
		}
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) (err error) {
		if out != "" {
			defer os.Stdout.Close()
			_ = os.Stdout.Sync()
		}
		return
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	out string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&out, "out", "o", "", "output file (default STDOUT)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".teamcity" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".teamcity")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
