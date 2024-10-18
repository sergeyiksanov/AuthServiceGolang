package entity

type Credentials struct {
	ID       int64  `gorm:"column_id:id,primaryKey"`
	Email    string `gorm:"column_id:email,unique"`
	Password string `gorm:"column_id:password"`
}

func (c *Credentials) TableName() string {
	return "credentials"
}
