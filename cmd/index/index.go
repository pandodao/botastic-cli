package index

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/pandodao/botastic-cli/cmd/core"
	"github.com/pandodao/botastic-go"
	"github.com/pandodao/tokenizer-go"
	"github.com/spf13/cobra"
)

var (
	action      string
	query       string
	indexesfile string
)

func NewCmdIndex() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index --act <action> [options]",
		Short: "create or search indexes",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			client := ctx.Value(core.CtxClient{}).(*botastic.Client)
			if action == "create" {
				if indexesfile == "" {
					cmd.PrintErrln("missing required flag: --file")
					os.Exit(-1)
				}

				buf, err := os.ReadFile(indexesfile)
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(-1)
				}

				indexes := &botastic.CreateIndexesRequest{}
				if err := json.Unmarshal(buf, &indexes); err != nil {
					cmd.PrintErrln(err)
					os.Exit(-1)
				}

				tokens := make([]int64, len(indexes.Items))
				for i, item := range indexes.Items {
					t := tokenizer.MustCalToken(item.Data)
					if t > 3500 {
						cmd.PrintErrf("data too long, objectId: %s\n", item.ObjectID)
						os.Exit(-1)
					}
					tokens[i] = t
				}

				req := botastic.CreateIndexesRequest{}
				tokenSum, lastOne, start := int64(0), false, 0
				for i := 0; i < len(indexes.Items)+1; i++ {
					token := int64(0)
					if i < len(tokens) {
						token = tokens[i]
					}
					if tokenSum+token > 8100 || lastOne {
						err := client.CreateIndexes(ctx, req)
						if err != nil {
							cmd.PrintErrln(err)
							os.Exit(-1)
						}
						cmd.Printf("📝 chunk %d~%d, %d indexes created, token: %d.\n", start, i-1, len(req.Items), tokenSum)
						req.Items = []*botastic.CreateIndexesItem{}
						tokenSum = 0
						start = i
					}
					if !lastOne {
						req.Items = append(req.Items, indexes.Items[i])
						tokenSum += token
						if i == len(indexes.Items)-1 {
							lastOne = true
						}
					}
				}

				cmd.Printf("✅ done. %d indexes created.\n", len(indexes.Items))

			} else if action == "search" {
				if query == "" {
					cmd.PrintErrln("missing required flag: --query")
					os.Exit(-1)
				}

				resp, err := client.SearchIndexes(ctx, botastic.SearchIndexesRequest{
					Keywords: query,
					N:        3,
				})
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(-1)
				}
				for ix, item := range resp.Items {
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
	cmd.Flags().StringVar(&indexesfile, "file", "", "indexes file path. Only valid when action is create")

	return cmd
}
