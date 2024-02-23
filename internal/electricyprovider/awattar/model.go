package awattar

type AwattarResponse struct {
	Object string        `json:"object"`
	Data   []AwattarData `json:"data"`
	Url    string        `json:"url"`
}

type AwattarData struct {
	Start       int     `json:"start_timestamp"`
	End         int     `json:"end_timestamp"`
	Marketprice float64 `json:"marketprice"`
	Unit        string  `json:"unit"`
}
