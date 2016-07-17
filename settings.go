package vk

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/codeship/go-retro"
)

var serverStatusNotReady = retro.NewStaticRetryableError(
	errors.New("Server status replied: wait"), 10, 3) // 30 seconds the long

// CB for callback
type (
	CBSettings struct {
		MsgNew     Bool `json:"message_new"`
		WCNew      Bool `json:"wall_reply_new"`
		WCEdit     Bool `json:"wall_reply_edit"`
		BPNew      Bool `json:"board_post_new"`
		BPEdit     Bool `json:"board_post_edit"`
		BPDelete   Bool `json:"board_post_delete"`
		BPRestore  Bool `json:"board_post_restore"`
		PhotoNew   Bool `json:"photo_new"`
		VideoNew   Bool `json:"video_new"`
		AudioNew   Bool `json:"audio_new"`
		PCNew      Bool `json:"photo_comment_new"`
		VCNew      Bool `json:"video_comment_new"`
		MCNew      Bool `json:"market_comment_new"`
		GroupJoin  Bool `json:"group_join"`
		GroupLeave Bool `json:"group_leave"`
		WPNew      Bool `json:"wall_post_new"`
	}

	CBServerSettings struct {
		Url string `json:"server_url"`
		Key string `json:"secret_key"`
	}

	CBConfirmCode struct {
		Code string `json:"code"`
	}
)

// SetServer will set the new URL for callback server. And automatically retry
// till it get status ok or failed. This is blocking call.
func SetServer(s *Session, surl string, gid int) error {
	v := url.Values{}
	v.Set("server_url", surl)
	v.Set("group_id", strconv.Itoa(gid))
	return retro.DoWithRetry(func() error {
		var r struct {
			Code  int    `json:"state_code"`
			State string `json:"state"`
		}
		err := s.CallAPI("groups.setCallbackServer", v, &r)
		if err != nil {
			return err
		}
		if r.Code == 3 || r.Code == 4 {
			return errors.New(fmt.Sprint("Set server error: ", r.State,
				" code: ", r.Code))
		} else if r.Code == 2 {
			return serverStatusNotReady
		}
		return nil
	})
}
