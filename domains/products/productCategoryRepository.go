package products

import (
	"fmt"

	. "github.com/SofyanHadiA/linq/core/database"
	. "github.com/SofyanHadiA/linq/core/repository"
	"github.com/SofyanHadiA/linq/core/utils"

	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

type productCategoryRepository struct {
	db IDB
}

func NewProductCategoryRepository(db IDB) IRepository {
	return productCategoryRepository{
		db: db,
	}
}

func (repo productCategoryRepository) CountAll() (int, error) {
	countQuery := "SELECT COUNT(*) FROM product_categories WHERE deleted = 0"

	var result int
	row, err := repo.db.ResolveSingle(countQuery)
	row.Scan(&result)
	if err != nil {
		return -1, err
	}
	return result, err
}

func (repo productCategoryRepository) IsExist(id uuid.UUID) (bool, error) {
	isExistQuery := "SELECT EXISTS(SELECT * FROM product_categories WHERE uid=? AND deleted = 0)"

	var result bool
	row, err := repo.db.ResolveSingle(isExistQuery, id)
	row.Scan(&result)
	return result, err
}

func (repo productCategoryRepository) GetAll(paging utils.Paging) (IModels, error) {
	query := "SELECT * FROM product_categories WHERE deleted=0 "

	if paging.Keyword != "" {
		query += ` AND (title LIKE '%?%') `
	}

	var columnMap string
	switch paging.Order {
	case 1:
		columnMap = "title"
	case 2:
		columnMap = "slug"
	case 3:
		columnMap = "description"
	default:
		columnMap = "created"
	}

	query += fmt.Sprintf(" ORDER BY %s %s ", columnMap, paging.OrderDir)

	if paging.Length > 0 {
		query += fmt.Sprintf(" LIMIT %d ", paging.Length)
	} else {
		query += " LIMIT 25 "
	}

	rows := &sqlx.Rows{}
	var err error

	if paging.Keyword != "" {
		rows, err = repo.db.Resolve(query, paging.Keyword)
	} else {
		rows, err = repo.db.Resolve(query)
	}
	if err != nil {
		return nil, err
	}

	result := ProductCategories{}

	for rows.Next() {
		var productCategory = &ProductCategory{}
		err := rows.StructScan(&productCategory)

		if err != nil {
			return nil, err
		}

		result = append(result, (*productCategory))
	}

	return &result, err
}

func (repo productCategoryRepository) Get(id uuid.UUID) (IModel, error) {
	selectQuery := "SELECT * FROM product_categories WHERE uid = ? AND deleted= 0 "

	productCategory := &ProductCategory{}
	rows, err := repo.db.ResolveSingle(selectQuery, id)
	if err != nil {
		return nil, err
	}
	rows.StructScan(productCategory)

	return productCategory, err
}

func (repo productCategoryRepository) Insert(model IModel) error {

	insertQuery := `INSERT INTO product_categories(uid, title, description, slug, created) VALUES 
		(:uid, :title, :description, :slug, now())`

	productCategory := model.(*ProductCategory)
	productCategory.Uid = uuid.NewV4()

	_, err := repo.db.Execute(insertQuery, productCategory)

	return err
}

func (repo productCategoryRepository) Update(model IModel) error {
	updateQuery := `UPDATE product_categories SET 
		title=:title, description=:description, slug=:slug, updated=now() WHERE uid=:uid`

	_, err := repo.db.Execute(updateQuery, model)

	return err
}

func (repo productCategoryRepository) Delete(model IModel) error {
	deleteQuery := "UPDATE product_categories SET deleted=1 WHERE uid=:uid"

	_, err := repo.db.Execute(deleteQuery, model)

	return err
}

func (repo productCategoryRepository) DeleteBulk(productCategories []uuid.UUID) error {
	deleteQuery := "UPDATE product_categories SET deleted=1 WHERE uid IN(?)"

	_, err := repo.db.ExecuteBulk(deleteQuery, productCategories)

	return err
}
