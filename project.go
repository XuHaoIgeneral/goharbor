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
	PATH_FMT_PROJECT_LIST        = "/api/projects"
	PATH_FMT_PROJECT_GET         = "/api/projects/%d"
	PATH_FMT_PROJECT_ADD_MEMBER  = "/api/projects/%d/members"
	PATH_FMT_PROJECT_DEL_MEMBER  = "/api/projects/%d/members/%d"
	PATH_FMT_PROJECT_GET_MEMBERs = "/api/projects/%d/members/?entityname="
)

type ProjectOption struct {
	Name   string
	Public string
	Owner  string
	PageOption
}

type ProjectInit struct {
	Public      bool
	ProjectName string
}

type ProjectMember struct {
	RoleId   int
	Username string
	Project  string
}

func (opt *ProjectOption) Urls() url.Values {
	v := opt.PageOption.Urls()

	switch opt.Public {
	case "1", "true", "True":
		v.Set("public", "1")
	case "0", "false", "False":
		v.Set("public", "0")
	default:

	}

	if opt.Name != "" {
		v.Set("name", opt.Name)
	}
	if opt.Owner != "" {
		v.Set("owner", opt.Owner)
	}

	return v
}

// todo : lower harbor public!!!!!!!!!!!
func (c *Client) CreateProject(ctx context.Context, opt *ProjectInit) (bool, error) {
	path := PATH_FMT_PROJECT_LIST
	if opt == nil {
		return false, ERROR_THE_GINSENG
	}

	metadata := "false"
	if opt.Public {
		metadata = "true"
	}
	ret := &models.InitProject{
		ProjectName: opt.ProjectName,
		Metadata: models.InitProjectMetadata{
			Public: metadata,
		},
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

func (c *Client) ListProjects(ctx context.Context, opt *ProjectOption) ([]*models.Project, error) {
	var ret []*models.Project

	path := PATH_FMT_PROJECT_LIST
	if opt != nil {
		path += "?" + opt.Urls().Encode()
	}

	// log.Print(path)
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

func (c *Client) GetProjectById(ctx context.Context, id int64) (*models.Project, error) {
	ret := new(models.Project)

	path := fmt.Sprintf(PATH_FMT_PROJECT_GET, id)
	req, err := http.NewRequest(http.MethodGet, c.host+path, nil)
	if err != nil {
		return ret, err
	}

	err = c.doJson(ctx, req, ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func (c *Client) GetProjectByName(ctx context.Context, projectName string) (*models.Project, error) {
	projects, err := c.ListProjects(ctx, &ProjectOption{Name: projectName})
	if err != nil {
		return nil, err
	}
	project := new(models.Project)
	switch len(projects) {
	case 0:
		return nil, NotFoundError
	case 1:
		if projects[0].Name != projectName {
			return nil, NotFoundError
		}
		project = projects[0]
	default:
		isProject := false
		for _, p := range projects {
			if p.Name == projectName {
				project = p
				isProject = true
				break
			}
		}
		if !isProject {
			return nil, NotFoundError
		}
	}
	return project, nil
}

func (c *Client) DeleteProjectById(ctx context.Context, id int64) (deleted bool, err error) {

	path := fmt.Sprintf(PATH_FMT_PROJECT_GET, id)
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

func (c *Client) DeleteProjectByName(ctx context.Context, project string) (deleted bool, err error) {
	projectModel, err := c.GetProjectByName(ctx, project)
	if err != nil {
		return false, err
	}
	projectId := projectModel.ProjectID
	return c.DeleteProjectById(ctx, projectId)
}

func (c *Client) AddMemberToProject(ctx context.Context, pm *ProjectMember) (bool, error) {
	// 检查project是否存在
	projects, err := c.GetProjectByName(ctx, pm.Project)
	if err != nil {
		return false, err
	}

	if ok, err := c.UserIsExist(ctx, pm.Username); err != nil && !ok {
		return false, err
	}

	type MemberUser struct {
		Username string `json:"username"`
	}
	ret := &struct {
		RoleId     int        `json:"role_id"`
		MemberUser MemberUser `json:"member_user"`
	}{
		RoleId: pm.RoleId,
		MemberUser: MemberUser{
			Username: pm.Username,
		},
	}
	// struct to string
	out, err := json.Marshal(ret)
	if err != nil {
		return false, err
	}
	str_out := string(out)
	req_body := strings.NewReader(str_out)
	path := fmt.Sprintf(PATH_FMT_PROJECT_ADD_MEMBER, projects.ProjectID)
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

func (c *Client) GetProjectMetadataList(ctx context.Context, project string) ([]*models.ProjectMember, error) {
	var ret []*models.ProjectMember
	projects, err := c.GetProjectByName(ctx, project)
	if err != nil {
		return ret, err
	}
	projectId := projects.ProjectID
	path := fmt.Sprintf(PATH_FMT_PROJECT_GET_MEMBERs, projectId)
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

func (c *Client) GetProjectMetadataByUser(ctx context.Context, username, project string) ([]*models.ProjectMember, error) {
	var ret []*models.ProjectMember
	projects, err := c.GetProjectByName(ctx, project)
	if err != nil {
		return ret, err
	}
	projectId := projects.ProjectID
	path := fmt.Sprintf(PATH_FMT_PROJECT_GET_MEMBERs+"%s", projectId, username)
	req, err := http.NewRequest(http.MethodGet, c.host+path, nil)
	if err != nil {
		return nil, err
	}

	err = c.doJson(ctx, req, &ret)
	if err != nil {
		return nil, err
	}

	retData := make([]*models.ProjectMember, 0)
	// filler
	if len(ret) != 1 {
		for _, v := range ret {
			if v.EntityName != username {
				continue
			}
			retData = append(retData, v)
			break
		}
	}
	if len(retData) != 1 || retData[0].EntityName != username {
		return nil, NotFoundError
	}

	return retData, nil
}

func (c *Client) DelMemberToProject(ctx context.Context, username, project string) (bool, error) {
	projects, err := c.GetProjectByName(ctx, project)
	if err != nil {
		return false, err
	}
	users, err := c.GetProjectMetadataByUser(ctx, username, project)
	if err != nil {
		return false, err
	}
	user := users[0]

	path := fmt.Sprintf(PATH_FMT_PROJECT_DEL_MEMBER, projects.ProjectID, user.Id)

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
