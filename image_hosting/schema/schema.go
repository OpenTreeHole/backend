package schema

// type LskyUploadResponse struct {
//	Status  bool   `json:"status"`
//	Message string `json:"message"`
//	Data    struct {
//		Key        string  `json:"key"`
//		Name       string  `json:"name"`
//		Pathname   string  `json:"pathname"`
//		OriginName string  `json:"origin_name"`
//		Size       float64 `json:"size"`
//		Mimetype   string  `json:"mimetype"`
//		Extension  string  `json:"extension"`
//		Md5        string  `json:"md5"`
//		Sha1       string  `json:"sha1"`
//		Links      struct {
//			Url              string `json:"url"`
//			Html             string `json:"html"`
//			Bbcode           string `json:"bbcode"`
//			Markdown         string `json:"markdown"`
//			MarkdownWithLink string `json:"markdown_with_link"`
//			ThumbnailUrl     string `json:"thumbnail_url"`
//		} `json:"links"`
//	} `json:"data"`
// }

type CheveretoImageInfo struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int    `json:"size,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Md5       string `json:"md5"`
	Filename  string `json:"filename"`
	Mime      string `json:"mime"`
	Url       string `json:"url"`
	Thumb     struct {
		Url string `json:"url"`
	} `json:"thumb"`
	DisplayUrl string `json:"display_url"`
}

type CheveretoUploadResponse struct {
	StatusCode int                `json:"status_code"`
	StatusTxt  string             `json:"status_txt"`
	Image      CheveretoImageInfo `json:"image"`
}
