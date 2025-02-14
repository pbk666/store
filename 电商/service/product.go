package service

import (
	"context"
	"database/sql"
	"decleration/dao"
	"decleration/model"
	"errors"
	"fmt"
	"strconv"
)

// 获取产品列表的业务逻辑
func GetProductList(ctx context.Context) ([]model.Product, error) {
	// 连接数据库
	db, err := dao.InitDB()
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	defer db.Close()

	// 查询商品数据
	rows, err := db.Query("SELECT id, name, description, type, price FROM products")
	if err != nil {
		return nil, fmt.Errorf("查询失败: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Type, &p.Price); err != nil {
			return nil, fmt.Errorf("数据解析失败: %w", err)
		}
		products = append(products, p)
	}

	return products, nil
}

// 查询商品（支持模糊搜索）
func SearchProduct(ctx context.Context, searchQuery string) ([]model.Product, error) {
	// 连接数据库
	db, err := dao.InitDB()
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	defer db.Close()

	query := `SELECT id, name, description, type, price FROM products WHERE name LIKE ? OR type LIKE ?`
	rows, err := db.Query(query, "%"+searchQuery+"%", "%"+searchQuery+"%")
	if err != nil {
		return nil, fmt.Errorf("查询失败: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Type, &p.Price); err != nil {
			return nil, fmt.Errorf("数据解析失败: %w", err)
		}
		products = append(products, p)
	}

	return products, nil
}

func AddCartService(ctx context.Context, userID string, idStr string) (model.Product, error) {
	var p model.Product

	// 检查 ID 是否为空
	if idStr == "" {
		return p, errors.New("id 不能为空")
	}

	// 转换 ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return p, errors.New("商品 ID 必须是整数")
	}

	// 连接数据库
	db, err := dao.InitDB()
	if err != nil {
		return p, errors.New("数据库连接失败")
	}
	defer db.Close()

	// 查询商品是否存在
	err = db.QueryRow("SELECT id, name, description, type, price FROM products WHERE id=?", id).
		Scan(&p.ID, &p.Name, &p.Description, &p.Type, &p.Price)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, errors.New("商品不存在")
		}
		return p, err
	}

	// 查询购物车是否已有该商品
	var cartCount int
	err = db.QueryRow("SELECT COUNT(*) FROM cart WHERE user_id = ? AND product_id = ?", userID, id).
		Scan(&cartCount)

	if err != nil {
		return p, errors.New("查询购物车失败")
	}

	if cartCount > 0 {
		_, err = db.Exec("UPDATE cart SET quantity = quantity + 1 WHERE user_id = ? AND product_id = ?", userID, id)
		if err != nil {
			return p, errors.New("更新购物车失败")
		}
	} else {
		_, err = db.Exec("INSERT INTO cart (user_id, product_id, quantity) VALUES (?, ?, ?)", userID, id, 1)
		if err != nil {
			return p, errors.New("添加商品到购物车失败")
		}
	}

	return p, nil
}

func GetCartList(userID string) ([]model.Cart, error) {
	// 连接数据库
	db, err := dao.InitDB()
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	defer db.Close()

	// 查询购物车数据
	rows, err := db.Query("SELECT id, user_id, product_id, quantity FROM cart WHERE user_id=?", userID)
	if err != nil {
		return nil, fmt.Errorf("查询购物车失败: %w", err)
	}
	defer rows.Close()

	// 遍历查询结果
	var cart []model.Cart
	for rows.Next() {
		var cartOne model.Cart
		if err = rows.Scan(&cartOne.ID, &cartOne.UserID, &cartOne.ProductID, &cartOne.Quantity); err != nil {
			return nil, fmt.Errorf("数据解析失败: %w", err)
		}
		cart = append(cart, cartOne)
	}

	return cart, nil
}

func GetProductInfo(productID string) (*model.Product, error) {
	// 连接数据库
	db, err := dao.InitDB()
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	defer db.Close()

	// 查询商品信息
	row := db.QueryRow("SELECT id, name, description, type, price FROM products WHERE id=?", productID)
	var product model.Product
	if err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Type, &product.Price); err != nil {
		return nil, err
	}

	return &product, nil
}

func GetTypeList(queryType string) ([]model.Product, error) {
	// 连接数据库
	db, err := dao.InitDB()
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	defer db.Close()

	// 查询商品列表
	rows, err := db.Query("SELECT id, name, description, type, price FROM products WHERE type=?", queryType)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %w", err)
	}
	defer rows.Close()

	// 解析查询结果
	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Type, &p.Price); err != nil {
			return nil, fmt.Errorf("数据解析失败: %w", err)
		}
		products = append(products, p)
	}

	return products, nil
}

// 给商品点赞
func ProductPraiseService(ctx context.Context, productID int) error {
	db, err := dao.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE products SET praise_count=praise_count + 1 WHERE id=?", productID)
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

// 推荐热度高的商品
func HeatProductService(ctx context.Context) ([]model.Product, error) {
	db, err := dao.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id,name,description,price FROM products ORDER BY heat DESC LIMIT 5")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	var product model.Product
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
