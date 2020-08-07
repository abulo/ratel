package elasticsearch

// IndexResponse --
type IndexResponse struct {
	HTTPStatusCode int    `json:"httpStatusCode"`
	Index          string `json:"_index"`
	Type           string `json:"_type"`
	ID             string `json:"_id"`
	Version        int    `json:"_version"`
	Result         string `json:"result"`
}

// Doc --
type Doc struct {
	Index string  `json:"_index"`
	Type  string  `json:"_type"`
	ID    string  `json:"_id"`
	Score float64 `json:"_score"`
}

// SearchHits --
type SearchHits struct {
	Total struct {
		Value    int    `json:"value"`
		Relation string `json:"relation"`
	} `json:"total"`
	MaxScore float64 `json:"max_score"`
}

// SearchResponse --
type SearchResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
}
