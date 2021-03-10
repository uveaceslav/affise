package model

type RequestEntity struct {
	URLs []string `json:"urls"`
}

type ProcessedEntity struct {
	URL  string      `json:"url"`
	Data interface{} `json:"data"`
}

type ResponseEntity struct {
	Data []ProcessedEntity `json:"data"`
}
