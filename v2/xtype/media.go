package xtype

type Image struct {
	URL       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Format    string `json:"format"`
	Size      int64  `json:"size"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
	Data      []byte `json:"data"`
}

type Video struct {
	URL      string `json:"url"`
	Format   string `json:"format"`
	Duration int    `json:"duration"`
	Size     int64  `json:"size"`
	Image    *Image `json:"image"`
	Name     string `json:"name"`
	Data     []byte `json:"data"`
}

type Audio struct {
	URL      string `json:"url"`
	Format   string `json:"format"`
	Duration int    `json:"duration"`
	Size     int64  `json:"size"`
	Name     string `json:"name"`
	Data     []byte `json:"data"`
}
