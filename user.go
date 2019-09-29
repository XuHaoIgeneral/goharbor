package harbor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XuHaoIgeneral/goharbor/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	PATH_FMT_USER_CREATE   = "/api/users"
	PATH_FMT_USER_LIST     = "/api/users"
	PATH_FMT_USER_SEARCH   = "/api/users/search"
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

func (c *Client) ListUser(ctx context.Context, opt *UserOption) ([]*models.Users, error) {
	var ret []*models.Users
	path := PATH_FMT_USER_LIST
	if opt != nil {
		path += "?" + opt.Urls().Encode()
	}
	//log.Print(path)
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
	opt := &UserOption{
		Username: username,
	}
	list_user, err := c.ListUser(ctx, opt)
	if err != nil {
		return false, err
	}
	if len(list_user) < 1 {
		return false, fmt.Errorf("username=%s is not find", username)
	}
	return true, nil
}

func (c *Client) GetUserByName(ctx context.Context, username string) (*models.Users, error) {
	users, err := c.ListUser(ctx, &UserOption{Username: username})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, NotFoundError
	}
	return users[0], nil
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
	//log.Print(path)
	req, err := http.NewRequest(http.MethodPost, c.host+path, req_body)
	if err != nil {
		return false, err
	}
	code, body, err := c.do(ctx, req)
	if err != nil {
		return false, err
	}
	defer body.Close()

	switch code {
	case 200, 201:
		return true, nil
	case 400:
		return false, ERROR_THE_FORMAT
	case 403:
		return false, ERROR_THE_PERMISSIONS
	case 415:
		return false, ERROR_THE_TYPE
	case 500:
		return false, ERROR_THE_SERVER
	default:
		return false, ERROR_THE_PKG
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
	//log.Print(path)
	req, err := http.NewRequest(http.MethodPut, c.host+path, req_body)
	if err != nil {
		return false, err
	}
	code, body, err := c.do(ctx, req)
	if err != nil {
		return false, err
	}
	defer body.Close()

	switch code {
	case 200, 201:
		return true, nil
	case 400:
		return false, fmt.Errorf("Invalid user ID; Old password is blank; New password is blank.")
	case 401:
		return false, fmt.Errorf("Don't have authority to change password. Please check login status.")
	case 403:
		return false, fmt.Errorf("The caller does not have permission to update the password of the user with" +
			" given ID, or the old password in request body is not correct.")
	case 415:
		return false, ERROR_THE_TYPE
	case 500:
		return false, ERROR_THE_SERVER
	default:
		return false, ERROR_THE_PKG
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

	fmt.Println(code)

	body_byte, err := ioutil.ReadAll(body)

	body_str := string(body_byte)
	fmt.Println(body_str)
	switch code {
	case 200, 201:
		return true, nil
	default:
		return false, err
	}
}

// TODO:Harbor's pwd manage is `Chaos`,pls commit ur code!!!
func (c *Client) ResetPwd(ctx context.Context, username string, newPassword string) (reseted bool, err error) {
	return c.UpdateUserPwd(ctx, username, &UserUpPwd{NewPassword: newPassword, OldPassword: newPassword})
}
