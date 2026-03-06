package business

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dev2choiz/api-skeleton/entity"
	"github.com/dev2choiz/api-skeleton/mocks/mockcache"
	"github.com/dev2choiz/api-skeleton/mocks/mockrepository"
	"github.com/dev2choiz/api-skeleton/pkg/errapp"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
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
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if token == "" {
				t.Fatal("expected JWT token, got empty string")
			}
		})
	}
}

func TestBusiness_ValidateToken(t *testing.T) {
	ctx := context.Background()
	jwtSecret := "mysecret"

	testUser := entity.User{
		ID:       "123",
		Username: "jaskier",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       testUser.ID,
		"username": testUser.Username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	st := mockrepository.NewMockRepository(t)
	ca := mockcache.NewMockCache(t)
	ca.EXPECT().
		GetString(ctx, "validate_token_"+testUser.ID).
		Return(`{"id": "123", "username": "jaskier"}`, nil)

	b := NewBusiness(st, ca, jwtSecret)

	got, err := b.ValidateToken(req)
	assert.NoError(t, err)

	assert.Equal(t, got.ID, testUser.ID)
	assert.Equal(t, got.Username, testUser.Username)
}
