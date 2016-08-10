package vk

import (
	"encoding/json"
	"net/url"
	"strconv"
)

const (
	PostSrc_VK       = "vk"
	PostSrc_Widget   = "widget"
	PostSrc_Api      = "api"
	PostSrc_Rss      = "rss"
	PostSrc_Sms      = "sms"
	PostType_Post    = "post"
	PostType_Suggest = "suggest"
)

type (
	// https://vk.com/dev/post
	Post struct {
		Id           int    `json:"id"`
		FromId       int    `json:"from_id"`
		OwnerId      int    `json:"owner_id"`
		Date         int64  `json:"date"`
		PostType     string `json:"post_type"`
		Text         string `json:"text"`
		CanEdit      Bool   `json:"can_edit"`
		CreatedBy    int    `json:"created_by"`
		CanDel       Bool   `json:"can_delete"`
		ReplyOwnerId int    `json:"reply_owner_id"`
		ReplyPostId  int    `json:"reply_post_id"`
		FriendsOnly  int    `json:"friends_only"`
		Comments     struct {
			Count   int  `json:"count"`
			CanPost Bool `json:"can_post"`
		} `json:"comments"`
		Likes struct {
			Count      int  `json:"count"`
			UserLikes  Bool `json:"user_likes"`
			CanLike    Bool `json:"can_like"`
			CanPublish Bool `json:"can_publish"`
		} `json:"likes"`
		Reposts struct {
			Count        int  `json:"count"`
			UserReposted Bool `json:"user_reposted"`
		} `json:"reposts"`
		PostSrc struct {
			Type string `json:"type"`
			Data string `json:"data"`
		} `json:"post_source"`
		Attachments []*Attachment `json:"attachments"`
		G           *Geo          `json:"geo"`
		SignerId    int           `json:"signer_id"`
		CopyHistory []Post        `json:"copy_history"`
		CanPin      Bool          `json:"can_pin"`
		IsPinned    Bool          `json:"is_pinned"`
	}

	Geo struct {
		Type  string `json:"type"`
		Coor  string `json:"coordinates"`
		Place struct {
			Id        int     `json:"id"`
			Title     string  `json:"title"`
			Latitude  float32 `json:"latitude"`
			Longitude float32 `json:"longitude"`
			Created   int64   `json:"created"`
			Icon      string  `json:"icon"`
			Country   string  `json:"country"`
			City      string  `json:"city"`
		} `json:"place"`
		Showmap Bool `json:"showmap"`
	}

	// https://vk.com/dev/comment_object
	Comment struct {
		Id          int           `json:"id"`
		FromId      int           `json:"from_id"`
		Date        int64         `json:"date"`
		Text        string        `json:"text"`
		ReplyToUid  int           `json:"reply_to_uid"`
		ReplyToCid  int           `json:"reply_to_cid"`
		Attachments []*Attachment `json:"attachments"`
		Likes       struct {
			Count     int  `json:"count"`
			UserLikes Bool `json:"user_likes"`
			CanLike   Bool `json:"can_like"`
		} `json:"likes"`
		PostOwnerId int `json:"post_owner_id"`
		PostId      int `json:"post_id"`
	}

	// https://vk.com/dev/notifications.get
	Topic struct {
		Id        int    `json:"id"`
		OwnerId   int    `json:"owner_id"`
		Title     string `json:"title"`
		Created   int64  `json:"created"`
		CreatedBy int    `json:"created_by"`
		Updated   int64  `json:"updated"`
		UpdatedBy int    `json:"updated_by"`
		IsClosed  Bool   `json:"is_closed"`
		IsFixed   Bool   `json:"is_fixed"`
		Comments  int    `json:"comments"`
	}
)

// Post wall to either user's wall or user's page (community) wall. Current arg
// wallOwnerOptions consist of 'idInt int' and 'isCommunity bool' options.
// idInt is the user ID or community ID. isCommunity is the flag to indicate if
// the idInt refering a normal user ID or community ID. m is the message that
// want to be posted on the wall. If wallOwnerOptions is not supplied then this
// will assume the current user owner of the session token. So, a group token
// will return error as only user token is valid for wall posting action.
func (s *Session) WallPost(m, a string, ownerOpts ...interface{}) (int, error) {
	vals := make(url.Values)
	vals.Set("message", m)
	if ownerOptions(vals, ownerOpts...) {
		vals.Set("from_group", "1")
	}
	if a != "" {
		vals.Set("attachments", a)
	}
	var n struct {
		PostId int `json:"post_id"`
	}
	if err := s.CallAPI("wall.post", vals, &n); err != nil {
		return 0, err
	}
	return n.PostId, nil
}

func (s *Session) WallPostEdit(id int, m, a string, ownerOpts ...interface{}) error {
	vals := make(url.Values)
	vals.Set("post_id", strconv.Itoa(id))
	vals.Set("message", m)
	ownerOptions(vals, ownerOpts...)
	if a != "" {
		vals.Set("attachments", a)
	}
	var r json.RawMessage
	if err := s.CallAPI("wall.edit", vals, &r); err != nil {
		return err
	}
	return nil
}

func wallPinDel(s *Session, act string, id int, ownerOpts ...interface{}) error {
	var r json.RawMessage
	vals := make(url.Values)
	vals.Set("post_id", strconv.Itoa(id))
	ownerOptions(vals, ownerOpts...)
	if err := s.CallAPI(act, vals, &r); err != nil {
		return err
	}
	return nil
}

func (s *Session) WallPin(id int, ownerOpts ...interface{}) error {
	return wallPinDel(s, "wall.pin", id, ownerOpts...)
}

func (s *Session) WallUnpin(id int, ownerOpts ...interface{}) error {
	return wallPinDel(s, "wall.unpin", id, ownerOpts...)
}

func (s *Session) WallDelete(id int, ownerOpts ...interface{}) error {
	return wallPinDel(s, "wall.delete", id, ownerOpts...)
}

func (s *Session) WallRestore(id int, ownerOpts ...interface{}) error {
	return wallPinDel(s, "wall.restore", id, ownerOpts...)
}
