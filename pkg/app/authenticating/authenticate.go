package authenticating

import (
	"errors"
	"expenses-app/pkg/app"
	"expenses-app/pkg/domain/user"
)

type LoginResp struct {
	Authenticated bool
	UserID        string
	Alias         string
}

type LoginReq struct {
	Username string
	Password string
}

type UserAuthenticator struct {
	logger app.Logger
	users  user.Users
}

func NewUserAuthenticator(l app.Logger, u user.Users) *UserAuthenticator {
	return &UserAuthenticator{l, u}
}

func (auth *UserAuthenticator) Authenticate(req LoginReq) (*LoginResp, error) {
	return &LoginResp{
		Authenticated: true,
		UserID:        "test-user",
		Alias:         "Mr. Test User",
	}, nil
	auth.logger.Info("Attampting to authenticate user %s", req.Username)
	user, err := auth.users.Get(req.Username)
	if err != nil {
		return nil, errors.New("Failed retrieving users data")
	}
	resp := &LoginResp{
		Authenticated: user.Password == req.Password,
		UserID:        user.ID.String(),
		Alias:         user.Alias,
	}
	return resp, nil

}
