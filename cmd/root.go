/*
Copyright © 2021 tacusci ltd

Licensed under the GNU GENERAL PUBLIC LICENSE Version 3 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.gnu.org/licenses/gpl-3.0.html

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tacusci/berrycms/web/config"
	"github.com/tacusci/berrycms/web/server"
)

var opts = config.Options{}

var rootCmd = &cobra.Command{
	Use:   "berry",
	Short: "A modern extensible CMS",
	// long description pending
	// 	Long: `A longer description that spans multiple lines and likely contains
	// examples and usage of using your application. For example:

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		startup(opts)
	},
}

func startup(o config.Options) {
	// launches CMS, blocks main thread and waits for shutdown
	<-server.Bootup(o)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().UintVarP(&opts.Port, "port", "p", 8080, "Port to listen for HTTP requests on")
	rootCmd.Flags().StringVarP(&opts.Addr, "address", "a", "0.0.0.0", "IP address to listen against if multiple network adapters")
	rootCmd.Flags().StringVarP(&opts.Sql, "database", "d", "sqlite", "Database server type to try to connect to [sqlite/mysql]")
	rootCmd.Flags().StringVar(&opts.SqlUsername, "dbuser", "berryadmin", "Database server username, ignored if using sqlite")
	rootCmd.Flags().StringVar(&opts.SqlPassword, "dbpass", "", "Database server password, ignored if using sqlite")
	rootCmd.Flags().StringVar(&opts.SqlAddress, "dbaddr", "/", "Database server location, ignored if using sqlite")
	rootCmd.Flags().StringVar(&opts.ActivityLogLoc, "actlog", "", "Activity/access log file location")
	rootCmd.Flags().StringVarP(&opts.AdminHiddenPassword, "adminhiddenpassword", "x", "", "URI prefix to hide admin pages behind")
}
