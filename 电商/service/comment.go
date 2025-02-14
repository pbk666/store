package service

import (
	"context"
	"database/sql"
	"decleration/dao"
	"decleration/model"
	"errors"
	"fmt"
)

// 获取商品评论
func GetCommentService(ctx context.Context, productID string) (*model.Comment, error) {
	db, err := dao.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT comment_id, content, user_id, nickname, praise_count, is_praised, product_id FROM comment WHERE product_id=?", productID)
	var com model.Comment
	if err := row.Scan(&com.CommentID, &com.Content, &com.UserID, &com.Nickname, &com.PraiseCount, &com.IsPraised, &com.ProductID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("没有找到该商品或该商品没有评论")
		}
		return nil, errors.New("数据解析失败")
	}
	return &com, nil
}

// 发送评论
func PostCommentService(ctx context.Context, comment model.Comment) error {
	//连接到数据库
	db, err := dao.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	//开启事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	//插入评论
	query := "INSERT INTO comment ( content, user_id, nickname, praise_count, is_praised, product_id) VALUES ( ?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, comment.Content, comment.UserID, comment.Nickname, comment.PraiseCount, comment.IsPraised, comment.ProductID)
	if err != nil {
		tx.Rollback()
		return err
	}

	//更新商品评论数
	updateQuery := "UPDATE products SET comment_count = comment_count + 1 WHERE id=?"
	if _, err := db.Exec(updateQuery, comment.ProductID); err != nil {
		tx.Rollback()
		return err
	}

	// 更新热度
	err = dao.UpdateProductHeat(comment.ProductID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("更新热度失败: %w", err)
	}

	//提交事务
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func DeleteCommentService(ctx context.Context, commentID string) error {
	db, err := dao.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	return dao.DeleteComment(db, commentID)
}

func UpdateCommentService(ctx context.Context, commentID, content string) error {
	db, err := dao.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	result, err := db.Exec("UPDATE comment SET content=? WHERE comment_id=?", content, commentID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("评论不存在或内容未更改")
	}
	return nil
}

func CommentPraiseService(ctx context.Context, commentID, status string) error {
	db, err := dao.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	var isPraised int
	row := db.QueryRow("SELECT is_praised FROM comment WHERE comment_id=?", commentID)
	if err := row.Scan(&isPraised); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("评论不存在")
		}
		return err
	}

	if isPraised != 0 {
		return errors.New("您已经点赞或点踩")
	}

	_, err = db.Exec("UPDATE comment SET is_praised=? WHERE comment_id=?", status, commentID)
	if err != nil {
		return err
	}
	return nil
}
