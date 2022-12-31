package xtype

type Gender int

const (
	GenderNone = iota
	GenderMale
	GenderFemale
)

type PhoneNumber struct {
	Code      int32  `json:"code"`
	Number    int64  `json:"number"`
	Extension string `json:"extension,omitempty"`
}
