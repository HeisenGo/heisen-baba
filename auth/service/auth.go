package service

import (
	"authservice/internal/user"
	"authservice/pkg/jwt"
	"context"
	"time"

	jwt2 "github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userOps                *user.Ops
	secret                 []byte
	tokenExpiration        uint
	refreshTokenExpiration uint
}

func NewAuthService(userOps *user.Ops, secret []byte,
	tokenExpiration uint, refreshTokenExpiration uint) *AuthService {
	return &AuthService{
		userOps:                userOps,
		secret:                 secret,
		tokenExpiration:        tokenExpiration,
		refreshTokenExpiration: refreshTokenExpiration,
	}
}

type UserToken struct {
	AuthorizationToken string
	RefreshToken       string
	ExpiresAt          int64
}

func (s *AuthService) CreateUser(ctx context.Context, user *user.User) (*user.User, error) {
	return s.userOps.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, email, pass string) (*UserToken, error) {
	fetchedUser, err := s.userOps.GetUserByEmailAndPassword(ctx, email, pass)
	if err != nil {
		return nil, err
	}

	// calc expiration time values
	var (
		authExp    = time.Now().Add(time.Minute * time.Duration(s.tokenExpiration))
		refreshExp = time.Now().Add(time.Minute * time.Duration(s.refreshTokenExpiration))
	)

	authToken, err := jwt.CreateToken(s.secret, s.userClaims(fetchedUser, authExp))
	if err != nil {
		return nil, err // todo
	}

	refreshToken, err := jwt.CreateToken(s.secret, s.userClaims(fetchedUser, refreshExp))
	if err != nil {
		return nil, err // todo
	}

	return &UserToken{
		AuthorizationToken: authToken,
		RefreshToken:       refreshToken,
		ExpiresAt:          authExp.Unix(),
	}, nil
}

func (s *AuthService) RefreshAuth(ctx context.Context, refreshToken string) (*UserToken, error) {
	claim, err := jwt.ParseToken(refreshToken, s.secret)
	if err != nil {
		return nil, err
	}

	u, err := s.userOps.GetUserByID(ctx, claim.UserID)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, user.ErrUserNotFound
	}

	// calc expiration time values
	var (
		authExp = time.Now().Add(time.Minute * time.Duration(s.tokenExpiration))
	)

	authToken, err := jwt.CreateToken(s.secret, s.userClaims(u, authExp))
	if err != nil {
		return nil, err // todo
	}

	return &UserToken{
		AuthorizationToken: authToken,
		RefreshToken:       refreshToken,
		ExpiresAt:          authExp.UnixMilli(),
	}, nil
}

func (s *AuthService) userClaims(user *user.User, exp time.Time) *jwt.UserClaims {
	return &jwt.UserClaims{
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: &jwt2.NumericDate{
				Time: exp,
			},
		},
		UserID: user.ID,
		Role:   user.Role.String(),
	}
}

package user

import (
    "context"
    "errors"

    "github.com/your-repo/auth/pkg/adapters/storage/entities"
    "github.com/your-repo/auth/pkg/adapters/storage"
    "golang.org/x/crypto/bcrypt"
)

type Service struct {
    storage *storage.Storage
}

func NewService(storage *storage.Storage) *Service {
    return &Service{storage: storage}
}

func (s *Service) CreateUser(ctx context.Context, username, email, password string) (*entities.User, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &entities.User{
        Username:     username,
        Email:        email,
        PasswordHash: string(hashedPassword),
    }

    err = s.storage.CreateUser(ctx, user)
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (s *Service) AuthenticateUser(ctx context.Context, username, password string) (*entities.User, error) {
    user, err := s.storage.GetUserByUsername(ctx, username)
    if err != nil {
        return nil, err
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    return user, nil
}

// Add more methods for user management, role assignment, etc.