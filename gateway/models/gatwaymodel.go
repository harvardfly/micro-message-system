package models

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type (
	GateWay struct {
		Id         int64
		Token      string     `json:"token" gorm:"varchar(11) notnull 'token'"`
		ImAddress  string     `json:"imAddress" gorm:"varchar(60) notnull 'im_address'"`
		ServerName string     `json:"server_name" gorm:"varchar(60) notnull 'server_name'"`
		Topic      string     `json:"topic" gorm:"varchar(60) notnull 'topic'"`
		CreateTime *time.Time `json:"createTime" gorm:"DateTime 'create_time'"`
		UpdateTime *time.Time `json:"updateTime" gorm:"DateTime 'update_time'"`
	}
	GateWayModel struct {
		mysql *gorm.DB
	}
)

func (g *GateWay) TableName() string {
	return "gateway"
}
func NewGateWayModel(mysql *gorm.DB) *GateWayModel {

	return &GateWayModel{mysql: mysql}
}

func (m *GateWayModel) FindByToken(token string) (*GateWay, error) {
	g := new(GateWay)

	if err := m.mysql.Where("token = ?", token).First(g).Error; err != nil {
		return nil, err
	}
	return g, nil
}

// 检查gateway中此用户的im websocket地址是否存在
func (m *GateWayModel) FindByServerNameTokenAddressTopic(serverName, topic, token, address string) (*GateWay, error) {
	g := new(GateWay)
	if err := m.mysql.Where(
		"token = ? and im_address =? and topic = ? and server_name=?",
		token,
		address,
		topic,
		serverName,
	).First(g).Error; err != nil {
		return nil, err
	}
	return g, nil
}

func (m *GateWayModel) Insert(g *GateWay) (*GateWay, error) {
	has, err := m.FindByServerNameTokenAddressTopic(g.ServerName, g.Topic, g.Token, g.ImAddress)
	if has != nil && has.Id > 0 && err == nil {
		return has, nil
	}
	// 记录不存在则插入
	if err := m.mysql.Create(g).Error; err != nil {
		return nil, err
	}
	return g, nil
}

func (m *GateWayModel) FindByImAddress(imAddress string) ([]*GateWay, error) {
	gs := []*GateWay(nil)
	// find 传入数组或者map的地址  返回数组/map
	if err := m.mysql.Where("im_address = ?", imAddress).Find(&gs).Error; err != nil {
		return nil, err
	}
	return gs, nil
}
