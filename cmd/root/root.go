package root

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/pandodao/botastic-cli/cmd/core"
	"github.com/pandodao/botastic-cli/cmd/index"
	"github.com/pandodao/botastic-cli/cmd/scan"

	"github.com/spf13/cobra"
)

func NewCmdRoot(version string) *cobra.Command {
	var opt struct {
		host string
	}

	cmd := &cobra.Command{
		Use:           "botastic-cli <command> <subcommand> [flags]",
		Short:         "bc",
		Long:          `a command line tool for botastic`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if cmd.Flags().Changed("host") {
				u, err := url.Parse(opt.host)
				if err != nil {
					return err
				}

				if u.Scheme == "" {
					u.Scheme = "https"
				}

				ctx = context.WithValue(ctx, core.CtxHost{}, u.String())
			} else {
				ctx = context.WithValue(ctx, core.CtxHost{}, opt.host)
			}

			appID := os.Getenv("BOTASTIC_APP_ID")
			appSecret := os.Getenv("BOTASTIC_SECRET")
			if appID == "" || appSecret == "" {
				fmt.Println("environment variable BOTASTIC_APP_ID or BOTASTIC_SECRET is empty")
				os.Exit(1)
			}

			ctx = context.WithValue(ctx, core.CtxBotasticAuth{}, fmt.Sprintf("%s:%s", appID, appSecret))
			cmd.SetContext(ctx)

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&opt.host, "host", "https://botastic-api.pando.im/api", "custom api host")

	cmd.AddCommand(scan.NewCmdScan())
	cmd.AddCommand(index.NewCmdIndex())

	return cmd
}
