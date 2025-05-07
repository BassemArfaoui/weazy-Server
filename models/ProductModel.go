package models

type Product struct {
	ID                 string `gorm:"primaryKey;column:id" json:"id"`
	Gender             string `gorm:"column:gender" json:"gender"`
	MasterCategory     string `gorm:"column:mastercategory" json:"mastercategory"`
	SubCategory        string `gorm:"column:subcategory" json:"subcategory"`
	ArticleType        string `gorm:"column:articletype" json:"articletype"`
	BaseColour         string `gorm:"column:basecolour" json:"basecolour"`
	Season             string `gorm:"column:season" json:"season"`
	Year               int    `gorm:"column:year" json:"year"`
	Usage              string `gorm:"column:usage" json:"usage"`
	ProductDisplayName string `gorm:"column:productdisplayname" json:"productdisplayname"`
	Link               string `gorm:"column:link" json:"link"`
	IsLiked            bool   `gorm:"-" json:"is_liked"`
}
