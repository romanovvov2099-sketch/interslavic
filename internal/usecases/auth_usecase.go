package usecases

import (
	"context"
	"errors"
	"interslavic/internal/auth"
	"interslavic/internal/database"
	"interslavic/internal/models"
	"interslavic/logging"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo database.UserRepository
	jwtCfg   *auth.JWTConfig
	logger   *logging.ModuleLogger
}

func NewAuthUseCase(
	userRepo database.UserRepository,
	jwtCfg *auth.JWTConfig,
	logger *logging.ModuleLogger,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
		jwtCfg:   jwtCfg,
		logger:   logging.NewModuleLogger("AUTH", "USECASE", logger),
	}
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type AuthResponse struct {
	User      *models.User    `json:"user"`
	TokenPair *auth.TokenPair `json:"tokens"`
}

func (uc *AuthUseCase) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	user, err := uc.userRepo.FindByLogin(ctx, req.Login)
	if err != nil {
		uc.logger.Error("user not found", logging.ErrAttr(err))
		return nil, errors.New("invalid user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		uc.logger.Error("invalid password", logging.ErrAttr(err))
		return nil, errors.New("invalid password")
	}

	tokenPair, err := uc.jwtCfg.GenerateTokenPair(user.ID, user.Login, user.Role)
	if err != nil {
		uc.logger.Error("failed to generate tokens", logging.ErrAttr(err))
		return nil, errors.New("internal server error")
	}

	if err := uc.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		uc.logger.Error("failed to update last login", logging.ErrAttr(err))
	}

	user.Password = ""

	return &AuthResponse{
		User:      user,
		TokenPair: tokenPair,
	}, nil
}

func (uc *AuthUseCase) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.Error("failed to hash password", logging.ErrAttr(err))
		return nil, errors.New("internal server error")
	}

	role := req.Role
	if role == "" {
		role = "student"
	}

	user := &models.User{
		Fullname: req.Fullname,
		Email:    req.Email,
		Login:    req.Login,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		uc.logger.Error("failed to create user", logging.ErrAttr(err))
		return nil, errors.New("user with this login or email already exists")
	}

	tokenPair, err := uc.jwtCfg.GenerateTokenPair(user.ID, user.Login, user.Role)
	if err != nil {
		uc.logger.Error("failed to generate tokens", logging.ErrAttr(err))
		return nil, errors.New("internal server error")
	}

	user.Password = ""

	return &AuthResponse{
		User:      user,
		TokenPair: tokenPair,
	}, nil
}

func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenPair, error) {
	claims, err := uc.jwtCfg.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := uc.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	tokenPair, err := uc.jwtCfg.GenerateTokenPair(user.ID, user.Login, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	return tokenPair, nil
}

func (uc *AuthUseCase) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
