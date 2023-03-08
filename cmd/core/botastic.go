package core

type (
	IndicesRequest struct {
		Items []*IndicesRequestItem `json:"items"`
	}
	IndicesRequestItem struct {
		ObjectID   string `json:"object_id"`
		Category   string `json:"category"`
		Data       string `json:"data"`
		Properties string `json:"properties"`
	}

	SearchResult struct {
		Ts   int64               `json:"ts"`
		Data []*SearchResultItem `json:"data"`
	}

	SearchResultItem struct {
		ObjectID   string  `json:"object_id"`
		Data       string  `json:"data"`
		Score      float64 `json:"score"`
		CreatedAt  int64   `json:"created_at"`
		Category   string  `json:"category"`
		Properties string  `json:"properties"`
	}
)

type (
	CtxHost         struct{}
	CtxBotasticAuth struct{}
)
