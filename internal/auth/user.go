package auth

import (
	"errors"

	"elf/internal/core"
)

var (
	ErrDuplicateSub     = errors.New("An authenticated user with that 'Sub' already exists")
	ErrUserDoesNotExist = errors.New("No user with that 'Id' exists")
)

type Profile struct {
	Sub     string `json:"sub"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
}

type AuthenticatedUser struct {
	User    core.User
	Profile Profile
}

type AuthenticatedUserStore interface {
	Create(profile Profile) (*AuthenticatedUser, error)
	Read(id int) (*AuthenticatedUser, error)
}

type ArrayAuthenticatedUserStore struct {
	userIdCounter int
	users         []AuthenticatedUser
}

func (s *ArrayAuthenticatedUserStore) Create(profile Profile) (user *AuthenticatedUser, err error) {
	for _, u := range s.users {
		if u.Profile.Sub == profile.Sub {
			return nil, ErrDuplicateSub
		}
	}

	user = &AuthenticatedUser{
		User:    core.User{Id: s.nextUserId()},
		Profile: profile,
	}
	return user, nil
}

func (s *ArrayAuthenticatedUserStore) Read(id int) (*AuthenticatedUser, error) {
	for _, u := range s.users {
		if u.User.Id == id {
			return &u, nil
		}
	}
	return nil, ErrUserDoesNotExist
}

func (s *ArrayAuthenticatedUserStore) nextUserId() int {
	id := s.userIdCounter
	s.userIdCounter = s.userIdCounter + 1
	return id
}
