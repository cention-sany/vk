package vk

import (
	"net/url"
	"strconv"
	"strings"
)

const (
	MsgUnreadOnly  = 1
	MsgNotChat     = 2
	MsgFromFriends = 4
)

// https://vk.com/dev/message
type (
	Messages struct {
		Count int
		Items []*Message
	}

	Message struct {
		Id          int           `json:"id"`
		UserId      int           `json:"user_id"`
		Date        int64         `json:"date"`
		ReadState   Bool          `json:"read_state"`
		Out         Bool          `json:"out"`
		Title       string        `json:"title"`
		Body        string        `json:"body"`
		Attachments []*Attachment `json:"attachments"`
		Emoji       Bool          `json:"emoji"`
		Deleted     Bool          `json:"deleted"`
		Fwd         []*Message    `json:"fwd_messages"`
		G           *Geo          `json:"geo"`
		// for group chat
		ChatId     int    `json:"chat_id"`
		ChatActive []int  `json:"chat_active"`
		UsersCount int    `json:"users_count"`
		AdminId    int    `json:"admin_id"`
		Photo50    string `json:"photo_50"`
		Photo100   string `json:"photo_100"`
		Photo200   string `json:"photo_200"`
	}
)

// MsgsGet implements method https://vk.com/dev/messages.get
func (s *Session) MsgsGet(out bool, filters int, lastId int) (*Messages, error) {
	outStr := "0"
	if out {
		outStr = "1"
	}
	vals := make(url.Values)
	vals.Set("out", outStr)
	vals.Set("filters", strconv.Itoa(filters))
	vals.Set("count", "200")
	vals.Set("last_message_id", strconv.Itoa(lastId))

	var m Messages
	if err := s.CallAPI("messages.get", vals, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// MsgsMarkAsRead implements method https://vk.com/dev/messages.markAsRead
func (s *Session) MsgsMarkAsRead(ids []int) (Bool, error) {
	var b Bool
	ss := make([]string, len(ids))
	for i, _ := range ids {
		ss[i] = strconv.Itoa(ids[i])
	}
	v := url.Values{}
	v.Set("message_ids", strings.Join(ss, ","))
	if err := s.CallAPI("messages.markAsRead", v, &b); err != nil {
		return Bool(false), err
	}
	return b, nil
}
