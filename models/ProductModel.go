package models



type Product struct {
	ID                string `json:"id"`
	Gender            string `json:"gender"`
	MasterCategory    string `json:"mastercategory"`
	SubCategory       string `json:"subcategory"`
	ArticleType       string `json:"articletype"`
	BaseColour        string `json:"basecolour"`
	Season            string `json:"season"`
	Year              int    `json:"year"`
	Usage             string `json:"usage"`
	ProductDisplayName string `json:"productdisplayname"`
	Link              string `json:"link"`
}
