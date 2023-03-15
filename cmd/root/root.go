package root

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/pandodao/botastic-cli/cmd/core"
	"github.com/pandodao/botastic-cli/cmd/index"
	"github.com/pandodao/botastic-cli/cmd/scan"
	"github.com/pandodao/botastic-go"

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

			host := opt.host
			if cmd.Flags().Changed("host") {
				u, err := url.Parse(opt.host)
				if err != nil {
					return err
				}

				if u.Scheme == "" {
					u.Scheme = "https"
				}

				host = u.String()
			}

			appID := os.Getenv("BOTASTIC_APP_ID")
			appSecret := os.Getenv("BOTASTIC_SECRET")
			if appID == "" {
				fmt.Println("environment variable BOTASTIC_APP_ID is empty")
				os.Exit(1)
			}
			if appSecret == "" {
				fmt.Println("environment variable BOTASTIC_SECRET is empty, you cannot create indexes")
			}

			botastic := botastic.New(appID, appSecret, botastic.WithHost(host))

			ctx = context.WithValue(ctx, core.CtxClient{}, botastic)
			cmd.SetContext(ctx)

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&opt.host, "host", "https://botastic-api.pando.im", "custom api host")

	cmd.AddCommand(scan.NewCmdScan())
	cmd.AddCommand(index.NewCmdIndex())

	return cmd
}
