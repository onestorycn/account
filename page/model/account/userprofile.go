package account

import (
	"sync"
	"account/page/service/mysql"
	"account/library/config"
	"account/library"
	"errors"
	"time"
	"strconv"
)

var (
	tableName  = "user_profile"
	once       sync.Once
	serviceMap = make(map[string]*AccountService, 0)
	lock       = &sync.Mutex{}
)

type UserProfile struct {
	Id          int64  `gorm:"PRIMARY_KEY; AUTO_INCREMENT" json:"id"`
	Openid      string `gorm:"not null" json:"openid"`
	Passid      string `gorm:"unique; not null" json:"passid"`
	Email       string `gorm:"unique; not null" json:"email"`
	Phone       int64  `gorm:"unique; not null" json:"phone"`
	Password    string `gorm:"not null" json:"password"`
	Update_time int64  `gorm:"not null" json:"update_time"`
	Nick_name   string `gorm:"not null" json:"nick_name"`
	Avatar      string `gorm:"not null" json:"avatar"`
	Ext         string `gorm:"not null" json:"ext"`
	Active      int    `gorm:"not null" json:"active"`
}

func (userprofile *UserProfile) TableName() string {
	return tableName
}

type AccountService struct {
	DbInstance *mysql.MysqlDbInfo
	env        string
}

func LoadAccountService() (accountService *AccountService) {
	var env string
	envGet := config.GetServiceEnv("env")
	if envGet == nil {
		env = "prod"
	} else {
		env = envGet.(string)
	}
	if _, ok := serviceMap[env]; ok {
		return serviceMap[env]
	}
	lock.Lock()
	defer lock.Unlock()
	if _, ok := serviceMap[env]; !ok {
		mysqlInstance := mysql.LoadMysqlConn(env)
		if mysqlInstance == nil || mysqlInstance.Conn == nil {
			return nil
		} else {
			sockService := new(AccountService)
			sockService.DbInstance = mysqlInstance
			sockService.env = env
			serviceMap[env] = sockService
		}
	}
	return serviceMap[env]
}

func (accountService *AccountService) InsertNewUser(user *UserProfile) error {
	conn := accountService.DbInstance.CheckAndReturnConn()
	getId := accountService.generatePassIdByTime()
	user.Passid = getId
	user.Openid = getId
	res := conn.Create(&user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (accountService *AccountService) GetUserById(passId, openId string) (*UserProfile, error) {
	userProfile := new(UserProfile)
	conn := accountService.DbInstance.CheckAndReturnConn()
	if conn == nil {
		return nil , errors.New("get db connection fail")
	}
	if len(openId) < 1 || openId == ""{
		if len(passId) < 1 {
			return nil, errors.New("empty passid")
		}else{
			res := conn.Where("passid = ?", passId).First(&userProfile)
			if res.Error != nil {
				return nil, res.Error
			}
			return userProfile, nil
		}
		return nil, errors.New("empty openid")
	}else{
		res := conn.Where("openid = ?", openId).First(&userProfile)
		if res.Error != nil {
			return nil, res.Error
		}
	}
	return userProfile, nil
}

func (accountService *AccountService) GetUserByPhone(phone int64) (*UserProfile, error) {
	conn := accountService.DbInstance.CheckAndReturnConn()
	if conn == nil {
		return nil , errors.New("get db connection fail")
	}
	if phone == 0 {
		return nil, errors.New("empty phone")
	}
	userProfile := new(UserProfile)
	res := conn.Where("phone = ?", phone).First(&userProfile)
	if res.Error != nil {
		return nil, res.Error
	}
	return userProfile, nil
}

func (accountService *AccountService) UpdateUser(Password, Email, NickName, Avatar, Ext , PassId string, Phone int64, isDeActive bool) (int64, error) {
	var updates = make(map[string]interface{}, 0)
	if resBool := library.IsEmpty(Password); !resBool {
		updates["password"] = Password
	}
	if resBool := library.IsEmpty(Email); !resBool {
		updates["email"] = Email
	}
	if resBool := library.IsEmpty(NickName); !resBool {
		updates["nick_name"] = NickName
	}
	if resBool := library.IsEmpty(Avatar); !resBool {
		updates["avatar"] = Avatar
	}
	if resBool := library.IsEmpty(Ext); !resBool {
		updates["ext"] = Ext
	}
	if resBool := library.IsEmpty(Phone); !resBool {
		updates["phone"] = Phone
	}

	if len(updates) == 0 {
		return 0, errors.New("no field to update")
	}
	if isDeActive {
		updates["active"] = 0
	}else{
		updates["active"] = 1
	}
	updates["update_time"] = time.Now().Unix()

	conn := accountService.DbInstance.CheckAndReturnConn()
	if conn == nil {
		return 0 , errors.New("get db connection fail")
	}
	user := new(UserProfile)
	upRes := conn.Model(&user).Where("passid = ?", PassId).Updates(updates)
	if upRes.Error != nil {
		return 0, upRes.Error
	}
	return upRes.RowsAffected, nil
}

func (accountService *AccountService) CheckUserLoginByPhone(phone int64, password string) (*UserProfile, error) {
	userProfile := new(UserProfile)
	conn := accountService.DbInstance.CheckAndReturnConn()
	if conn == nil {
		return nil , errors.New("get db connection fail")
	}
	res := conn.Where("phone = ? and password = ?", phone, password).First(&userProfile)
	if res.Error != nil {
		return nil, res.Error
	}
	return userProfile, nil
}

func  (accountService *AccountService) CheckUserLoginByEmail(email, password string) (*UserProfile, error) {
	userProfile := new(UserProfile)
	conn := accountService.DbInstance.CheckAndReturnConn()
	if conn == nil {
		return nil , errors.New("get db connection fail")
	}
	res := conn.Where("email = ? and password = ?", email, password).First(&userProfile)
	if res.Error != nil {
		return nil, res.Error
	}
	return userProfile, nil
}

func (accountService *AccountService)generatePassIdByTime() string {
	id := time.Now().Unix() * 10000 + int64(library.RandInt(1,9999))
	return strconv.FormatInt(id, 10)
}
