package api

import (
	"context"
	"database/sql"
	"decleration/model"
	"decleration/service"
	"decleration/utils1"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
)

// 查询商品评论
func GetComment(ctx context.Context, c *app.RequestContext) {
	//查询商品id
	productID := c.DefaultQuery("product_id", "")

	//获取评论
	com, err := service.GetCommentService(ctx, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, utils.H{"message": "没有找到该商品或该商品没有评论"})
			return
		}
		c.JSON(http.StatusInternalServerError, utils.H{"message": "数据解析失败"})
		return
	}

	//查询成功
	c.JSON(http.StatusOK, utils.H{
		"info":     "success",
		"status":   10000,
		"comments": com,
	})
}

// 发送评论
func PostComment(ctx context.Context, c *app.RequestContext) {
	//提取解析token
	auth, _ := utils1.ExtractToken(c)
	userID, _ := utils1.ValidateToken(auth)

	//获取解析数据
	var com struct {
		ProductID int    `json:"product_id"`
		Content   string `json:"content"`
		CommentID string `json:"comment_id"`
	}
	if err := c.Bind(&com); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{"message": "解析失败"})
		return
	}

	// 组装 Comment 结构体
	comment := model.Comment{
		CommentID:   com.CommentID,
		Content:     com.Content,
		UserID:      fmt.Sprintf("%v", userID), // 确保 UserID 是 string
		Nickname:    "julia",                   // 从 token 解析
		PraiseCount: 0,                         // 新建评论默认 0 赞
		IsPraised:   0,                         // 默认未点赞
		ProductID:   com.ProductID,
	}

	//进行评论
	if err := service.PostCommentService(ctx, comment); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, utils.H{"message": "评论失败"})
		return
	}

	//评论成功
	c.JSON(http.StatusOK, utils.H{
		"info":    "success",
		"status":  10000,
		"message": "评论成功",
		"comment": comment})
}

// 删除评论
func DeleteComment(ctx context.Context, c *app.RequestContext) {
	//获取，解析token
	auth, _ := utils1.ExtractToken(c)
	utils1.ValidateToken(auth)

	//获取数据,删除
	commentID := c.DefaultQuery("comment_id", "")
	if err := service.DeleteCommentService(ctx, commentID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"message": "删除失败"})
		return
	}

	//操作成功
	c.JSON(http.StatusOK, utils.H{
		"info":    "success",
		"status":  10000,
		"message": "删除成功"})
}

// 更新评论
func UpdateComment(ctx context.Context, c *app.RequestContext) {
	//提取token
	auth, _ := utils1.ExtractToken(c)
	//验证token
	_, err := utils1.ValidateToken(auth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.H{"message": "无效的Token"})
		return
	}

	// 获取,解析参数
	var request struct {
		CommentID string `json:"comment_id"`
		Content   string `json:"content"`
	}
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"message": "数据解析失败"})
		return
	}

	// 校验参数
	if request.CommentID == "" || request.Content == "" {
		c.JSON(http.StatusBadRequest, utils.H{"message": "comment_id/评论内容不能为空"})
		return
	}

	// 更新评论
	if err := service.UpdateCommentService(ctx, request.CommentID, request.Content); err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"message": "更新评论失败"})
		return
	}

	//更改评论成功
	c.JSON(http.StatusOK, utils.H{
		"info":    "success",
		"status":  10000,
		"message": "修改成功"})
}

// 给评论点赞/点踩
func CommentPraise(ctx context.Context, c *app.RequestContext) {
	utils1.ExtractToken(c)

	// 解析 form-data 参数
	var request struct {
		CommentID string `json:"comment_id"`
		Status    string `json:"status"`
	}
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{"message": "数据解析失败"})
		return
	}
	// 确保 CommentID 不能为空
	if request.CommentID == "" || request.Status == "" {
		c.JSON(http.StatusBadRequest, utils.H{"message": "comment_id/status 不能为空"})
		return
	}

	if err := service.CommentPraiseService(ctx, request.CommentID, request.Status); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{"message": "操作失败"})
		return
	}

	//操作成功
	c.JSON(http.StatusOK, utils.H{
		"info":    "success",
		"status":  10000,
		"message": "点赞成功"})
}
