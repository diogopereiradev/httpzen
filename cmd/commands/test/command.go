package test_command

import (
	"fmt"
	"time"

	request_module "github.com/diogopereiradev/httpzen/internal/request"
	"github.com/spf13/cobra"
)

func Executor() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command for testing the application functionality",
		Run: func(cmd *cobra.Command, args []string) {
			info := request_module.HandleRequest(request_module.RequestOptions{
				Timeout: 10 * time.Second,
				Url:     "https://google.com",
				Method:  "GET",
			})
			fmt.Println(info)
		},
	}
	return cmd
}
