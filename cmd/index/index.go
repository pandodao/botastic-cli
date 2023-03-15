package index

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/pandodao/botastic-cli/cmd/core"
	"github.com/pandodao/botastic-go"
	"github.com/spf13/cobra"
)

var (
	action      string
	query       string
	indicesfile string
)

func NewCmdIndex() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index --act <action> [options]",
		Short: "create or search indices",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			client := ctx.Value(core.CtxClient{}).(*botastic.Client)
			if action == "create" {
				if indicesfile == "" {
					cmd.PrintErrln("missing required flag: --file")
					os.Exit(-1)
				}

				buf, err := os.ReadFile(indicesfile)
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(-1)
				}

				indices := &botastic.CreateIndicesRequest{}
				if err := json.Unmarshal(buf, &indices); err != nil {
					cmd.PrintErrln(err)
					os.Exit(-1)
				}

				// split indices into chunks
				chunks := make([]botastic.CreateIndicesRequest, 0)
				chunkSize := 128
				for i := 0; i < len(indices.Items); i += chunkSize {
					end := i + chunkSize
					if end > len(indices.Items) {
						end = len(indices.Items)
					}
					chunks = append(chunks, botastic.CreateIndicesRequest{
						Items: indices.Items[i:end],
					})
				}

				for ix, chunk := range chunks {
					err := client.CreateIndices(ctx, chunk)
					if err != nil {
						cmd.PrintErrln(err)
						continue
					}
					cmd.Printf("📝 chunk %d, %d indices created.\n", ix+1, len(chunk.Items))
				}

				cmd.Printf("✅ done. %d indices created.\n", len(indices.Items))

			} else if action == "search" {
				if query == "" {
					cmd.PrintErrln("missing required flag: --query")
					os.Exit(-1)
				}

				resp, err := client.SearchIndices(ctx, botastic.SearchIndicesRequest{
					Keywords: query,
					N:        3,
				})
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(-1)
				}
				for ix, item := range resp.Indices {
					cmd.Printf("💡 Result #%d (%f):\n%s\nprop: %s\n\n", ix+1, item.Score, strings.TrimSpace(item.Data), item.Properties)
				}
			} else {
				cmd.Help()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&action, "act", "", "action to perform. [create, search]")
	cmd.Flags().StringVar(&query, "query", "", "query to search. Only valid when action is search")
	cmd.Flags().StringVar(&indicesfile, "file", "", "indices file path. Only valid when action is create")

	return cmd
}
