package vk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	// attachment type
	AT_Photo       = "photo"
	AT_Sticker     = "sticker"
	AT_PostedPhoto = "posted_photo"
	AT_Video       = "video"
	AT_Audio       = "audio"
	AT_Doc         = "doc"
	AT_Graffiti    = "graffiti"
	AT_Url         = "url"
	AT_Link        = "link"
	AT_Note        = "note"
	AT_App         = "app"
	AT_Poll        = "poll"
	AT_Page        = "page"
)

type (
	Attachment struct {
		Type string          `json:"type"`
		Pl   *Poll           `json:"poll"`
		P    *Photo          `json:"photo"`
		S    *Sticker        `json:"sticker"`
		V    *Video          `json:"video"`
		A    *Audio          `json:"audio"`
		D    *Doc            `json:"doc"`
		G    json.RawMessage `json:"geo"`
	}

	ReceiveContent string

	Poll struct {
		Id       int    `json:"id"`
		OwnerId  int    `json:"owner_id"`
		Created  int64  `json:"created"`
		Question string `json:"question"`
		Votes    int    `json:"votes"`
		AnswerId int    `json:"answer_id"`
		Answers  []struct {
			Id    int     `json:"id"`
			Text  string  `json:"text"`
			Votes int     `json:"votes"`
			Rate  float32 `json:"rate"`
		} `json:"answers"`
		Anonymous int `json:"anonymous"`
	}

	Photo struct {
		ReceiveContent
		Id        int    `json:"id"`
		AlbumId   int    `json:"album_id"`
		OwnerId   int    `json:"owner_id"`
		Photo75   string `json:"photo_75"`
		Photo130  string `json:"photo_130"`
		Photo604  string `json:"photo_604"`
		Photo807  string `json:"photo_807"`
		Photo1280 string `json:"photo_1280"`
		Photo2560 string `json:"photo_2560"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		Text      string `json:"text"`
		Date      int64  `json:"date"`
		AccessKey string `json:"access_key"`
	}

	Sticker struct {
		ReceiveContent
		Id        int    `json:"id"`
		ProductId int    `json:"product_id"`
		Photo64   string `json:"photo_64"`
		Photo128  string `json:"photo_128"`
		Photo256  string `json:"photo_256"`
		Photo352  string `json:"photo_352"`
		Photo512  string `json:"photo_512"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
	}

	Video struct {
		ReceiveContent
		Id          int    `json:"id"`
		OwnerId     int    `json:"owner_id"`
		Title       string `json:"title"`
		Duration    int    `json:"duration"`
		Description string `json:"description"`
		Date        int64  `json:"date"`
		Views       int    `json:"views"`
		Comments    int    `json:"comments"`
		Photo130    string `json:"photo_130"`
		Photo320    string `json:"photo_320"`
		Photo800    string `json:"photo_800"`
		AccessKey   string `json:"access_key"`
		CanEdit     Bool   `json:"can_edit"`
		CanAdd      Bool   `json:"can_add"`
	}

	// Audio refer to audio.go

	Doc struct {
		ReceiveContent
		Id        int    `json:"id"`
		OwnerId   int    `json:"owner_id"`
		Title     string `json:"title"`
		Size      int    `json:"size"`
		Ext       string `json:"ext"`
		Url       string `json:"url"`
		Date      int64  `json:"date"`
		Type      int    `json:"type"`
		AccessKey string `json:"access_key"`
	}
)

// Another approach to get the Attachment object from raw bytes string.
// r is considered given to this method.
func GetAttachment(r []byte) (*Attachment, error) {
	v := new(Attachment)
	if err := json.Unmarshal(r, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (r ReceiveContent) String() string {
	return string(r)
}

func (r ReceiveContent) Content() (io.ReadCloser, error) {
	if r == "" {
		return nil, errors.New("can not get attachment content because empty url.")
	}
	req, err := http.NewRequest("GET", r.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// Implement json.Unmarshaler, so the parsing on Post can be parsed directly.
// TODO: use Decoder struct to implment attachment unmarshal.
func (a *Attachment) UnmarshalJSON(b []byte) error {
	var v struct {
		Type string          `json:"type"`
		Pl   *Poll           `json:"poll"`
		P    *Photo          `json:"photo"`
		S    *Sticker        `json:"sticker"`
		V    *Video          `json:"video"`
		A    *Audio          `json:"audio"`
		D    *Doc            `json:"doc"`
		G    json.RawMessage `json:"geo"`
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	a.Type = v.Type
	a.Pl = v.Pl
	a.P = v.P
	a.S = v.S
	a.V = v.V
	a.A = v.A
	a.D = v.D
	a.G = v.G
	return nil
}

func (a *Attachment) Photo() (*Photo, error) {
	if a.Type != AT_Photo {
		return nil, errors.New("vk: not photo json")
	}
	v := a.P
	var s, url string
	if v.Photo2560 != "" {
		url = v.Photo2560
	} else if v.Photo1280 != "" {
		url = v.Photo1280
	} else if v.Photo807 != "" {
		url = v.Photo807
	} else if v.Photo604 != "" {
		url = v.Photo604
	} else if v.Photo130 != "" {
		url = v.Photo130
	} else if v.Photo75 != "" {
		url = v.Photo75
	}
	if url != "" {
		if v.AccessKey != "" {
			if n := strings.Index(url, "?"); n < 0 {
				s = fmt.Sprint("?access_key=", v.AccessKey)
			} else {
				s = fmt.Sprint("&access_key=", v.AccessKey)
			}
			url = fmt.Sprint(url, s)
		}
		v.ReceiveContent = ReceiveContent(url)
	}
	return v, nil
}

func (a *Attachment) Sticker() (*Sticker, error) {
	if a.Type != AT_Sticker {
		return nil, errors.New("vk: not sticker json")
	}
	v := a.S
	var url string
	if v.Photo512 != "" {
		url = v.Photo512
	} else if v.Photo352 != "" {
		url = v.Photo352
	} else if v.Photo256 != "" {
		url = v.Photo256
	} else if v.Photo128 != "" {
		url = v.Photo128
	} else if v.Photo64 != "" {
		url = v.Photo64
	}
	if url != "" {
		v.ReceiveContent = ReceiveContent(url)
	}
	return v, nil
}

// Video will only retrieve the preview image of the video. Only direct
// authorization is allowed to get the whole video content from VK
// server.
func (a *Attachment) Video() (*Video, error) {
	if a.Type != AT_Video {
		return nil, errors.New("vk: not video json")
	}
	v := a.V
	var s, url string
	if v.Photo800 != "" {
		url = v.Photo800
	} else if v.Photo320 != "" {
		url = v.Photo320
	} else if v.Photo130 != "" {
		url = v.Photo130
	}
	if url != "" {
		if v.AccessKey != "" {
			if n := strings.Index(url, "?"); n < 0 {
				s = fmt.Sprint("?access_key=", v.AccessKey)
			} else {
				s = fmt.Sprint("&access_key=", v.AccessKey)
			}
			url = fmt.Sprint(url, s)
		}
		v.ReceiveContent = ReceiveContent(url)
	}
	return v, nil
}

func (a *Attachment) Audio() (*Audio, error) {
	if a.Type != AT_Audio {
		return nil, errors.New("vk: not audio json")
	}
	v := a.A
	if v.Url != "" {
		v.ReceiveContent = ReceiveContent(v.Url)
	}
	return v, nil
}

func (a *Attachment) Doc() (*Doc, error) {
	if a.Type != AT_Doc {
		return nil, errors.New("vk: not doc json")
	}
	v := a.D
	var s, url string
	if v.Url != "" {
		url = v.Url
		if v.AccessKey != "" {
			if n := strings.Index(url, "?"); n < 0 {
				s = fmt.Sprint("?access_key=", v.AccessKey)
			} else {
				s = fmt.Sprint("&access_key=", v.AccessKey)
			}
			url = fmt.Sprint(url, s)
		}
		v.ReceiveContent = ReceiveContent(url)
	}
	return v, nil
}

func (a *Attachment) Poll() (*Poll, error) {
	if a.Type != AT_Poll {
		return nil, errors.New("vk: not poll json")
	}
	return a.Pl, nil
}
