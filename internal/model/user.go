package model

const (
	UserAuthorityUser  = 0
	UserAuthorityAdmin = 255
)

type User struct {
	Common
	GithubID  int64  `json:"github_id,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
	Sid       string `json:"sid,omitempty"`
	Authority uint8  `json:"authority,omitempty"`
}
