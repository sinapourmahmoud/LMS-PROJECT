package values

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func NewAccessTokenResponse(accessToken string) *AccessTokenResponse {

	return &AccessTokenResponse{
		AccessToken: accessToken,
	}
}
