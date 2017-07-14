//Author: Peter Nagy <https://peternagy.ie>
//Since: 07, 2017
//Description: --
package auth

import (
	"../common"
	"errors"
	"net"
	"strings"
)

const (
	c_id     = "CD5251A185EA50A8F9E61B77545"
	c_secret = "2017809D5E86FC2AD56E2B26660C9C637299CA127430DFBB81849145A1EE8052"
)

type Model struct {
}

type UserCredentials struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type User struct {
	ID       string `json:"_id" valid:"-"`
	Username string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required"`
}

func (m *Model) CreateUser(u *User) (bool, error) {
	u.ID = common.GetStringHash(u.Username)
	pb, e := common.CreateBcryptPassword(u.Password)

	if e != nil {
		common.CheckError("auth_model::CreateUser", e, false)
		return false, errors.New("Error while hashing password")
	}

	u.Password = string(pb)

	rc := common.GetRedisClient()
	if exists, err := rc.KeyExists(u.ID); exists || err != nil {
		return false, errors.New("User already registered")
	}

	ub, err := common.Encode(u)
	if err != nil {
		common.CheckError("auth_model::CreateUser ", err, false)
		return false, errors.New("Error while encoding")
	}

	ok := rc.SetKey(u.ID, ub, 0)
	if !ok {
		return false, errors.New("Error while storing entry")
	}

	return true, nil
}

func (m *Model) GetUser(uname string) (*User, error) {
	id := common.GetStringHash(uname)
	rc := common.GetRedisClient()

	ub := rc.GetKey(id)

	if len(ub) == 0 {
		return nil, errors.New("Failed to get user")
	}

	var u *User
	if err := common.Decode(ub, u); err != nil {
		return nil, errors.New("Failed to decode user data")
	}

	return u, nil
}

func (m *Model) AuthUser(uc *UserCredentials) (*common.TokenModel, error) {
	if !strings.EqualFold(uc.ClientID, c_id) || !strings.EqualFold(uc.ClientSecret, c_secret) {
		return nil, errors.New("Invalid client_id or client_secret")
	}

	rc := common.GetRedisClient()
	ub := rc.GetKey(common.GetStringHash(uc.Username))
	if len(ub) == 0 {
		return nil, errors.New("Could not find user")
	}

	var u *User
	if err := common.Decode(ub, u); err != nil {
		return nil, errors.New("Could not decode user")
	}

	if !common.ValidateBcryptPassword(u.Password, uc.Password) {
		return nil, errors.New("Password does not match")
	}

	tm := common.CreateTokens(u.ID, net.ParseIP("127.0.0.1"), 0)

	if tm == nil {
		return nil, errors.New("Error while generating token")
	}

	return tm, nil
}
