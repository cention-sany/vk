package vk

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type Audio struct {
	Id       int    `json:"id"`
	Owner    int    `json:"owner_id"`
	Artist   string `json:"artist"`
	Title    string `json:"title"`
	Duration int    `json:"duration"` // in sec // TODO: parse to time.Duration
	Url      string `json:"url"`
	Lyrics   int    `json:"lyrics_id"`
	Album    int    `json:"album_id"`
	Genre    int    `json:"genre_id"`
}

type Playlist struct {
	Id    int    `json:"id"`
	Owner int    `json:"owner_id"`
	Title string `json:"title"`
}

func (s *Session) AudioSearch(qu string, count int, autoCompl, perfOnly bool) ([]Audio, error) {
	if qu == "" {
		return nil, errors.New("you must provide a search query")
	}

	vals := make(url.Values)
	vals.Set("q", qu)
	if count > 0 {
		vals.Set("count", fmt.Sprint(count))
	}
	vals.Set("sort", "2")
	if autoCompl {
		vals.Set("auto_complete", "1")
	} else {
		vals.Set("auto_complete", "0")
	}
	if perfOnly {
		vals.Set("performer_only", "1")
	} else {
		vals.Set("performer_only", "0")
	}

	var audio []Audio
	list := ApiList{
		Items: &audio,
	}
	if err := s.CallAPI("audio.search", vals, &list); err != nil {
		return nil, err
	}
	return audio, nil
}

func (s *Session) AudioGetById(ids [][2]int) ([]Audio, error) {
	if len(ids) == 0 {
		return nil, errors.New("you must pass at least one pair of ids")
	}

	var audios = make([]string, len(ids))
	for i, v := range ids {
		audios[i] = fmt.Sprintf("%d_%d", v[0], v[1])
	}

	vals := make(url.Values)
	vals.Set("audios", strings.Join(audios, ","))

	var audio []Audio
	list := ApiList{
		Items: &audio,
	}
	if err := s.CallAPI("audio.getById", vals, &list); err != nil {
		return nil, err
	}
	return audio, nil
}

func (s *Session) AudioGetAlbums(owner int, offset, count int) ([]Playlist, error) {
	vals := make(url.Values)
	if owner != 0 {
		vals.Set("owner_id", fmt.Sprint(owner))
	}
	if offset > 0 {
		vals.Set("offset", fmt.Sprint(offset))
	}
	if count > 0 {
		vals.Set("count", fmt.Sprint(count))
	}

	var plists []Playlist
	list := ApiList{
		Items: &plists,
	}
	if err := s.CallAPI("audio.getAlbums", vals, &list); err != nil {
		return nil, err
	}
	return plists, nil
}

func (s *Session) audioGetFromAny(vals url.Values, offset, count int) ([]Audio, error) {
	if offset > 0 {
		vals.Set("offset", fmt.Sprint(offset))
	}
	if count > 0 {
		vals.Set("count", fmt.Sprint(count))
	}

	var audio []Audio
	list := ApiList{
		Items: &audio,
	}
	if err := s.CallAPI("audio.get", vals, &list); err != nil {
		return nil, err
	}
	return audio, nil
}

func (s *Session) AudioGetFromAlbum(album int, offset, count int) ([]Audio, error) {
	if album <= 0 {
		return nil, errors.New("incorrect album id")
	}

	vals := make(url.Values)
	vals.Set("album_id", fmt.Sprint(album))
	return s.audioGetFromAny(vals, offset, count)
}

func (s *Session) AudioGetFromUser(user int, offset, count int) ([]Audio, error) {
	if user <= 0 {
		return nil, errors.New("incorrect user id")
	}

	vals := make(url.Values)
	vals.Set("owner_id", fmt.Sprint(user))
	return s.audioGetFromAny(vals, offset, count)
}

func (s *Session) AudioGet(ids ...int) ([]Audio, error) {
	if len(ids) == 0 {
		return nil, errors.New("you must pass at least one audio id")
	}

	vals := make(url.Values)
	vals.Set("audio_ids", IdList(ids).String())
	return s.audioGetFromAny(vals, 0, len(ids))
}
