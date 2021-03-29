package models

import "time"

type GroupAuthority struct {
	GroupId       *Group    `orm:"column(group_id);rel(fk)" description:"그룹"`
	AuthorityCd   string    `orm:"column(authority_cd);size(8)" description:"권한"`
	GroupNm       string    `orm:"column(group_nm);size(48)" description:"그룹"`
	AuthorityCdNm string    `orm:"column(authority_cd_nm);size(48)" description:"권한"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp)" description:"등록일시"`
	CreatedBy     string    `orm:"column(created_by);size(48)" description:"등록자"`
	UpdatedAt     time.Time `orm:"column(updated_at);type(timestamp)" description:"수정일시"`
	UpdatedBy     string    `orm:"column(updated_by);size(48)" description:"수정자"`
}
