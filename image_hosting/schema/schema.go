package schema

type CheveretoImageInfo struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int    `json:"size,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	// Md5       string `json:"md5"`
	Filename string `json:"filename"`
	Mime     string `json:"mime"`
	Url      string `json:"url"`
	// Thumb     struct {
	// 	Url string `json:"url"`
	// } `json:"thumb"`
	DisplayUrl string `json:"display_url"`
}

type CheveretoUploadResponse struct {
	StatusCode int                `json:"status_code"`
	StatusTxt  string             `json:"status_txt"`
	Image      CheveretoImageInfo `json:"image"`
}
