package handler

import (
	"backend/internal/handler/auth"
	"backend/internal/svc"
)

func RegisterHandlers(s *svc.ServiceContext) {
	// === Authentication Routes ===
	authz := s.Echo.Group("/auth")
	authz.POST("/signup", auth.SignUp(s))
	authz.POST("/signin", auth.SignIn(s))
	authz.POST("/password-forgot", auth.ForgotPassword(s))
	authz.POST("/reset-password", auth.ResetPassword(s))
	authz.POST("/verify", auth.VerifyEmail(s))
	authz.POST("/refresh-token", auth.RefreshToken(s))
}
