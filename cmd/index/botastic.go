package index

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pandodao/botastic-cli/cmd/core"
)

func createIndices(ctx context.Context, payload *core.IndicesRequest) ([]byte, error) {
	apiBase := ctx.Value(core.CtxHost{}).(string)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/indices", apiBase), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	buildHeader(ctx, req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

func searchIndices(ctx context.Context, query string, n int) (*core.SearchResult, error) {
	apiBase := ctx.Value(core.CtxHost{}).(string)

	params := url.Values{}
	params.Add("keywords", query)
	params.Add("n", strconv.Itoa(n))

	encoded := params.Encode()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/indices/search?%s", apiBase, encoded), nil)
	if err != nil {
		return nil, err
	}

	buildHeader(ctx, req)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &core.SearchResult{}
	json.Unmarshal(body, &result)

	return result, nil
}

func buildHeader(ctx context.Context, req *http.Request) error {
	authString := ctx.Value(core.CtxBotasticAuth{}).(string)
	parts := strings.Split(authString, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid auth string: %s", authString)
	}

	req.Header.Set("X-BOTASTIC-APPID", parts[0])
	req.Header.Set("X-BOTASTIC-SECRET", parts[1])
	req.Header.Set("Content-Type", "application/json")
	return nil
}
