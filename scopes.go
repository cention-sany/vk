package vk

const (
	audStr = "audio"
	vidStr = "video"
)

// Scope is an access scope from https://vk.com/dev/permissions
type Scope int

func (s Scope) String() string {
	switch s {
	case 1:
		return "notify"
	case 2:
		return "friends"
	case 4:
		return "photos"
	case 8:
		return audStr
	case 16:
		return vidStr
	case 32: // deprecated
		return "offers"
	case 64: // deprecated
		return "questions"
	case 128:
		return "pages"
	case 1024:
		return "status"
	case 2048:
		return "notes"
	case 4096:
		return "messages"
	case 8192:
		return "wall"
	case 32768:
		return "ads"
	case 65536:
		return "offline"
	case 131072:
		return "docs"
	case 262144:
		return "groups"
	case 524288:
		return "notifications"
	case 1048576:
		return "stats"
	case 4194304:
		return "email"
	default:
		return ""
	}
}

// List of known access scopes from https://vk.com/dev/permissions
const (
	ScopeNotify        = Scope(1)
	ScopeFriends       = Scope(2)
	ScopePhotos        = Scope(4)
	ScopeAudio         = Scope(8)
	ScopeVideo         = Scope(16)
	ScopeDocs          = Scope(131072)
	ScopeNotes         = Scope(2048)
	ScopePages         = Scope(128)
	ScopeStatus        = Scope(1024)
	ScopeOffers        = Scope(32)
	ScopeQuestions     = Scope(64)
	ScopeWall          = Scope(8192)
	ScopeGroups        = Scope(262144)
	ScopeMessages      = Scope(4096)
	ScopeEmail         = Scope(4194304)
	ScopeNotifications = Scope(524288)
	ScopeStats         = Scope(1048576)
	ScopeAds           = Scope(32768)
	ScopeOffline       = Scope(65536)
)
