package harbor

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/XuHaoIgeneral/goharbor/models"
)

const (
	PATH_FMT_USER_CREATE   = "/api/users"
	PATH_FMT_USER_LIST     = "/api/users"
	PATH_FMT_USER_SEARCH   = "/api/users/search?username=%s"
	PATH_FMT_USER_PASSWORD = "/api/users/%d/password"
	PATH_FMT_USER_DELETE   = "/api/users/%d"
)

type UserInit struct {
	Username string
	Password string
	Realname string
	Email    string
}

type UserOption struct {
	Username string
	Email    string
	PageOption
}

type UserUpPwd struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (opt *UserOption) Urls() url.Values {
	v := opt.PageOption.Urls()

	if opt.Username != "" {
		v.Set("username", opt.Username)
	}
	if opt.Email != "" {
		v.Set("email", opt.Email)
	}
	return v
}

func (c *Client) SearchUser(ctx context.Context, username string) ([]*models.SearchUser, error) {
	var ret []*models.SearchUser
	path := fmt.Sprintf(PATH_FMT_USER_SEARCH, username)
	req, err := http.NewRequest(http.MethodGet, c.host+path, nil)
	if err != nil {
		return ret, err
	}
	err = c.doJson(ctx, req, &ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func (c *Client) ListUser(ctx context.Context, opt *UserOption) ([]*models.Users, error) {
	var ret []*models.Users
	path := PATH_FMT_USER_LIST
	if opt != nil {
		path += "?" + opt.Urls().Encode()
	}
	req, err := http.NewRequest(http.MethodGet, c.host+path, nil)
	if err != nil {
		return ret, err
	}
	err = c.doJson(ctx, req, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func (c *Client) UserIsExist(ctx context.Context, username string) (bool, error) {
	_, err := c.GetUserByName(ctx, username)
	if err != nil {
		return false, fmt.Errorf("username=%s is not find", username)
	}
	return true, nil
}

func (c *Client) GetUserByName(ctx context.Context, username string) (*models.SearchUser, error) {
	users, err := c.SearchUser(ctx, username)
	if err != nil {
		return nil, err
	}
	// 加入全匹配规则
	user := new(models.SearchUser)
	switch len(users) {
	case 0:
		return nil, NotFoundError
	case 1:
		if users[0].Username != username {
			return nil, NotFoundError
		}
		user = users[0]
	default:
		isUser := false
		for _, u := range users {
			if u.Username == username {
				user = u
				isUser = true
				break
			}
		}
		if !isUser {
			return nil, NotFoundError
		}
	}
	return user, nil
}

func (c *Client) CreateUser(ctx context.Context, opt *UserInit) (bool, error) {

	path := PATH_FMT_USER_CREATE
	if opt == nil {
		return false, ERROR_THE_GINSENG
	}

	ret := &models.InitUser{
		Username: opt.Username,
		Password: opt.Password,
		Realname: opt.Realname,
		Email:    opt.Email,
	}

	// struct to string
	out, err := json.Marshal(ret)
	if err != nil {
		return false, err
	}

	str_out := string(out)
	req_body := strings.NewReader(str_out)
	req, err := http.NewRequest(http.MethodPost, c.host+path, req_body)
	if err != nil {
		return false, err
	}
	code, body, err := c.do(ctx, req)
	if err != nil {
		return false, err
	}
	defer body.Close()

	body_byte, err := ioutil.ReadAll(body)

	body_str := string(body_byte)

	switch code {
	case 200, 201:
		return true, nil
	default:
		return false, fmt.Errorf("error info : %s", body_str)
	}
}

// TODO:harbor api is false. the same as update password and reset password
func (c *Client) UpdateUserPwd(ctx context.Context, username string, upwd *UserUpPwd) (bool, error) {

	if upwd == nil {
		return false, ERROR_THE_GINSENG
	}
	user, err := c.GetUserByName(ctx, username)
	if err != nil {
		return false, err
	}
	userId := user.UserId

	path := fmt.Sprintf(PATH_FMT_USER_PASSWORD, userId)

	out, err := json.Marshal(upwd)
	str_out := string(out)
	req_body := strings.NewReader(str_out)
	req, err := http.NewRequest(http.MethodPut, c.host+path, req_body)
	if err != nil {
		return false, err
	}
	code, body, err := c.do(ctx, req)
	if err != nil {
		return false, err
	}
	defer body.Close()

	body_byte, err := ioutil.ReadAll(body)

	body_str := string(body_byte)

	switch code {
	case 200, 201:
		return true, nil
	default:
		return false, fmt.Errorf("error info : %s", body_str)
	}
}

func (c *Client) DeleteUser(ctx context.Context, username string) (deleted bool, err error) {

	user, err := c.GetUserByName(ctx, username)
	if err != nil {
		return false, err
	}
	userId := user.UserId

	path := fmt.Sprintf(PATH_FMT_USER_DELETE, userId)

	req, err := http.NewRequest(http.MethodDelete, c.host+path, nil)
	if err != nil {
		return false, err
	}
	code, body, err := c.do(ctx, req)
	if err != nil {
		return false, err
	}

	defer body.Close()

	body_byte, err := ioutil.ReadAll(body)

	body_str := string(body_byte)

	switch code {
	case 200, 201:
		return true, nil
	default:
		return false, fmt.Errorf("error info : %s", body_str)
	}
}

// TODO:Harbor's pwd manage is `Chaos`,pls commit ur code!!!
func (c *Client) ResetPwd(ctx context.Context, username string, newPassword string) (reseted bool, err error) {
	return c.UpdateUserPwd(ctx, username, &UserUpPwd{NewPassword: newPassword, OldPassword: newPassword})
}
