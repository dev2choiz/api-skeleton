package business

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/dev2choiz/api-skeleton/entity"
	"github.com/dev2choiz/api-skeleton/mocks/mockcache"
	"github.com/dev2choiz/api-skeleton/mocks/mockrepository"
	"github.com/dev2choiz/api-skeleton/pkg/errapp"
)

func TestBusiness_Register(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		inputUser entity.User
		mockSetup func(m *mockrepository.MockRepository)
		wantErr   error
	}{
		{
			name:      "empty username",
			inputUser: entity.User{Username: "", Password: "secret"},
			mockSetup: func(m *mockrepository.MockRepository) {},
			wantErr:   errapp.ErrAppBadRequest,
		},
		{
			name:      "empty password",
			inputUser: entity.User{Username: "ciri", Password: ""},
			mockSetup: func(m *mockrepository.MockRepository) {},
			wantErr:   errapp.ErrAppBadRequest,
		},
		{
			name:      "username already exists",
			inputUser: entity.User{Username: "ciri", Password: "secret"},
			mockSetup: func(m *mockrepository.MockRepository) {
				m.EXPECT().
					InsertUser(ctx, mock.AnythingOfType("entity.User")).
					Return(entity.User{}, errapp.ErrAppConflict)
			},
			wantErr: errapp.ErrAppConflict,
		},
		{
			name:      "successful registration",
			inputUser: entity.User{Username: "ciri", Password: "secret"},
			mockSetup: func(m *mockrepository.MockRepository) {
				m.EXPECT().
					InsertUser(ctx, mock.AnythingOfType("entity.User")).
					Return(entity.User{}, nil)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := mockrepository.NewMockRepository(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepository)
			}

			b := NewBusiness(mockRepository, nil, "jwtsecret")
			_, gotErr := b.Register(ctx, tt.inputUser.Username, tt.inputUser.Password)
			if tt.wantErr != nil {
				assert.ErrorIs(t, gotErr, tt.wantErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}

func TestBusiness_Authenticate(t *testing.T) {
	ctx := context.Background()
	jwtSecret := "secret"

	tests := []struct {
		name      string
		username  string
		password  string
		mockSetup func(m *mockrepository.MockRepository)
		wantErr   bool
	}{
		{
			name:     "user not found",
			username: "fitz",
			password: "secret",
			mockSetup: func(m *mockrepository.MockRepository) {
				m.EXPECT().
					GetOneUser(ctx, entity.UserFilters{Username: new("fitz")}).
					Return(entity.User{}, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name:     "wrong password",
			username: "bee",
			password: "wrong",
			mockSetup: func(m *mockrepository.MockRepository) {
				hashed, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)

				m.EXPECT().
					GetOneUser(ctx, entity.UserFilters{Username: new("bee")}).
					Return(entity.User{
						ID:       "123",
						Username: "bee",
						Password: string(hashed),
					}, nil)
			},
			wantErr: true,
		},
		{
			name:     "successful authentication",
			username: "altea",
			password: "secret",
			mockSetup: func(m *mockrepository.MockRepository) {
				hashed, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)

				m.EXPECT().
					GetOneUser(ctx, entity.UserFilters{Username: &[]string{"altea"}[0]}).
					Return(entity.User{
						ID:       "123",
						Username: "altea",
						Password: string(hashed),
					}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := mockrepository.NewMockRepository(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepository)
			}

			b := NewBusiness(mockRepository, nil, jwtSecret)

			token, err := b.Authenticate(ctx, tt.username, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				if token == "" {
					t.Fatal("expected JWT token, got empty string")
				}
			}
		})
	}
}

func TestBusiness_ValidateToken(t *testing.T) {
	jwtSecret := "secret"

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	tokWithoutID, err := tok.SignedString([]byte(jwtSecret))
	assert.NoError(t, err)

	tok = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"id": 123})
	tokWithoutIntID, err := tok.SignedString([]byte(jwtSecret))
	assert.NoError(t, err)

	tok = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"id": "123"})
	validTok, err := tok.SignedString([]byte(jwtSecret))
	assert.NoError(t, err)

	tests := []struct {
		name       string
		token      string
		setupMocks func(t *testing.T) Business
		wantErr    bool
		wantUserID string
	}{
		{
			name:       "empty token",
			token:      "",
			setupMocks: func(t *testing.T) Business { return NewBusiness(nil, nil, jwtSecret) },
			wantErr:    true,
		},
		{
			name:       "invalid token string",
			token:      "invalidToken",
			setupMocks: func(t *testing.T) Business { return NewBusiness(nil, nil, jwtSecret) },
			wantErr:    true,
		},
		{
			name:       "valid token but claims missing id",
			token:      tokWithoutID,
			setupMocks: func(t *testing.T) Business { return NewBusiness(nil, nil, jwtSecret) },
			wantErr:    true,
		},
		{
			name:       "valid token with and id, but the wrong type",
			token:      tokWithoutIntID,
			setupMocks: func(t *testing.T) Business { return NewBusiness(nil, nil, jwtSecret) },
			wantErr:    true,
		},
		{
			name:  "valid token + user in cache",
			token: validTok,
			setupMocks: func(t *testing.T) Business {
				red := mockcache.NewMockCache(t)
				red.EXPECT().GetString(mock.Anything, mock.Anything).Return(`{"id": "123"}`, nil)

				return NewBusiness(nil, red, jwtSecret)
			},
			wantUserID: "123",
		},
		{
			name:  "valid token + user not in cache + error with the repository",
			token: validTok,
			setupMocks: func(t *testing.T) Business {
				red := mockcache.NewMockCache(t)
				red.EXPECT().GetString(mock.Anything, mock.Anything).Return("", errors.New("not found"))

				rep := mockrepository.NewMockRepository(t)
				rep.EXPECT().GetUser(mock.Anything, "123").Return(entity.User{}, errors.New("some error"))

				return NewBusiness(rep, red, jwtSecret)
			},
			wantErr: true,
		},
		{
			name:  "valid token + user not in cache + no error with the repository",
			token: validTok,
			setupMocks: func(t *testing.T) Business {
				red := mockcache.NewMockCache(t)
				red.EXPECT().GetString(mock.Anything, mock.Anything).Return("", errors.New("not found"))
				red.EXPECT().SetJSON(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some error"))

				rep := mockrepository.NewMockRepository(t)
				rep.EXPECT().GetUser(mock.Anything, "123").Return(entity.User{ID: "123"}, nil)

				return NewBusiness(rep, red, jwtSecret)
			},
			wantErr:    false,
			wantUserID: "123",
		},
		{
			name:  "valid token + user not in cache the first check + but in the cache at the second check",
			token: validTok,
			setupMocks: func(t *testing.T) Business {
				red := mockcache.NewMockCache(t)
				red.EXPECT().GetString(mock.Anything, mock.Anything).Return("", errors.New("not found")).Once()
				red.EXPECT().GetString(mock.Anything, mock.Anything).Return(`{"id": "123"}`, nil).Once()

				rep := mockrepository.NewMockRepository(t)

				return NewBusiness(rep, red, jwtSecret)
			},
			wantUserID: "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.setupMocks(t)

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "Bearer "+tt.token)

			user, err := b.ValidateToken(req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, user.ID, tt.wantUserID)
			}
		})
	}
}
