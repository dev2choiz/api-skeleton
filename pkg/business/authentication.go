package business

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/dev2choiz/api-skeleton/entity"
	"github.com/dev2choiz/api-skeleton/pkg/cache"
	"github.com/dev2choiz/api-skeleton/pkg/errapp"
	"github.com/dev2choiz/api-skeleton/pkg/logger"
)

var (
	bcryptCost         = bcrypt.DefaultCost
	tokenJWTExpiration = 24000 * time.Hour // TODO: to adjust
	cacheTTL           = 60 * time.Minute
)

var ErrInvalidJWTToken = errors.New("invalid jwt token")

// Register creates a new user with the given username and password.
// It hashes the password using bcrypt before persisting the user.
// It returns the created user or an error if validation or persistence fails.
func (b *business) Register(ctx context.Context, username, password string) (entity.User, error) {
	if username == "" {
		return entity.User{}, errapp.NewBadRequest("username cannot be empty")
	}

	if password == "" {
		return entity.User{}, errapp.NewBadRequest("password cannot be empty")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		logger.Get(ctx).Error("bcrypt hash failed", zap.Error(err))
		return entity.User{}, errapp.WrapInternal(err)
	}

	user := entity.User{
		Username: username,
		Password: string(passwordHash),
	}

	return b.repository.InsertUser(ctx, user)
}

// Authenticate verifies the provided credentials and returns a signed JWT token.
// It returns an error if the username does not exist or the password is invalid.
func (b *business) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := b.repository.GetOneUser(ctx, entity.UserFilters{Username: &username})
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errapp.WrapBadRequest(err, "username / password incorrect")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(tokenJWTExpiration).Unix(),
		"nbf": time.Now().Unix(),
	})

	return token.SignedString([]byte(b.jwtSecret))
}

// ValidateToken validates the JWT token from the Authorization header.
// It returns the associated user if the token is valid.
// It returns an error if the token is invalid, expired, or the user cannot be found.
func (b *business) ValidateToken(r *http.Request) (entity.User, error) {
	ctx := r.Context()
	log := logger.Get(ctx)

	tokStr := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
	if tokStr == "" {
		return entity.User{}, errapp.WrapBadRequest(ErrInvalidJWTToken)
	}

	token, err := jwt.Parse(tokStr, func(token *jwt.Token) (any, error) {
		return []byte(b.jwtSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return entity.User{}, errapp.WrapBadRequest(ErrInvalidJWTToken)
	}

	if !token.Valid {
		return entity.User{}, errapp.WrapBadRequest(ErrInvalidJWTToken)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return entity.User{}, errapp.WrapBadRequest(ErrInvalidJWTToken)
	}

	idRaw, ok := claims["id"]
	if !ok {
		return entity.User{}, errapp.WrapBadRequest(ErrInvalidJWTToken)
	}

	id, ok := idRaw.(string)
	if !ok || id == "" {
		return entity.User{}, errapp.WrapBadRequest(ErrInvalidJWTToken)
	}

	var user entity.User
	cacheKey := fmt.Sprintf("auth:user:%s", id)

	user, err = cache.GetJSON[entity.User](ctx, b.redis, cacheKey)
	if err == nil {
		return user, nil
	}

	v, err, _ := b.sf.Do(cacheKey, func() (any, error) {
		// retry to fetch the user in the cache
		user, err = cache.GetJSON[entity.User](ctx, b.redis, cacheKey)
		if err == nil {
			return user, nil
		}

		// not in the cache. try to get the user in the db
		user, err = b.repository.GetUser(ctx, id)
		if err != nil {
			log.Error("failed to get the user", zap.Error(err))
			return entity.User{}, fmt.Errorf("failed to get the user: %w", err)
		}

		err = b.redis.SetJSON(ctx, cacheKey, user, cacheTTL)
		if err != nil {
			log.Error("failed to save the user in redis", zap.Error(err))
		}

		return user, nil
	})

	return v.(entity.User), err
}
