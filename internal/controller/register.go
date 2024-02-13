package controller

import (
	"root-histoty-service/internal/middleware"
)

func (s *Server) RegisterRoutes() {
	s.r.Use(middleware.Logger(s.logger))
	s.r.POST("/user/register", s.Register)
	s.r.GET("/user/authorize", s.Authorize)
	s.r.GET("/user/refresh-tokens", s.RefreshToken)

	e := s.r.Group("/user")
	e.Use(middleware.Auth(s.cfg.SecretWord, s.logger))
	e.GET("/info", s.GetUserInfo)

	//s.r.Use(middleware.Auth(s.cfg.SecretWord, s.logger))
}
