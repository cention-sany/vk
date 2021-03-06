package vk

import (
	"encoding/json"
	"io"
)

type (
	// Group contains community information (https://vk.com/dev/fields_groups)
	Group struct {
		// default fields
		Id          int    `json:"id"`
		Name        string `json:"name"`
		ScreenName  string `json:"screen_name"`
		IsClosed    Bool   `json:"is_closed"`
		IsAdmin     Bool   `json:"is_admin"`
		AdminLvl    int    `json:"admin_level"`
		IsMember    Bool   `json:"is_member"`
		Type        string `json:"type"`
		Photo       string `json:"photo"`
		PhotoMedium string `json:"photo_medium"`
		PhotoBig    string `json:"photo_big"`
		// optional fields
		City    int `json:"city,omitempty"`
		Country int `json:"country,omitempty"`
		Place   struct {
			Pid       int     `json:"pid"`
			Title     string  `json:"title"`
			Latitude  float32 `json:"latitude"`
			Longitude float32 `json:"longitude"`
			Type      string  `json:"type"`
			Country   int     `json:"country"`
			City      int     `json:"city"`
			Address   string  `json:"address"`
		} `json:"place,omitempty"`
		Desc           string          `json:"description,omitempty"`
		Wiki           string          `json:"wiki_page,omitempty"`
		MembersCount   int             `json:"members_count,omitempty"`
		Counters       json.RawMessage `json:"counters,omitempty"` // TODO: counters
		StartDate      int64           `json:"start_date,omitempty"`
		EndDate        int64           `json:"end_date,omitempty"`
		CanPost        Bool            `json:"can_post,omitempty"`
		CanSeeAllPosts Bool            `json:"can_see_all_posts,omitempty"`
		Activity       string          `json:"activity,omitempty"`
		Status         string          `json:"status,omitempty"`
		Contacts       string          `json:"contacts,omitempty"`
	}

	SmallGroups struct {
		Count int           `json:"count"`
		Items []*SmallGroup `json:"items"`
	}

	SmallGroup struct {
		Id      int    `json:"id"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		IsAdmin Bool   `json:"is_admin"`
	}
)

type GroupLooper interface {
	More() (*SmallGroup, error)
	Size() int
}

func (s *SmallGroups) Size() int {
	return s.Count
}

func (s *SmallGroups) More() (*SmallGroup, error) {
	if len(s.Items) <= 0 {
		return nil, io.EOF
	}
	sg := s.Items[0]
	s.Items = s.Items[1:]
	return sg, nil
}

func NewGroupLoop(s *SmallGroups) GroupLooper {
	return s
}
