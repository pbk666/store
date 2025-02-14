package dao

import (
	"context"
)

func GetAverageRating(ctx context.Context, productID int) (float64, error) {
	db, err := InitDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var avgRating float64
	query := "SELECT COALESCE(AVG(score), 0) FROM rating_product WHERE product_id = ?"
	err = db.QueryRow(query, productID).Scan(&avgRating)
	if err != nil {
		return 0, err
	}
	return avgRating, nil
}
func UpdateProductHeat(productID int) error {
	db, _ := InitDB()
	query := `UPDATE products SET heat = ( purchases * 5 + praise_count * 3 + comment_count * 2 + rating * 10)WHERE id = ?`
	_, err := db.Exec(query, productID)
	return err
}
