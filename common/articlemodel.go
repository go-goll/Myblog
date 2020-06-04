package common

type Article struct {
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	UserId      string    `json:"-" gorm:"index"`
	CategoryId  string    `json:"-" gorm:"index"`
	CreateUser  *User     `json:"createUser" gorm:"foreignkey:UserId" `
	Category    *Category `json:"category" gorm:"foreignkey:CategoryId"`
	Tags        []*Tag    `json:"tags" gorm:"many2many:article_tags"`
	BaseModel
}
