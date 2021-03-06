package models

import (
	"time"
)

type Project struct {
	ProjectID    int64             `orm:"pk;auto;column(project_id)" json:"project_id"`
	OwnerID      int               `orm:"column(owner_id)" json:"owner_id"`
	Name         string            `orm:"column(name)" json:"name"`
	CreationTime time.Time         `orm:"column(creation_time);auto_now_add" json:"creation_time"`
	UpdateTime   time.Time         `orm:"column(update_time);auto_now" json:"update_time"`
	Deleted      interface{}       `orm:"column(deleted)" json:"deleted"`
	OwnerName    string            `orm:"-" json:"owner_name"`
	Togglable    bool              `orm:"-" json:"togglable"`
	Role         int               `orm:"-" json:"current_user_role_id"`
	RepoCount    int64             `orm:"-" json:"repo_count"`
	ChartCount   uint64            `orm:"-" json:"chart_count"`
	Metadata     map[string]string `orm:"-" json:"metadata"`
}

type InitProject struct {
	ProjectName string              `json:"project_name"`
	Metadata    InitProjectMetadata `json:"metadata"`
}

type InitProjectMetadata struct {
	Public string `json:"public"`
}

type ProjectMembers struct {
	ProjectMemberEntity []ProjectMember
}
type ProjectMember struct {
	Id         int    `json:"id"`
	ProjectId  int    `json:"project_id"`
	EntityName string `json:"entity_name"`
	RoleName   string `json:"role_name"`
	RoleId     int    `json:"role_id"`
	EntityId   int    `json:"entity_id"`
	EntityType string `json:"entity_type"`
}

func (p *Project) IsDeleted() bool {
	return getBool(p.Deleted)
}
