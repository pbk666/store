package service

import (
	"context"
	"decleration/dao"
	"decleration/model"
	"fmt"
)

// 下单
func CreateOrder(ctx context.Context, request model.Order) (int64, error) {
	// 连接数据库
	db, err := dao.InitDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// 开启事务
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// 创建订单
	query := "INSERT INTO orders (user_id, address, total) VALUES (?, ?, ?)"
	result, err := tx.Exec(query, request.UserID, request.Address, request.Total)
	if err != nil {
		tx.Rollback()
		fmt.Println("下单失败:", err)
		return 0, err
	}

	// 获取订单号
	orderID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		fmt.Println("生成订单号失败:", err)
		return 0, err
	}

	// 插入订单详情前，检查 product_id 是否存在
	checkProductQuery := "SELECT COUNT(*) FROM products WHERE id = ?"
	for _, item := range request.Orders {
		var count int
		err = db.QueryRow(checkProductQuery, item.ProductID).Scan(&count)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if count == 0 {
			tx.Rollback()
			return 0, fmt.Errorf("商品ID %d 不存在", item.ProductID)
		}
	}

	// 插入订单详情
	orderItemQuery := "INSERT INTO order_items (order_id, product_id, quantity) VALUES (?, ?, ?)"
	for _, item := range request.Orders {
		_, err = tx.Exec(orderItemQuery, orderID, item.ProductID, item.Quantity)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// 批量更新商品购买数
	productUpdateQuery := "UPDATE products SET purchases = purchases + ? WHERE id = ?"
	for _, item := range request.Orders {
		_, err = tx.Exec(productUpdateQuery, item.Quantity, item.ProductID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	// **在事务提交后调用 `UpdateProductHeat`**
	go func() {
		for _, item := range request.Orders {
			if err := dao.UpdateProductHeat(item.ProductID); err != nil {
				fmt.Printf("更新商品 %d 热度失败: %v\n", item.ProductID, err)
			}
		}
	}()

	return orderID, nil
}

// 评分
func RatingProduct(ctx context.Context, review model.Rating) error {
	db, err := dao.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 检查 product_id 是否存在
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", review.ProductID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("商品ID %d 不存在", review.ProductID)
	}

	// 插入评分
	query := "INSERT INTO rating_product (user_id, product_id, score ) VALUES (?, ?, ?)"
	_, err = db.Exec(query, review.UserID, review.ProductID, review.Score)
	if err != nil {
		return err
	}

	return nil
}

// 更新商品评分
func UpdateProductRating(ctx context.Context, productID int) error {
	db, err := dao.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	avgRating, _ := dao.GetAverageRating(ctx, productID)
	query := "UPDATE products SET rating = ? where id = ?"
	_, err = db.Exec(query, avgRating, productID)
	if err != nil {
		return err
	}

	// 更新热度
	err = dao.UpdateProductHeat(productID)
	if err != nil {
		return fmt.Errorf("更新热度失败: %w", err)
	}

	return nil
}
