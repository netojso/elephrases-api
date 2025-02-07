package domain

type Login struct {
	Email    string
	Password string
}

type Register struct {
	Email    string
	Password string
}

type Session struct {
	AccessToken  string
	RefreshToken string
}

func (s *Session) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"access_token":  s.AccessToken,
		"refresh_token": s.RefreshToken,
	}
}
