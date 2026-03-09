package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dev2choiz/api-skeleton/entity"
	"github.com/dev2choiz/api-skeleton/pkg/errapp"
	"github.com/dev2choiz/api-skeleton/pkg/fixtures"
)

func Test_repository_GetUser(t *testing.T) {
	resetDB(t)
	s := newTestRepository()

	tests := []struct {
		name    string
		id      string
		want    entity.User
		wantErr error
	}{
		{
			"should return Geralt user",
			fixtures.Users[0].ID,
			fixtures.Users[0],
			nil,
		},
		{
			"should  return error for no existing id",
			"557d6749-3c88-4034-bfda-660ee628ab28",
			entity.User{},
			errapp.ErrAppNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := s.GetUser(context.Background(), tt.id)
			if tt.wantErr != nil {
				assert.ErrorIs(t, gotErr, tt.wantErr)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func Test_repository_GetUsers(t *testing.T) {
	resetDB(t)
	s := newTestRepository()

	tests := []struct {
		name    string
		filters entity.UserFilters
		want    []entity.User
		wantLen *int
	}{
		{
			name: "filter by no existing ID",
			filters: entity.UserFilters{
				ID: new("d264b37e-1f1b-49c5-8b46-75adf7e644ee"),
			},
			want: []entity.User{},
		},
		{
			name: "filter by existing ID",
			filters: entity.UserFilters{
				ID: new("38dfac55-8ff6-4c9d-8916-9016052f50cd"),
			},
			want: []entity.User{fixtures.Users[0]},
		},
		{
			name: "filter by existing Username",
			filters: entity.UserFilters{
				Username: new("geralt"),
			},
			want: []entity.User{fixtures.Users[0]},
		},
		{
			name: "filter by existing Firstname",
			filters: entity.UserFilters{
				Firstname: new("Geralt"),
			},
			want: []entity.User{fixtures.Users[0]},
		},
		{
			name: "filter by existing Lastname",
			filters: entity.UserFilters{
				Lastname: new("de Riv"),
			},
			want: []entity.User{fixtures.Users[0]},
		},
		{
			name:    "no filter returns all users",
			filters: entity.UserFilters{},
			want:    fixtures.Users,
		},
		{
			name: "no filter but with limit",
			filters: entity.UserFilters{
				Limit: new(2),
			},
			wantLen: new(2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := s.GetUsers(context.Background(), tt.filters)
			assert.NoError(t, gotErr)

			if tt.want != nil {
				assert.Len(t, got, len(tt.want))

				for i := range got {
					assert.Equal(t, got[i].Username, tt.want[i].Username)
				}
			} else {
				assert.Len(t, got, *tt.wantLen)
			}
		})
	}
}

func Test_repository_GetOneUser(t *testing.T) {
	resetDB(t)
	s := newTestRepository()

	tests := []struct {
		name    string
		filters entity.UserFilters
		want    entity.User
		wantErr error
	}{
		{
			name: "no existing user",
			filters: entity.UserFilters{
				Firstname: new("Altea"),
			},
			want:    entity.User{},
			wantErr: ErrNotFound,
		},
		{
			name:    "too many results",
			filters: entity.UserFilters{},
			want:    entity.User{},
			wantErr: ErrTooManuResults,
		},
		{
			name: "existing user by Username",
			filters: entity.UserFilters{
				Firstname: new("Richard"),
			},
			want:    fixtures.Users[5],
			wantErr: nil,
		},
		{
			name: "existing user by ID",
			filters: entity.UserFilters{
				ID: &fixtures.Users[0].ID,
			},
			want:    fixtures.Users[0],
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := s.GetOneUser(context.Background(), tt.filters)
			if tt.wantErr != nil {
				assert.ErrorIs(t, gotErr, tt.wantErr)
			} else {
				assert.NoError(t, gotErr)

				if tt.wantErr != nil {
					t.Fatal("GetOneUser() succeeded unexpectedly")
				}

				if got.ID != tt.want.ID || got.Username != tt.want.Username {
					t.Fatalf("GetOneUser() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}

func Test_repository_InsertUser(t *testing.T) {
	resetDB(t)
	s := newTestRepository()

	tests := []struct {
		name    string
		user    entity.User
		wantErr error
	}{
		{
			name: "insert new user with existing username",
			user: entity.User{
				Username: "ciri",
				Password: "hashedpassword",
			},
			wantErr: errapp.ErrAppConflict,
		},
		{
			name: "insert user without username",
			user: entity.User{
				Username: "",
				Password: "hashedpassword",
			},
			wantErr: errapp.ErrAppConflict,
		},
		{
			name: "insert user without password",
			user: entity.User{
				Username: "triss",
				Password: "",
			},
			wantErr: errapp.ErrAppConflict,
		},
		{
			name: "insert new user successfully",
			user: entity.User{
				Username: "polgara",
				Password: "hashedpassword",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := s.InsertUser(context.Background(), tt.user)
			if tt.wantErr != nil {
				assert.ErrorIs(t, gotErr, tt.wantErr)
			} else {
				assert.NoError(t, gotErr)

				if got.Username != tt.user.Username {
					t.Fatalf("InsertUser() = %+v, want Username %v", got, tt.user.Username)
				}

				if got.ID == "" {
					t.Fatal("InsertUser() returned empty ID")
				}
			}
		})
	}
}
