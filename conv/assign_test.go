package conv_test

import (
	"encoding/json"
	"testing"
	"time"

	"code.olapie.com/sugar/v2/jsonutil"
	"code.olapie.com/sugar/v2/must"
)

type Image struct {
	Width  int
	Height int
	Link   string
}

type Topic struct {
	Title      string
	CoverImage *Image
	MoreImages []*Image
	CreatedAt  time.Time
}

func TestAssign(t *testing.T) {
	params := map[string]any{
		"title": "this is title",
		"cover_image": map[string]any{
			"w":    100,
			"h":    200,
			"link": "https://www.image.com",
		},
		"more_images": []map[string]any{
			{
				"w":    100,
				"h":    200,
				"link": "https://www.image.com",
			},
		},
		"created_at": "2020-12-06T12:46:15.134526-05:00",
	}

	var topic Topic
	var i any = topic
	err := rt.Assign(&i, params)
	if err != nil {
		t.FailNow()
	}
	t.Logf("%#v", topic)
	t.Logf("%#v", i)
}

func TestAssignSlice(t *testing.T) {
	params := map[string]any{
		"title": "this is title",
		"cover_image": map[string]any{
			"w":    100,
			"h":    200,
			"link": "https://www.image.com",
		},
		"more_images": []map[string]any{
			{
				"w":    100,
				"h":    200,
				"link": "https://www.image.com",
			},
		},
	}

	values := []any{params}
	var topics []*Topic
	err := rt.Assign(&topics, values)
	if err != nil || len(topics) == 0 {
		t.FailNow()
	}
}

type User struct {
	Id       int
	Name     string
	OpenAuth *OpenAuth
}

type OpenAuth struct {
	Provider string
	OpenID   string
}

type UserInfo struct {
	Id       int
	Name     string
	OpenAuth *OpenAuthInfo
}

type OpenAuthInfo struct {
	Provider string
	OpenID   string
}

func TestAssignStruct(t *testing.T) {
	user := &User{}
	userInfo := &UserInfo{
		Id:   1,
		Name: "tom",
		OpenAuth: &OpenAuthInfo{
			Provider: "wechat",
			OpenID:   "open_id_123",
		},
	}

	err := rt.Assign(user, userInfo)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%#v", user)
}

func TestAssignJSONToStruct(t *testing.T) {
	type Item struct {
		ID        int64     `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	}
	tm := time.Now().Add(time.Hour)
	i := new(Item)
	jsonMap := map[string]any{"id": 10, "created_at": tm}
	err := rt.Assign(i, jsonMap)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if !tm.Equal(i.CreatedAt) {
		t.Errorf("Expected CreatedAt to be %v instead of %v", tm, i.CreatedAt)
		t.FailNow()
	}

	i = new(Item)
	jsonStr := jsonutil.ToString(tm)
	jsonStr = jsonStr[1 : len(jsonStr)-1]
	jsonMap = map[string]any{"id": 10, "created_at": jsonStr}
	err = rt.Assign(i, jsonMap)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if !tm.Equal(i.CreatedAt) {
		t.Errorf("Expected CreatedAt to be %v instead of %v", tm, i.CreatedAt)
		t.FailNow()
	}
}

type jsonMap map[int]int

func (j *jsonMap) UnmarshalJSON(bytes []byte) error {
	var m map[string]int
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		(*j)[must.ToInt(k)] = v
	}
	return nil
}

var _ json.Unmarshaler = (*jsonMap)(nil)

func TestAssignMap(t *testing.T) {
	var dst jsonMap
	var src = map[string]int{"1": 2}
	err := rt.Assign(&dst, src)
	if err != nil {
		t.Error(err)
	}
	s1, s2 := jsonutil.ToString(dst), jsonutil.ToString(src)
	if s1 != s2 {
		t.Errorf("mismatch: \n%s\n%s\n", s1, s2)
	}
}
