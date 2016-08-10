package vk

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	NameCaseNom = "nom"
	NameCaseGen = "gen"
	NameCaseDat = "dat"
	NameCaseAcc = "acc"
	NameCaseIns = "ins"
	NameCaseAbl = "abl"
)

var (
	// NameCases is a list of name cases available for VK
	NameCases = []string{NameCaseNom, NameCaseGen, NameCaseDat, NameCaseAcc, NameCaseIns, NameCaseAbl}
)

const (
	FieldPhotoId         = "photo_id"
	FieldVerified        = "verified"
	FieldBlacklisted     = "blacklisted"
	FieldSex             = "sex"
	FieldBirthDate       = "bdate"
	FieldCity            = "city"
	FieldCountry         = "country"
	FieldHomeTown        = "home_town"
	FieldPhoto50         = "photo_50"
	FieldPhoto100        = "photo_100"
	FieldPhoto200        = "photo_200"
	FieldPhoto200Orig    = "photo_200_orig"
	FieldPhoto400Orig    = "photo_400_orig"
	FieldPhotoMax        = "photo_max"
	FieldPhotoMaxOrig    = "photo_max_orig"
	FieldOnline          = "online"
	FieldLists           = "lists"
	FieldDomain          = "domain"
	FieldHasMobile       = "has_mobile"
	FieldContacts        = "contacts"
	FieldSite            = "site"
	FieldEducation       = "education"
	FieldUniversities    = "universities"
	FieldSchools         = "schools"
	FieldStatus          = "status"
	FieldLastSeen        = "last_seen"
	FieldFollowersCount  = "followers_count"
	FieldCommonCount     = "common_count"
	FieldCounters        = "counters"
	FieldOccupation      = "occupation"
	FieldNickName        = "nickname"
	FieldRelatives       = "relatives"
	FieldRelation        = "relation"
	FieldPersonal        = "personal"
	FieldConnections     = "connections"
	FieldExports         = "exports"
	FieldWallComments    = "wall_comments"
	FieldActivities      = "activities"
	FieldInterests       = "interests"
	FieldMusic           = "music"
	FieldMovies          = "movies"
	FieldTv              = "tv"
	FieldBooks           = "books"
	FieldGames           = "games"
	FieldAbout           = "about"
	FieldQuotes          = "quotes"
	FieldCanPost         = "can_post"
	FieldCanSeeAllPosts  = "can_see_all_posts"
	FieldCanSeeAudio     = "can_see_audio"
	FieldCanWritePrivate = "can_write_private_message"
	FieldTimeZone        = "timezone"
	FieldScreenName      = "screen_name"
	FieldMaidenName      = "maiden_name"
)

const (
	SexUnknown = 0
	SexFemale  = 1
	SexMale    = 2
)

type (
	// User contains user information (https://vk.com/dev/fields)
	User struct {
		Id          int    `json:"id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Deactivated string `json:"deactivated"`
		Hidden      Bool   `json:"hidden"`

		PhotoId      string   `json:"photo_id,omitempty"`
		Verified     Bool     `json:"verified,omitempty"`
		Blacklisted  Bool     `json:"blacklisted,omitempty"`
		Sex          int      `json:"sex,omitempty"`
		Birthdate    string   `json:"bdate,omitempty"`
		City         GeoPlace `json:"city,omitempty"`
		Country      GeoPlace `json:"country,omitempty"`
		HomeTown     string   `json:"home_town,omitempty"`
		Photo50      string   `json:"photo_50,omitempty"`
		Photo100     string   `json:"photo_100,omitempty"`
		Photo200     string   `json:"photo_200,omitempty"`
		Photo200Orig string   `json:"photo_200_orig,omitempty"`
		Photo400Orig string   `json:"photo_400_orig,omitempty"`
		PhotoMax     string   `json:"photo_max,omitempty"`
		PhotoMaxOrig string   `json:"photo_max_orig,omitempty"`
		Online       Bool     `json:"online,omitempty"`
		OnlineMobile Bool     `json:"online_mobile,omitempty"`
		OnlineApp    Bool     `json:"online_app,omitempty"`
		Lists        []int    `json:"lists,omitempty"`
		Domain       string   `json:"domain,omitempty"`
		HasMobile    Bool     `json:"has_mobile,omitempty"`
		Contacts     struct {
			Mobile string `json:"mobile_phone,omitempty"`
			Home   string `json:"home_phone,omitempty"`
		} `json:"contacts,omitempty"`
		Site      string `json:"site,omitempty"`
		Education struct {
			University     int    `json:"university,omitempty"`
			UniversityName string `json:"university_name,omitempty"`
			Faculty        int    `json:"faculty,omitempty"`
			FacultyName    string `json:"faculty_name,omitempty"`
			Graduation     int    `json:"graduation,omitempty"`
		} `json:"education,omitempty"`
		Universities   []University `json:"universities,omitempty"`
		Schools        []School     `json:"schools,omitempty"`
		Status         string       `json:"status,omitempty"`
		StatusAudio    interface{}  `json:"status_audio,omitempty"` // TODO: status_audio
		LastSeen       PlatformInfo `json:"last_seen,omitempty"`
		FollowersCount int          `json:"followers_count,omitempty"`
		CommonCount    int          `json:"common_count,omitempty"`
		Counters       struct {
			Albums        int `json:"albums,omitempty"`
			Videos        int `json:"videos,omitempty"`
			Audios        int `json:"audios,omitempty"`
			Photos        int `json:"photos,omitempty"`
			Notes         int `json:"notes,omitempty"`
			Friends       int `json:"friends,omitempty"`
			Groups        int `json:"groups,omitempty"`
			OnlineFriends int `json:"online_friends,omitempty"`
			MutualFriends int `json:"mutual_friends,omitempty"`
			UserVideos    int `json:"user_videos,omitempty"`
			Followers     int `json:"followers,omitempty"`
			UserPhotos    int `json:"user_photos,omitempty"`
			Subscriptions int `json:"subscriptions,omitempty"`
		} `json:"counters,omitempty"`
		Occupation struct {
			Type string `json:"type,omitempty"`
			Id   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"occupation,omitempty"`
		NickName  string     `json:"nickname"`
		Relatives []Relative `json:"relatives,omitempty"`
		Relation  int        `json:"relation,omitempty"` // TODO: constants for relation
		Personal  struct {   // TODO: constants for personal info
			Political  int      `json:"political,omitempty"`
			Langs      []string `json:"langs,omitempty"`
			Religion   string   `json:"religion,omitempty"`
			InspiredBy string   `json:"inspired_by,omitempty"`
			PeopleMain int      `json:"people_main,omitempty"`
			LifeMain   int      `json:"life_main,omitempty"`
			Smoking    int      `json:"smoking,omitempty"`
			Alcohol    int      `json:"alcohol,omitempty"`
		} `json:"personal,omitempty"`
		Connections     interface{} `json:"connections,omitempty"` // TODO: connections
		Exports         interface{} `json:"exports,omitempty"`     // TODO: exports
		WallComments    Bool        `json:"wall_comments,omitempty"`
		Activities      string      `json:"activities,omitempty"`
		Interests       string      `json:"interests,omitempty"`
		Music           string      `json:"music,omitempty"`
		Movies          string      `json:"movies,omitempty"`
		Tv              string      `json:"tv,omitempty"`
		Books           string      `json:"books,omitempty"`
		Games           string      `json:"games,omitempty"`
		About           string      `json:"about,omitempty"`
		Quotes          string      `json:"quotes,omitempty"`
		CanPost         Bool        `json:"can_post,omitempty"`
		CanSeeAllPosts  Bool        `json:"can_see_all_posts,omitempty"`
		CanSeeAudio     Bool        `json:"can_see_audio,omitempty"`
		CanWritePrivate Bool        `json:"can_write_private_message,omitempty"`
		TimeZone        int         `json:"timezone,omitempty"`
		ScreenName      string      `json:"screen_name,omitempty"`
		MaidenName      string      `json:"maiden_name,omitempty"`
	}
	SmallUser struct {
		Id         int    `json:"id"`
		ScreenName string `json:"screen_name"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Timezone   int    `json:"timezone"`
	}
	SmallUsers []SmallUser
	// GeoPlace contains geographical information like City, Country
	GeoPlace struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
	}
	// PlatformInfo contains information about time and platform
	PlatformInfo struct {
		Time     EpochTime `json:"time"`
		Platform int       `json:"platform"`
	}
	// University contains information about the university
	University struct {
		Id              int    `json:"id"`
		Country         int    `json:"country"`
		City            int    `json:"city"`
		Name            string `json:"name"`
		Faculty         int    `json:"faculty"`
		FacultyName     string `json:"faculty_name"`
		Chair           int    `json:"chair"`
		ChairName       string `json:"chair_name"`
		Graduation      int    `json:"graduation"`
		EducationForm   string `json:"education_form"`
		EducationStatus string `json:"education_status"`
	}
	// School contains information about schools
	School struct {
		Id         int    `json:"id"`
		Country    int    `json:"country"`
		City       int    `json:"city"`
		Name       string `json:"name"`
		YearFrom   int    `json:"year_from"`
		YearTo     int    `json:"year_to"`
		Class      string `json:"class"`
		TypeStr    string `json:"type_str,omitempty"`
		Speciality string `json:"speciality,omitempty"`
	}
	// Relative contains information about relative to the user
	Relative struct {
		Id   int    `json:"id"`   // negative id describes non-existing users (possibly prepared id if they will register)
		Type string `json:"type"` // like `parent`, `grandparent`, `sibling`
		Name string `json:"name,omitempty"`
	}
)

// UsersGet implements method http://vk.com/dev/users.get
func (s *Session) UsersGet(userIds []int, fields []string, nameCase string) ([]User, error) {
	if len(userIds) == 0 {
		return nil, errors.New("you must pass at least one id or screen_name")
	}
	if nameCase == "" {
		nameCase = NameCaseNom
	}
	if !ElemInSlice(nameCase, NameCases) {
		return nil, errors.New("the only available name cases are: " + strings.Join(NameCases, ", "))
	}

	vals := make(url.Values)
	vals.Set("user_ids", IdList(userIds).String())
	vals.Set("fields", strings.Join(fields, ","))
	vals.Set("name_case", nameCase)

	var users []User

	if err := s.CallAPI("users.get", vals, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// User returns current user info with call to UsersGet
func (s *Session) User(fields []string, nameCase string) (User, error) {
	var u User
	list, err := s.UsersGet([]int{s.UserID}, fields, nameCase)
	if err != nil {
		return u, err
	}
	if len(list) == 0 {
		return u, errors.New("empty response")
	}
	u = list[0]
	return u, nil
}

// UsersIsAppUser implements https://vk.com/dev/users.isAppUser
func (s *Session) UsersIsAppUser(user int) (bool, error) {
	if user < 0 {
		return false, errors.New("incorrect user id")
	}

	vals := make(url.Values)
	if user > 0 {
		vals.Set("user_id", fmt.Sprint(user))
	}

	var res Bool
	if err := s.CallAPI("users.isAppUser", vals, &res); err != nil {
		return false, err
	}
	return bool(res), nil
}

// UsersGetFollowers implements https://vk.com/dev/users.getFollowers
func (s *Session) UsersGetFollowers(user int, fields []string, nameCase string, offset, count int) ([]User, error) {
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

	if err := s.CallAPI("users.getFollowers", vals, &list); err != nil {
		return nil, err
	}
	return users, nil
}
