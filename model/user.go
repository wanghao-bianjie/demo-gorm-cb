package model

type User struct {
	Id          uint   `gorm:"column:id;primaryKey"`
	Name        string `gorm:"column:name"`
	PhoneNumber string `gorm:"column:phone_number"`
	Address     string `gorm:"column:address"`
	IdNo        string `gorm:"column:id_no"`
	UpdateAt    int64  `gorm:"column:update_at"`
}

func (u *User) TableName() string {
	return "user"
}

func GetUserColumn() UserColumn {
	return userColumn
}

var userColumn = UserColumn{
	Id:          "id",
	Name:        "name",
	PhoneNumber: "phone_number",
	Address:     "address",
	IdNo:        "id_no",
	UpdateAt:    "update_at",
}

type UserColumn struct {
	Id          string
	Name        string
	PhoneNumber string
	Address     string
	IdNo        string
	UpdateAt    string
}
