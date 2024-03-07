package entity

type Product struct {
	ID       int64
	Name     string
	Category Category
}

// type ProductView struct {
// 	ID           int64
// 	Name         string
// 	CategoryName string
// }

type ProductCategoryListItem struct {
	ID   int64
	Name string
}

type AddProductDTO struct {
	ProductName string
	CategoryID  int64
}

type AddOrUpdateProductDTO struct {
	ProductName  string `json:"title"`
	CategoryName string `json:"category"`
}

type UpdateProductNameDTO struct {
	ProductID int64
	NewName   string
}

type UpdateProductCategoryDTO struct {
	ProductID     int64
	OldCategoryID int64
	NewCategoryID int64
}
