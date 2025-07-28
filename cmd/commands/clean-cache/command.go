package clean_cache_command

import (
	"os"

	ip_cache_module "github.com/diogopereiradev/httpzen/internal/cache"
	logger_module "github.com/diogopereiradev/httpzen/internal/logger"
	"github.com/spf13/cobra"
)

var Exit = os.Exit
var LoggerSuccess = logger_module.Success

var IpClearCache = ip_cache_module.ClearCache

func Init(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "cleancache",
		Short: "Clear the app cache, like the IP cache and other temporary data",
		Run: func(cmd *cobra.Command, args []string) {
			IpClearCache()
			LoggerSuccess("Your cache was cleared successfully!", 50)
			Exit(0)
		},
	}
	rootCmd.AddCommand(cmd)
}
