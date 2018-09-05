package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var DB *gorm.DB

func init() {
	dbUser := "root"
	dbPassWd := "root"
	host := "localhost"
	dbName := "juejin"
	DB, _ = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		dbUser, dbPassWd, host, dbName, "Asia%2FShanghai"))
	DB.DB().SetMaxIdleConns(100)
	DB.DB().SetMaxOpenConns(500)
	DB.DB().SetConnMaxLifetime(time.Minute)
	DB.LogMode(true)

	//DB.AutoMigrate(&User{})
}

// 两种登录方式
// 手机号和密码
// 手机号和验证码
type User struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ObjectID  string    `json:"objectID"` // 掘金用户唯一 ID
	Username  string    `json:"username"`
	Company   string    `json:"company"`
	Followed  int       `gorm:"type:tinyint;not null;default:0" json:"followed"` // 我是否关注
	Checked   int       `gorm:"type:tinyint;not null;default:0" json:"checked"`  // 是否检测该账号
}

func (u *User) Create() error {
	return DB.Create(u).Error
}

func (User) TableName() string {
	return "users"
}

func (u *User) TxCreate(tx *gorm.DB) error {
	err := DB.Create(u).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

func (u *User) GetByID() error {
	if 0 == u.ID {
		return fmt.Errorf("miss primary key")
	}
	return DB.First(u, u.ID).Error
}

func (u *User) FindByObjectID() bool {
	if "" == u.ObjectID {
		return true
	}
	return DB.First(u, "object_id = ?", u.ObjectID).RecordNotFound()
}

func (u *User) Updates(params *User) error {
	if 0 == u.ID {
		return fmt.Errorf("miss primary key")
	}
	return DB.Model(u).Updates(params).Error
}

func (u *User) UpdateChecked() error {
	if 0 == u.ID {
		return fmt.Errorf("miss primary key")
	}
	return DB.Model(u).Update("checked", 1).Error
}

func (u *User) UpdateFollowed() error {
	if 0 == u.ID {
		return fmt.Errorf("miss primary key")
	}
	return DB.Model(u).Update("followed", 1).Error
}

//func (u *User) UpdateXXX(p interface) error {
//	return DB.Model(u).Update("xxx", p).Error
//}

func (u *User) Delete() error {
	if 0 == u.ID {
		return fmt.Errorf("miss primary key")
	}
	return DB.Model(u).Update("display", 0).Error
}

func (u *User) Recovery() error {
	if 0 == u.ID {
		return fmt.Errorf("miss primary key")
	}
	return DB.Model(u).Update("display", 1).Error
}

func (u *User) RealDelete() error {
	if 0 == u.ID {
		return fmt.Errorf("miss primary key")
	}
	return DB.Unscoped().Delete(u).Error
}
