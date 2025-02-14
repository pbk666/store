package main

import (
	"decleration/api"
	"decleration/dao"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	dao.InitDB()
	h := server.Default()
	// 用户相关路由分组
	userGroup := h.Group("/user")
	{
		userGroup.POST("/register", api.Register)
		userGroup.GET("/login", api.Login)
		userGroup.GET("/token/refresh", api.RefreshToken)
		userGroup.PUT("/password", api.ChangePassword)
		userGroup.GET("/info/:user_id", api.GetUserInfo)
		userGroup.PUT("/info", api.ChangeUserInfo)
	}
	//商品相关路由分组
	productGroup := h.Group("/product")
	{
		productGroup.GET("/list", api.GetProductList)
		productGroup.GET("/search", api.SearchProduct)
		productGroup.PUT("/addCart", api.AddCart)
		productGroup.GET("/cart", api.CartList)
		productGroup.GET("/info/:product_id", api.GetProductInfo)
		productGroup.GET("/:type", api.GetTypeList)
		productGroup.PUT("/praise", api.ProductPraise)
		productGroup.GET("/heat", api.HeatProduct)
	}
	//评论相关
	commentGroup := h.Group("/comment")
	{
		commentGroup.GET("/:product_id", api.GetComment)
		commentGroup.POST("/:product_id", api.PostComment)
		commentGroup.DELETE("/:comment_id", api.DeleteComment)
		commentGroup.PUT("/:comment_id", api.UpdateComment)
		commentGroup.PUT("/praise", api.CommentPraise)
	}
	//操作相关
	operateGroup := h.Group("/operate")
	{
		operateGroup.POST("/order", api.Order)
		operateGroup.POST("/rating", api.Rating)
	}
	h.Spin()
}
