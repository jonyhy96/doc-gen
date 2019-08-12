package models

// User user
// apidoc motherfucker
type User struct {
	Name        string `json:"name" bson:"name" description:"用户名称"`
	Age         int8   `json:"age" bson:"age" description:"用户年龄"`
	GirlFriends Girl   `json:"girl" bson:"girl"`
}

// Girl girl
type Girl struct {
	Name string `json:"name" bson:"name" description:"名字"`
	Age  string `json:"age" bson:"age" description:"年龄"`
}

// CreateUser create user
// @name 创建用户
// @route  /v1/users post
// @in object models.User 用户
// @out string id 用户id
func CreateUser() {

}

// GetUser get user
// @name 获取所有用户
// @route  /v1/users get
// @in string name 名称
// @in string age 年龄
// @out object []models.User 用户列表
func GetUser() {

}
