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
		Type string `json:"type"`
		raw  []byte
	}

	ReceiveContent string

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
	v.raw = r
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
	v := struct {
		Type string `json:"type"`
	}{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	a.Type = v.Type
	copy(a.raw, b)
	return nil
}

func (a *Attachment) Photo() (*Photo, error) {
	if a.Type != AT_Photo {
		return nil, errors.New("vk: not photo json")
	}
	v := new(Photo)
	if err := json.Unmarshal(a.raw, v); err != nil {
		return nil, err
	}
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

// Video will only retrieve the preview image of the video. Only direct
// authorization is allowed to get the whole video content from VK
// server.
func (a *Attachment) Video() (*Video, error) {
	if a.Type != AT_Video {
		return nil, errors.New("vk: not video json")
	}
	v := new(Video)
	if err := json.Unmarshal(a.raw, v); err != nil {
		return nil, err
	}
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
	v := new(Audio)
	if err := json.Unmarshal(a.raw, v); err != nil {
		return nil, err
	}
	if v.Url != "" {
		v.ReceiveContent = ReceiveContent(v.Url)
	}
	return v, nil
}

func (a *Attachment) Doc() (*Doc, error) {
	if a.Type != AT_Doc {
		return nil, errors.New("vk: not doc json")
	}
	v := new(Doc)
	if err := json.Unmarshal(a.raw, v); err != nil {
		return nil, err
	}
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
