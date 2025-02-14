package api

import (
	"context"
	"decleration/model"
	"decleration/service"
	"decleration/utils1"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
)

func Order(ctx context.Context, c *app.RequestContext) {
	//提取token
	utils1.ExtractToken(c)

	//解析参数
	var request model.Order
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{"message": "数据解析失败"})
		return
	}

	//下单,获取订单号
	OrderID, err := service.CreateOrder(ctx, request)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, utils.H{"message": "下单失败"})
		return
	}

	//下单成功
	c.JSON(http.StatusOK, utils.H{
		"info":     "success",
		"status":   10000,
		"message":  "下单成功",
		"order_id": OrderID,
	})
}

// 给商品评分
func Rating(ctx context.Context, c *app.RequestContext) {
	var request model.Rating
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{"message": "数据解析失败"})
		return
	}

	if request.Score < 1 || request.Score > 5 {
		c.JSON(http.StatusBadRequest, utils.H{"message": "评分必须是1到5"})
		return
	}

	if err := service.RatingProduct(ctx, request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, utils.H{"message": "评分失败"})
		return
	}
	if err := service.UpdateProductRating(ctx, request.ProductID); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, utils.H{"message": "刷新商品评分失败"})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"info":   "success",
		"status": 10000,
	})
}
