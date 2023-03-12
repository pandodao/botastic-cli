package scan

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pandodao/botastic-cli/cmd/core"

	"github.com/spf13/cobra"
)

var (
	scanDir        string
	fileType       string
	scanMode       string
	supportedTypes = []string{"txt", "md"}
)

func NewCmdScan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan --dir <message> --type <type>",
		Short: "scan a directory for files with a specific type and create indices for them",
		RunE: func(cmd *cobra.Command, args []string) error {
			if scanDir == "" {
				cmd.PrintErrln("missing required flag: --dir")
				os.Exit(-1)
			}

			if fileType == "" {
				cmd.PrintErrln("missing required flag: --type")
				cmd.PrintErrln("supported types:", strings.Join(supportedTypes, ", "))
				os.Exit(-1)
			}

			files, err := scanDirectory(scanDir, fileType)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(-1)
			}

			req := &core.IndicesRequest{}
			req.Items = make([]*core.IndicesRequestItem, 0)
			for _, file := range files {
				var items []*core.IndicesRequestItem
				switch fileType {
				case "md":
					{
						extractFn := extractMardownFileByLine
						if scanMode == "paragraph" {
							extractFn = extractMardownFileByParagraph
						}
						items, err = extractFn(file)
						if err != nil {
							cmd.PrintErrln(err)
							continue
						}
					}
				case "txt":
					{
					}
				}

				cmd.Printf("🔍 scan file: %s, %d indices\n", file, len(items))
				req.Items = append(req.Items, items...)
			}

			if err = saveIndices(req); err != nil {
				cmd.PrintErrln(err)
				os.Exit(-1)
			}

			cmd.Printf("✅ done. save to file: %s\n", "indices.json")
			return nil
		},
	}

	cmd.Flags().StringVar(&scanDir, "dir", "", "the directory to be scanned for files.")
	cmd.Flags().StringVar(&fileType, "type", "", "the file type to be scanned for.")
	cmd.Flags().StringVar(&scanMode, "mode", "line", "the scan mode. supported modes: 'line', 'paragraph'.")

	return cmd
}

func scanDirectory(root, fileType string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == "."+fileType {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func saveIndices(req *core.IndicesRequest) error {
	f, err := os.Create("indices.json")
	if err != nil {
		return err
	}
	defer f.Close()

	buf, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
