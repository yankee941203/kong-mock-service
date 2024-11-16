package model

type MetaInfo struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	TotalPages int    `json:"total_pages"`
	TotalItems int    `json:"total_items"`
	PrevPage   string `json:"prev_page"`
	NextPage   string `json:"next_page"`
	LastPage   string `json:"last_page"`
}

type ResponseInfo struct {
	Meta MetaInfo      `json:"meta,omitempty"`
	Data []ServiceInfo `json:"data"`
}
