package vk

import (
	"net/url"
	"strconv"
)

type LikeType int

const (
	likesAdd = "likes.add"
	likesDel = "likes.delete"
)

const (
	LikesPost LikeType = iota
	LikesComment
	LikesPhoto
	LikesAudio
	LikesVideo
	LikesNote
	LikesPhotoComment
	LikesVideoComment
	LikesTopicComment
)

var likeType = []string{
	"post",
	"comment",
	"photo",
	"audio",
	"video",
	"note",
	"photo_comment",
	"video_comment",
	"topic_comment",
}

func (l LikeType) String() string {
	return likeType[int(l)]
}

func likeUnlike(s *Session, act string, t LikeType, id int,
	likesOptions ...interface{}) (int, error) {
	vals := url.Values{}
	vals.Set("type", t.String())
	vals.Set("item_id", strconv.Itoa(id))
	if likesOptions != nil {
		optSize := len(likesOptions)
		if optSize > 1 {
			if idInt, ok := likesOptions[0].(int); ok && idInt > 0 {
				if isCommunity, ok := likesOptions[1].(bool); ok {
					if isCommunity {
						idInt = -1 * idInt
					}
				}
				vals.Set("owner_id", strconv.Itoa(idInt))
			}
		}
		if optSize > 2 {
			if ak, ok := likesOptions[2].(string); ok && ak != "" {
				vals.Set("access_key", ak)
			}
		}
	}
	var n struct {
		Likes int `json:"likes"`
	}
	if err := s.CallAPI(act, vals, &n); err != nil {
		return 0, err
	}
	return n.Likes, nil
}

func (s *Session) Likes(t LikeType, id int, likesOptions ...interface{}) (int, error) {
	return likeUnlike(s, likesAdd, t, id, likesOptions...)
}

func (s *Session) Unlikes(t LikeType, id int, likesOptions ...interface{}) (int, error) {
	return likeUnlike(s, likesDel, t, id, likesOptions...)
}
