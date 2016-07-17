// media objects in VKontakte - video and photo (see audio.go for audio)

package vk

type (
	// https://new.vk.com/dev/photo
	PhotoObj struct {
		Pid        int    `json:"pid"`
		Aid        int    `json:"aid"` // album ID
		OwnerId    int    `json:"owner_id"`
		Src        string `json:"src"`
		SrcBig     string `json:"src_big"`
		SrcSmall   string `json:"src"`
		SrcxXBig   string `json:"src_xbig"`
		SrcxXXBig  string `json:"src_xxbig"`
		SrcxXXXBig string `json:"src_xxxbig"`
		Width      int    `json:"width"`
		Height     int    `json:"height"`
		Text       string `json:"text"`
		Created    int64  `json:"created"`
	}

	// https://new.vk.com/dev/video_object
	VideoObj struct {
		Id          int    `json:"id"`
		Vid         int    `json:"vid"`
		OwnerId     int    `json:"owner_id"`
		Title       string `json:"title"`
		Duration    int    `json:"duration"`
		Description string `json:"description"`
		Link        string `json:"link"`
		Photo130    string `json:"photo_130"`
		Photo320    string `json:"photo_320"`
		Photo640    string `json:"photo_640"`
		Date        int64  `json:"date"`
		AddDate     int64  `json:"adding_date"`
		Views       int    `json:"views"`
		Player      string `json:"player"`
	}
)
