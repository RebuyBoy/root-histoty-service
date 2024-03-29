package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"root-histoty-service/config"
	"root-histoty-service/internal"
	"root-histoty-service/internal/converter"
	"root-histoty-service/internal/dto/request"
	"root-histoty-service/internal/middleware"
	"strconv"
	"time"
)

type Server struct {
	cfg    *config.ServerConfig
	logger *logrus.Logger
	r      *echo.Echo
	p      internal.PlayerService
}

type tokenResponse struct {
	Message        string `json:"message"`
	AccessTokenTTL int64  `json:"access_token_ttl"`
}

func NewServer(cfg *config.ServerConfig, logger *logrus.Logger, p internal.PlayerService) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
		r:      echo.New(),
		p:      p,
	}
}

func (s *Server) StartRouter() {
	srv := http.Server{
		Addr:    "localhost:" + s.cfg.Port,
		Handler: middleware.Cors(s.r),
	}
	s.logger.Info("server is running....")
	err := srv.ListenAndServe()
	if err != nil {
		s.logger.Fatal(err)
	}
}

func (s *Server) Register(ctx echo.Context) error {

	data, err := s.parsePlayerData(ctx)
	if err != nil {
		s.logger.Error("error parsing context during creating player: ", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = s.p.Register(
		ctx.Request().Context(),
		converter.CreatePlayerRequestToPlayer(data),
	)

	if err != nil {
		s.logger.Error("error parsing context during creating player: ", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return ctx.JSON(http.StatusCreated, "player created")
}

func (s *Server) Authorize(ctx echo.Context) error {
	login := ctx.QueryParam("name")
	pinCodeStr := ctx.QueryParam("pin-code")
	pinCode, err := s.parsePinCode(pinCodeStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid pin code. Please provide a valid numeric pin code")
	}

	tokens, err := s.p.Authorize(ctx.Request().Context(), login, pinCode)
	ctx.SetCookie(&http.Cookie{
		Name:     "refresh-token",
		Value:    tokens.RefreshToken,
		Path:     "/user/refresh-tokens",
		HttpOnly: true,
		MaxAge:   180,
	})

	logrus.Info("tokens.AccessTokenExpireAt: ", tokens.AccessTokenExpireAt)
	logrus.Info(time.Unix(tokens.AccessTokenExpireAt, 0))
	ctx.SetCookie(&http.Cookie{
		Name:   "access-token",
		Value:  tokens.AccessToken,
		Path:   "/",
		MaxAge: 60,
	})

	return ctx.JSON(http.StatusOK, tokenResponse{
		Message:        "authorized",
		AccessTokenTTL: tokens.AccessTokenExpireAt,
	})
}

func (s *Server) RefreshToken(ctx echo.Context) error {
	refreshToken, err := ctx.Cookie("refresh-token")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "refresh token not found")
	}
	tokens, err := s.p.RefreshTokens(ctx.Request().Context(), refreshToken.Value)
	ctx.SetCookie(&http.Cookie{
		Name:  "refresh-token",
		Value: tokens.RefreshToken,
		Path:  "/user/refresh-token",
	})
	ctx.SetCookie(&http.Cookie{
		Name:   "access-token",
		Value:  tokens.AccessToken,
		Path:   "/",
		MaxAge: 60,
	})
	return ctx.JSON(http.StatusOK, "token refreshed")
}

func (s *Server) GetUserInfo(ctx echo.Context) error {
	userId := ctx.Get("userId").(string)
	return ctx.JSON(http.StatusOK, userId)

}

func (s *Server) parsePlayerData(ctx echo.Context) (*request.CreatePlayerRequest, error) {
	playerName := ctx.FormValue("name")
	pinCodeStr := ctx.FormValue("pinCode")
	pinCode, err := s.parsePinCode(pinCodeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid pin code. Please provide a valid numeric pin code")
	}
	avatarFileHeader, err := ctx.FormFile("avatar")
	if err != nil {
		return nil, fmt.Errorf("Could not get avatar file from form data: " + err.Error())
	}
	avatarBytes, err := readAvatarFile(avatarFileHeader)
	if err != nil {
		return nil, fmt.Errorf("Could not read avatar file: " + err.Error())
	}

	return &request.CreatePlayerRequest{
		Name:    playerName,
		PinCode: pinCode,
		Avatar:  avatarBytes,
	}, nil
}

func (s *Server) parsePinCode(pinCodeStr string) (int, error) {
	pinCode, err := strconv.Atoi(pinCodeStr)
	if err != nil {
		s.logger.Error("error parsing pin code: ", err)
		return -1, fmt.Errorf("Could not parse pin code: " + err.Error())
	}
	return pinCode, nil
}

func readAvatarFile(avatarFileHeader *multipart.FileHeader) ([]byte, error) {
	avatar, err := avatarFileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("Could not open avatar file: " + err.Error())
	}
	defer avatar.Close()

	avatarBytes, err := io.ReadAll(avatar)
	if err != nil {
		return nil, fmt.Errorf("Could not read avatar file: " + err.Error())
	}

	return avatarBytes, nil
}
