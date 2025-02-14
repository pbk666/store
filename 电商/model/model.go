package model

type Product struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	Comment     string  `json:"comment"`
	Dating      string  `json:"dating"`
	Favorite    int     `json:"favorite"`
	Purchase    int     `json:"purchase"`
	Heat        int     `json:"heat"`
}

type Cart struct {
	ID        uint `json:"id"`
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type Comment struct {
	CommentID   string `json:"comment_id"`
	Content     string `json:"content"`
	UserID      string `json:"user_id"`
	Nickname    string `json:"nickname"`
	PraiseCount int    `json:"praise_count"`
	IsPraised   int    `json:"is_praised"`
	ProductID   int    `json:"product_id"`
}

type Order struct {
	OrderID int         `json:"order_id"`
	UserID  uint        `json:"user_id"`
	Orders  []OrderItem `json:"orders"`
	Address string      `json:"address"`
	Total   uint        `json:"total"`
}
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	Id int `json:"id"`
	//Avatar string `json:"avatar"`
	Nickname     string `json:"nick_name"`
	Introduction string `json:"introduction"`
	Phone        string `json:"phone"`
	Qq           int    `json:"qq"`
	Gender       string `json:"gender"`
	Email        string `json:"email"`
	Birthday     string `json:"birthday"`
}
type OrderItem struct {
	ID        int `json:"id"`
	OrderID   int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Rating struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Score     int `json:"score"`
}
