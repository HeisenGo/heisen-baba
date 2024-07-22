package auth

import (
	"auth/internal/models"
	"auth/internal/tokens"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	pb "auth/api"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	pb.UnimplementedAuthServiceServer
	secret                 []byte
	tokenExpiration        uint
	refreshTokenExpiration uint
	db                     *gorm.DB
}

func NewService(db *gorm.DB, secret []byte) *Service {
	return &Service{db: db, secret: secret}
}

func (s *Service) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var user models.User
	if err := s.db.WithContext(ctx).Where("email = ?", req.Email).First(&user).Error; err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "request missing required field: Email")
	}
	if req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "request missing required field: Password")
	}

	user = models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{UserId: fmt.Sprintf("%d", user.ID)}, nil
}

func (s *Service) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user *models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	var (
		authExp    = time.Now().Add(time.Minute * time.Duration(s.tokenExpiration))
		refreshExp = time.Now().Add(time.Minute * time.Duration(s.refreshTokenExpiration))
	)

	authToken, err := tokens.CreateToken(s.secret, s.userClaims(user, authExp))
	if err != nil {
		return nil, err
	}

	refreshToken, err := tokens.CreateToken(s.secret, s.userClaims(user, refreshExp))
	if err != nil {
		return nil, err
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &pb.LoginResponse{Token: authToken, RefreshToken: refreshToken}, nil
}

func (s *Service) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	_, err := tokens.ParseToken(req.Token, s.secret)
	if err != nil {
		return &pb.ValidateTokenResponse{IsValid: false}, nil
	}
	return &pb.ValidateTokenResponse{IsValid: true}, nil
}

func (s *Service) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.LoginResponse, error) {
	claims, err := tokens.ParseToken(req.RefreshToken, s.secret)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	var user *models.User
	if err := s.db.Where("email = ?", claims.Email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	authExp := time.Now().Add(time.Minute * time.Duration(s.tokenExpiration))

	authToken, err := tokens.CreateToken(s.secret, s.userClaims(user, authExp))
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{Token: authToken, RefreshToken: req.RefreshToken}, nil
}

func (s *Service) userClaims(user *models.User, exp time.Time) *tokens.UserClaims {
	return &tokens.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: exp,
			},
		},
		UserID:       user.ID,
		IsSuperAdmin: user.IsSuperAdmin,
	}
}
