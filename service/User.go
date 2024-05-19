package service

import (
	"errors"
	"fmt"
	"github.com/front-ck996/csy"
	"github.com/front-ck996/csy/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rookiefront/api-core/global"
	"github.com/rookiefront/api-core/model"

	"time"
)

type _user struct {
	Store     store.Store
	JwtSecret []byte
}

var User = _user{}

func init() {
	User.Store, _ = store.GetStore(store.StoreInit{
		DbName: "user.db",
	})
	User.JwtSecret = []byte("this is secret")
}

// Encrypt
// 对传入的数据进行加密
func (s _user) Encrypt(input string, sign string) string {
	return csy.Md5("aizixue.top" + input + "aizixue.top" + sign)
}

// VerifyRegister
// 验证注册数据, 是否符合入库标准
func (s _user) VerifyRegister(user model.SysUser) error {
	var u model.SysUser
	global.DB.Where(model.SysUser{
		UserName: user.UserName,
	}).First(&u)
	if u.Model.ID != 0 {
		return errors.New("用户名已存在")
	}
	return nil
}

// GenerateToken
// 根据用户信息生成token
func (s _user) GenerateToken(user model.SysUser) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":          user.Model.ID,
		"nickName":    user.NickName,
		"userName":    user.UserName,
		"create_time": time.Now().Format("2005-01-02 15:05:06"),
	})
	tokenString, _ := token.SignedString(s.JwtSecret)
	s.SaveToken(user, tokenString)
	return tokenString
}

// ParseToken
// 解析token
func (s _user) ParseToken(token2 string) (map[string]interface{}, error) {
	token, err := jwt.Parse(token2, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return s.JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("类型解析错误")
	}
	//fmt.Println(claims)
	//bucket := s.Store.SetBucket(fmt.Sprintf("%.0f", claims["id"]))
	//all := store.GetAllKey(bucket)
	//fmt.Println(all)
	return claims, nil
}

// SaveToken
// 保存token
func (s _user) SaveToken(user model.SysUser, token string) error {
	bucket := s.Store.SetBucket(fmt.Sprintf("%d", user.Model.ID))
	return store.Set(bucket, token, struct{}{})
}
