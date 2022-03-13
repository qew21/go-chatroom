package db

import (
	"fmt"
	"simple-chatroom/pkg/pb"
	"simple-chatroom/pkg/util"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// Message 消息
type Message struct {
	Id           int64 `gorm:"not null; PRIMARY_KEY; AUTO_INCREMENT"`    // 自增主键
	UserId       int64     // 所属类型id
	RequestId    int64     // 请求id
	SenderType   int32     // 发送者类型
	SenderId     int64     // 发送者账户id
	ReceiverType int32     // 接收者账户id
	ReceiverId   int64     // 接收者id,如果是单聊信息，则为user_id，如果是群组消息，则为group_id
	ToUserIds    string    // 需要@的用户id列表，多个用户用，隔开
	Type         int       // 消息类型
	Content      []byte    // 消息内容
	Seq          string    // 消息同步序列
	SendTime     time.Time // 消息发送时间
	Status       int32     // 创建时间
}

type User struct {
	Id          int64  `gorm:"not null; PRIMARY_KEY; AUTO_INCREMENT"`    // 用户id
	PhoneNumber string // 手机号
	Nickname    string // 昵称
	Sex         int32  // 性别
	Password    string // 哈希密码
}

// Group 群组
type Group struct {
	Id           int64   `gorm:"not null; PRIMARY_KEY; AUTO_INCREMENT"`    // 群组id
	Name         string      // 组名
	Introduction string      // 群简介
	Members      []GroupUser `gorm:"-"` // 群组成员
}

type GroupUser struct {
	Id         int64    `gorm:"not null; PRIMARY_KEY; AUTO_INCREMENT"`     // 自增主键
	GroupId    int64  // 群组id
	UserId     int64  // 用户id
	MemberType int    // 群组类型
	Remarks    string // 备注
}

var (
	DB    *gorm.DB
	MySQL = "gouser:123456@tcp(127.0.0.1:3308)/im_db?charset=utf8&parseTime=true"
	TIME_LAYOUT = "2006-01-02 15:04:05"
)

func init() {
	InitMysql(MySQL)
}

// InitMysql 初始化MySQL
func InitMysql(dataSource string) {
	fmt.Println("init mysql")
	var err error
	DB, err = gorm.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	DB.SingularTable(true)
	DB.LogMode(true)
	fmt.Println("init mysql ok")
}

func SaveMessage(jmsg pb.SocketMessage) error {
	sendTime, _ := time.Parse(TIME_LAYOUT, jmsg.Time)
	message := Message{
		UserId:       0, 
		RequestId:    0, 
		SenderType:   0, 
		SenderId:     GetUserID(jmsg.User), 
		ReceiverType: 0, 
		ReceiverId:   GetGroupID(jmsg.Group), 
		ToUserIds:    "",                 
		Type:         0,                    
		Content:      []byte(jmsg.Content), 
		Seq:          util.GetSeq(0, 0),    
		SendTime:     sendTime,             
		Status:       0,                    
	}
	err := DB.Table("message").Create(&message).Error
	if err != nil {
		fmt.Println("SaveMessage", err)
	}
	return nil
}

func PullMessage(groupName string) []pb.ChatMessage {
	receiverID := GetGroupID(groupName)
	var chats []pb.ChatMessage
	err := DB.Table("message m").Select("u.nickname as user, DATE_FORMAT(m.send_time,'%Y-%m-%d %H:%i:%s') as time, m.content").
		Joins("join `user` u on m.sender_id = u.id").
		Where("receiver_id = ?", receiverID).
		Find(&chats).Error
	if err != nil {
		fmt.Println("PullMessage", err)
	}
	return chats
}

func CreateUser(name string) error {
	user := User{
		Nickname: name,
		PhoneNumber: "133" + strconv.FormatInt(time.Now().Unix(), 10),
		Sex: 0,
		Password: "PASS",
	}
	err := DB.Table("user").Create(&user).Error
	if err != nil {
		fmt.Println("CreateUser", err)
	}
	return err
}

func CreateGroup(name string) error {
	group := Group{
		Name: name,
		Introduction: "",
	}
	err := DB.Table("group").Create(&group).Error
	if err != nil {
		fmt.Println("CreateGroup", err)
	}
	return err
}

func GetUserID(userName string) int64 {
	var users []User
	err := DB.Table("user").Where("nickname = ?", userName).First(&users).Error
	if err != nil {
		fmt.Println(err)
	}
	return users[0].Id
}

func GetGroupID(groupName string) int64 {
	var groups []Group;
	err := DB.Table("group").Where("name = ?", groupName).First(&groups).Error
	if err != nil {
		fmt.Println(err)
	}
	return groups[0].Id
}


func JoinGroup(userName string, groupName string) error {
	userID := GetUserID(userName)
	groupID := GetGroupID(groupName)
	fmt.Println(userID)
	fmt.Println(groupID)
	groupUser := GroupUser{
		GroupId:    groupID,
		UserId:     userID,
		MemberType: 0,
		Remarks:    "",
	}
	err := DB.Table("group_user").Create(&groupUser).Error
	if err != nil {
		return err
	}
	return nil
}
