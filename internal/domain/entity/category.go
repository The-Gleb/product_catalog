package entity

type Category struct {
	ID   int64
	Name string
}

type AddCategoryDTO struct {
	Name string
}

type UpdateCategoryNameDTO struct {
	CategoryID int64
	NewName    string
}
