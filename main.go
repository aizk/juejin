package main

import (
	"github.com/parnurzeal/gorequest"
	"fmt"
	"encoding/json"
	"juejin/model"
	"time"
)

var Request *gorequest.SuperAgent
// 你的 UserID
var UID string
// token
var Token string

func main() {
	Request = gorequest.New()
	UID = ""
	Token = ""
	Request.Set("token", Token)

	// 刚开始需要一个初始节点
	err := Process("")
	fmt.Println(err)
	return
}

func Process(userID string) (err error) {
	// 种子用户
	if "" != userID {
		err = HandleUserFolloweeList(userID)
		if err != nil {
			return
		}
		err = HandleUserFollowerList(userID)
		if err != nil {
			return
		}
	}
	//for {
		//time.Sleep(time.Millisecond * 10)
		var users []*model.User
		err = model.DB.Where("checked = ?", 0).Find(&users).Error
		if err != nil {
			return
		}
		for _, user := range users {
			err = HandleUserFolloweeList(user.ObjectID)
			if err != nil {
				return
			}
			err = HandleUserFollowerList(user.ObjectID)
			if err != nil {
				return
			}
			err = Follow(user.ObjectID)
			if err != nil {
				return
			}
			err = user.UpdateFollowed()
			if err != nil {
				return
			}
			// 检查完一个 User
			err = user.UpdateChecked()
			if err != nil {
				return
			}
		}
	//}
	return
}

type FolloweeList struct {
	ObjectID string
	FollowerID string `json:"followerId"` // 关注者 ID
	Followee Followee `json:"followee"` // 关注人信息
}

type Followee struct {
	ObjectID string `json:"objectId"` // 被关注者 ID
	Username string `json:"username"`
	Company string `json:"company"`
}

type FolloweeRsp struct {
	D []*FolloweeList `json:"d"`
	M string `json:"m"`
	S int `json:"s"`
}

// 用户关注的用户列表
func HandleUserFolloweeList(userID string) (err error) {
	time.Sleep(time.Millisecond)
	_, body, err1 := Request.Get(fmt.Sprintf("https://follow-api-ms.juejin.im/v1/getUserFolloweeList?uid=%s&currentUid=%s&src=web", userID, UID)).End()
	if len(err1) > 0 {
		err = err1[0]
		return
	}
	var r FolloweeRsp
	json.Unmarshal([]byte(body), &r)
	// 拿到多个该用户关注者
	// 新用户加入数据库
	for _, follow := range r.D {
		u := model.User{
			ObjectID: follow.Followee.ObjectID,
			Username: follow.Followee.Username,
			Company: follow.Followee.Company,
		}
		if u.FindByObjectID() {
			err = u.Create()
			if err != nil {
				return
			}
		}
	}
	return
}

type FollowerList struct {
	ObjectID string
	FollowerID string `json:"followerId"` // 关注者 ID
	Follower Follower `json:"follower"` // 关注人信息
}

type Follower struct {
	ObjectID string `json:"objectId"` // 被关注者 ID
	Username string `json:"username"`
	Company string `json:"company"`
}

type FollowerRsp struct {
	D []*FollowerList `json:"d"`
	M string `json:"m"`
	S int `json:"s"`
}

// 关注用户的用户列表
func HandleUserFollowerList(userID string) (err error) {
	time.Sleep(time.Millisecond)
	_, body, err1 := Request.Get(fmt.Sprintf("https://follow-api-ms.juejin.im/v1/getUserFollowerList?uid=%s&currentUid=%s&src=web", userID, UID)).End()
	if len(err1) > 0 {
		err = err1[0]
		return
	}
	var r FollowerRsp
	json.Unmarshal([]byte(body), &r)
	// 拿到多个该用户关注者
	// 新用户加入数据库
	for _, follow := range r.D {
		u := model.User{
			ObjectID: follow.Follower.ObjectID,
			Username: follow.Follower.Username,
			Company: follow.Follower.Company,
		}
		if u.FindByObjectID() {
			err = u.Create()
			if err != nil {
				return
			}
		} else {
			fmt.Println("exist.")
		}
	}
	return
}

// https://follow-api-ms.juejin.im/v1/follow?
func Follow(userID string) (err error) {
	src := fmt.Sprintf("https://follow-api-ms.juejin.im/v1/follow?follower=%s&followee=%s&token=%s", UID, userID, Token) + "%3D%3D&device_id=1535505892825&src=web"
	_, _, err1 := Request.Get(src).End()
	if len(err1) > 0 {
		err = err1[0]
		return
	}
	return
}