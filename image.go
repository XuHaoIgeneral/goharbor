package harbor

import (
	`context`
	`fmt`
	`io/ioutil`
	`net/http`
	`strings`
)

const (
	IMAGE_GET      = "/api/repositories?project_id=%s&q=%s"
	IMAGE_NAGE_TAG = "/api/repositories/%s/tags/%s"
)

type Image struct {
	Project string
	Names   string
	Tag     string
}

func (c *Client) GetImageByNameAndTag(ctx context.Context, project, name, tag string) error {

	repo_name := project + "/" + name

	path := fmt.Sprintf(IMAGE_NAGE_TAG, repo_name, tag)

	req, err := http.NewRequest(http.MethodGet, c.host+path, nil)
	if err != nil {
		return err
	}

	code, body, err := c.do(ctx, req)
	if err != nil {
		return err
	}

	defer body.Close()
	body_byte, err := ioutil.ReadAll(body)
	body_str := string(body_byte)
	body_str = strings.Replace(body_str, " ", "", -1)
	body_str = strings.Replace(body_str, "\n", "", -1)
	switch code {
	case 200:
		return nil
	case 401, 403, 404:
		return fmt.Errorf(body_str)
	}
	return fmt.Errorf("err about get image by name")
}
