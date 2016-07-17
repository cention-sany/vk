package vk

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	JT_Join     = "join"
	JT_Unsure   = "unsure"
	JT_Accepted = "accepted"
	JT_Approved = "approved"
	JT_Request  = "request"
)

type (
	// ReceivedResult type
	ReceivedResult struct {
		Type    string          `json:"type"`
		Object  json.RawMessage `json:"object"`
		GroupId int             `json:"group_id"`
		Secret  string          `json:"secret"`
	}

	TopicDelete struct {
		Topic int `json:"topic_id"`
		Id    int `json:"id"`
	}

	GroupLeave struct {
		UserId string `json:"user_id"`
		Self   Bool   `json:"self"`
	}

	GroupJoin struct {
		UserId   string `json:"user_id"`
		JoinType string `json:"join_type"`
	}

	// https://new.vk.com/dev/callback_api
	TopicComment struct {
		Id           int           `json:"id"`
		FromId       int           `json:"from_id"`
		Date         int64         `json:"date"`
		Text         string        `json:"text"`
		TopicOwnerId int           `json:"topic_owner_id"`
		TopicId      int           `json:"topic_id"`
		Attachments  []*Attachment `json:"attachments"`
		Likes        struct {
			Count     int  `json:"count"`
			UserLikes Bool `json:"user_likes"`
			CanLike   Bool `json:"can_like"`
		} `json:"likes"`
	}

	Receive struct {
		secret string
	}
)

// ParseRequest function
func (rx *Receive) ParseRequest(r *http.Request) (res *ReceivedResult, err error) {
	defer r.Body.Close()
	res = &ReceivedResult{}
	unmarshal := json.NewDecoder(r.Body)
	if err = unmarshal.Decode(res); err != nil {
		return nil, err
	}
	if rx.secret != "" && res.Secret != rx.secret {
		return nil, errors.New("vk: invalid callback API secret")
	}
	return
}

func ParseRequest(r *http.Request) (res *ReceivedResult, err error) {
	defer r.Body.Close()
	res = &ReceivedResult{}
	unmarshal := json.NewDecoder(r.Body)
	if err = unmarshal.Decode(res); err != nil {
		return nil, err
	}
	return
}

func (rr *ReceivedResult) WallComment() (*Comment, error) {
	v := &Comment{}
	if err := unmarshaler(v, bytes.NewReader(rr.Object)); err != nil {
		return nil, err
	}
	return v, nil
}

func (rr *ReceivedResult) WallPost() (*Post, error) {
	v := &Post{}
	if err := unmarshaler(v, bytes.NewReader(rr.Object)); err != nil {
		return nil, err
	}
	return v, nil
}

func (rr *ReceivedResult) Audio() (*Audio, error) {
	v := &Audio{}
	if err := unmarshaler(v, bytes.NewReader(rr.Object)); err != nil {
		return nil, err
	}
	return v, nil
}

func (rr *ReceivedResult) Photo() (*Photo, error) {
	v := &Photo{}
	if err := unmarshaler(v, bytes.NewReader(rr.Object)); err != nil {
		return nil, err
	}
	return v, nil
}

func (rr *ReceivedResult) Video() (*Video, error) {
	v := &Video{}
	if err := unmarshaler(v, bytes.NewReader(rr.Object)); err != nil {
		return nil, err
	}
	return v, nil
}

func (rr *ReceivedResult) PM() (*Message, error) {
	v := &Message{}
	if err := unmarshaler(v, bytes.NewReader(rr.Object)); err != nil {
		return nil, err
	}
	return v, nil
}

func (rr *ReceivedResult) GetGroupLeave() (*GroupLeave, error) {
	v := &GroupLeave{}
	if err := unmarshaler(v, bytes.NewReader(rr.Object)); err != nil {
		return nil, err
	}
	return v, nil
}

func (rr *ReceivedResult) GetGroupJoin() (*GroupJoin, error) {
	v := &GroupJoin{}
	if err := unmarshaler(v, bytes.NewReader(rr.Object)); err != nil {
		return nil, err
	}
	return v, nil
}

func (rr *ReceivedResult) Board() (*TopicComment, error) {
	v := &TopicComment{}
	if err := unmarshaler(v, bytes.NewReader(rr.Object)); err != nil {
		return nil, err
	}
	return v, nil
}

func (rr *ReceivedResult) GetBoardDelete() (*TopicDelete, error) {
	v := &TopicDelete{}
	if err := unmarshaler(v, bytes.NewReader(rr.Object)); err != nil {
		return nil, err
	}
	return v, nil
}
