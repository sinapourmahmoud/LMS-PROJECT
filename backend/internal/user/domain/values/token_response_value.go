package values

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewTokenResponse(accessToken string, refreshToken string) *TokenResponse {
	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
