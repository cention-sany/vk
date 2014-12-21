package vk

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func (s *Session) FriendsGet(user int, offset, count int, fields []string, nameCase string) ([]User, error) {
	if user < 0 {
		return nil, errors.New("incorrect user id")
	}
	if nameCase == "" {
		nameCase = NameCaseNom
	}
	if !ElemInSlice(nameCase, NameCases) {
		return nil, errors.New("the only available name cases are: " + strings.Join(NameCases, ", "))
	}

	vals := make(url.Values)
	if user > 0 {
		vals.Set("user_id", fmt.Sprint(user))
	}
	vals.Set("fields", strings.Join(fields, ","))
	vals.Set("order", "name")
	vals.Set("name_case", nameCase)
	if offset > 0 {
		vals.Set("offset", fmt.Sprint(offset))
	}
	if count > 0 {
		vals.Set("count", fmt.Sprint(count))
	}

	var users []User
	list := ApiList{
		Items: &users,
	}
	if err := s.CallAPI("friends.get", vals, &list); err != nil {
		return nil, err
	}
	return users, nil
}
