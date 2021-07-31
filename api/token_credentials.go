package api

import model "FiberFinanceAPI/database/models"

// tokenCredentials returns access token, refresh token and when they all expire
func (s *Server) tokenCredentials(user model.User) (accessToken string, accessTokenExp int64, refreshToken string,
	refreshTokenExp int64, err error) {
	s.logs.WithField("func", "token_credentials.go -> tokenCredentials()").Debug()
	// Create Token for the valid user
	accessToken, err = s.token.CreateAccessToken(string(user.ID), s.config.TokenDuration)
	if err != nil {
		s.logs.WithError(err).Warn("cannot create token")
	}
	s.logs.Debug("Server side Access Token Created")
	accessTokenExp, err = s.token.AccessTokenExpiresAt(accessToken)
	if err != nil {
		s.logs.WithError(err).Warn()
	}
	refreshToken, err = s.token.CreateRefreshToken(string(user.ID), s.config.RefreshTokenDuration)
	if err != nil {
		s.logs.WithError(err).Warn()
	}
	s.logs.Debug("Server side Refresh Token Created")
	refreshTokenExp, err = s.token.RefreshTokenExpiresAt(refreshToken)
	if err != nil {
		s.logs.WithError(err).Warn("unable to get Refresh Token expires at")
	}
	s.logs.Debug("Token response successful")
	return
}
