package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type (
	Members struct {
		Id         int64
		Token      string    `json:"token" gorm:"varchar(11) notnull 'token'"`
		Username   string    `json:"username" gorm:"varchar(60) notnull 'username'"`
		Password   string    `json:"password" gorm:"varchar(60) notnull 'password'"`
		CreateTime time.Time `json:"createTime" gorm:"DateTime 'create_time'"`
		UpdateTime time.Time `json:"updateTime" gorm:"DateTime 'update_time'"`
	}
	MembersModel struct {
		mysql *gorm.DB
	}
)

func NewMembersModel(mysql *gorm.DB) *MembersModel {

	return &MembersModel{
		mysql: mysql,
	}
}

func (m *MembersModel) FindByToken(token string) (*Members, error) {
	member := new(Members)
	if err := m.mysql.Where("token=?", token).First(member).Error; err != nil {
		return nil, err
	}
	return member, nil
}

func (m *MembersModel) FindById(id int64) (*Members, error) {
	member := new(Members)
	if err := m.mysql.Where("id=?", id).First(member).Error; err != nil {
		return nil, err
	}
	return member, nil
}

func (m *MembersModel) FindByUserName(userName string) (*Members, error) {
	member := new(Members)
	fmt.Println(userName)
	err := m.mysql.Where("username=?", userName).First(member).Error
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (m *MembersModel) InsertMember(member *Members) (*Members, error) {
	if err := m.mysql.Create(member).Error; err != nil {
		return nil, err
	}
	return member, nil
}
