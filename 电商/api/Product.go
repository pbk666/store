package api

import (
	"context"
	"database/sql"
	"decleration/service"
	"decleration/utils1"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"strconv"
)

// 获取商品列表
func GetProductList(ctx context.Context, c *app.RequestContext) {
	products, err := service.GetProductList(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"info":     "success",
		"status":   10000,
		"products": products,
	})
}

// 搜索商品
func SearchProduct(ctx context.Context, c *app.RequestContext) {
	// 获取查询参数
	searchQuery := c.DefaultQuery("query", "") // 如果没有传 `query`，默认空字符串

	// 检查输入是否为空
	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, utils.H{"message": "请输入搜索关键词"})
		return
	}
	products, err := service.SearchProduct(ctx, searchQuery)
	if err != nil {
		c.JSON(500, utils.H{"message": "搜索失败"})
		return
	}

	// 返回查询结果
	c.JSON(http.StatusOK, utils.H{
		"info":     "success",
		"status":   10000,
		"products": products,
	})
}

// 添加商品到购物车
func AddCart(ctx context.Context, c *app.RequestContext) {
	auth, _ := utils1.ExtractToken(c)
	userID, _ := utils1.ValidateToken(auth)
	//获取商品id
	idStr := c.PostForm("id")
	product, err := service.AddCartService(ctx, userID, idStr)
	if err != nil {
		c.JSON(500, utils.H{"message": "添加失败"})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"indo":    "success",
		"status":  10000,
		"message": "添加成功",
		"product": product,
	})
}

// 获取购物车列表
func CartList(ctx context.Context, c *app.RequestContext) {
	// 获取并验证 token
	auth, _ := utils1.ExtractToken(c)
	utils1.ValidateToken(auth)

	// 获取 userID
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, utils.H{"message": "userID不能为空"})
		return
	}

	//获取购物车数据
	cart, err := service.GetCartList(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"message": "获取购物车数据失败"})
		return
	}

	// 查询成功
	c.JSON(http.StatusOK, utils.H{
		"info":    "success",
		"status":  10000,
		"message": "查询成功",
		"product": cart,
	})
}

// 获取商品信息
func GetProductInfo(ctx context.Context, c *app.RequestContext) {
	//获取商品id
	productID := c.Query("id")

	product, err := service.GetProductInfo(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{"message": "获取商品信息失败"})
	}

	//查询成功
	c.JSON(http.StatusOK, utils.H{
		"info":    "success",
		"status":  10000,
		"product": product,
	})
}

// 查询某一类型的商品
func GetTypeList(ctx context.Context, c *app.RequestContext) {
	//获取参数
	queryType := c.Query("type")

	// 调用 service 层查询
	products, err := service.GetTypeList(queryType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || len(products) == 0 {
			c.JSON(http.StatusNotFound, utils.H{"message": "商品未找到"})
			return
		}
		c.JSON(http.StatusInternalServerError, utils.H{"message": "查询失败"})
		return
	}

	//查询成功
	c.JSON(http.StatusOK, utils.H{
		"info":     "success",
		"status":   10000,
		"message":  "查询成功",
		"type":     queryType,
		"products": products,
	})
}

// 给商品点赞
func ProductPraise(ctx context.Context, c *app.RequestContext) {
	productIdStr := c.PostForm("id")
	productID, _ := strconv.Atoi(productIdStr)

	if err := service.ProductPraiseService(ctx, productID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"message": "点赞失败", "err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"info":    "success",
		"status":  10000,
		"message": "点赞成功",
	})
}

// 推荐热度高的商品
func HeatProduct(ctx context.Context, c *app.RequestContext) {
	products, err := service.HeatProductService(ctx)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, utils.H{"message": "获取热门商品失败"})
		return
	}

	c.JSON(http.StatusOK, utils.H{"status": 10000, "info": "success", "products": products})

}
