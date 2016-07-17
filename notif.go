package vk

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

const (
	maxNotifReturn         = 100
	NotifFollow            = "follow"
	NotifFriendAccepted    = "friend_accepted"
	NotifMention           = "mention"
	NotifMentionComments   = "mention_comments"
	NotifWall              = "wall"
	NotifCommentPost       = "comment_post"
	NotifCommentPhoto      = "comment_photo"
	NotifCommentVideo      = "comment_video"
	NotifReplyComment      = "reply_comment"
	NotifReplyCommentPhoto = "reply_comment_photo"
	NotifReplyCommentVideo = "reply_comment_video"
	NotifReplyTopic        = "reply_topic"
	NotifLikePost          = "like_post"
	NotifLikeComment       = "like_comment"
	NotifLikePhoto         = "like_photo"
	NotifLikeVideo         = "like_video"
	NotifLikeCommentPhoto  = "like_comment_photo"
	NotifLikeCommentVideo  = "like_comment_video"
	NotifLikeCommentTopic  = "like_comment_topic"
	NotifCopyPost          = "copy_post"
	NotifCopyPhoto         = "copy_photo"
	NotifCopyVideo         = "copy_video"
)

type Notifications struct {
	Count    int          `json:"count"`
	Items    []*NotifItem `json:"items"`
	Profiles []User       `json:"profiles"`
	Groups   []Group      `json:"groups"`
	LastV    int64        `json:"last_viewed"`
	NextF    string       `json:"next_from"`
}

type NotifItem struct {
	Type     string          `json:"type"`
	Date     int64           `json:"date"`
	Parent   json.RawMessage `json:"parent,omitempty"`
	Feedback json.RawMessage `json:"feedback"`
}

// parent decode
func (n *NotifItem) ParentIsPost() bool {
	if n.Type == NotifMentionComments || n.Type == NotifCommentPost ||
		n.Type == NotifLikePost || n.Type == NotifCopyPost {
		return true
	}
	return false
}

func (n *NotifItem) ParentIsComment() bool {
	if n.Type == NotifReplyComment || n.Type == NotifCommentVideo ||
		n.Type == NotifCommentPhoto || n.Type == NotifLikeComment ||
		n.Type == NotifLikeCommentPhoto || n.Type == NotifLikeCommentVideo ||
		n.Type == NotifLikeCommentTopic {
		return true
	}
	return false
}

func (n *NotifItem) ParentAsPost() (*Post, error) {
	v := new(Post)
	err := json.Unmarshal(n.Parent, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (n *NotifItem) ParentAsComment() (*Comment, error) {
	v := new(Comment)
	err := json.Unmarshal(n.Parent, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (n *NotifItem) ParentAsTopic() (*Topic, error) {
	v := new(Topic)
	err := json.Unmarshal(n.Parent, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// feedback decode
func (n *NotifItem) FeedbackIsPost() bool {
	if n.Type == NotifMention || n.Type == NotifWall {
		return true
	}
	return false
}

func (n *NotifItem) FeedbackIsComment() bool {
	if n.Type == NotifMentionComments || n.Type == NotifCommentPost ||
		n.Type == NotifCommentPhoto || n.Type == NotifCommentVideo ||
		n.Type == NotifReplyComment || n.Type == NotifReplyCommentPhoto ||
		n.Type == NotifReplyCommentVideo || n.Type == NotifReplyTopic {
		return true
	}
	return false
}

func (n *NotifItem) FeedbackAsPost() (*Post, error) {
	v := new(Post)
	err := json.Unmarshal(n.Feedback, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (n *NotifItem) FeedbackAsComment() (*Comment, error) {
	v := new(Comment)
	err := json.Unmarshal(n.Feedback, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// NotifGet implements method https://vk.com/dev/notifications.get.
// start and end is start_time and end_time field respectively.
func (s *Session) NotifGet(startFrm string, filters []string, start, end int64) (*Notifications, error) {
	vals := make(url.Values)
	vals.Set("start_from", startFrm)
	vals.Set("filters", strings.Join(filters, ","))
	vals.Set("start_time", strconv.FormatInt(start, 10))
	vals.Set("end_time", strconv.FormatInt(end, 10))

	var n Notifications

	if err := s.CallAPI("notifications.get", vals, &n); err != nil {
		return nil, err
	}
	return &n, nil
}

// NotifMarkAsViewed implements method https://vk.com/dev/notifications.markAsViewed
func (s *Session) NotifMarkAsViewed() (Bool, error) {
	var b Bool
	if err := s.CallAPI("notifications.markAsViewed", url.Values{}, &b); err != nil {
		return Bool(false), err
	}
	return b, nil
}
