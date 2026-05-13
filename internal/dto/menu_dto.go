package dto

type CreateMenuRequest struct {
	CategoryID  string `form:"category_id" binding:"required"`
	Name        string `form:"name" binding:"required,min=2,max=100"`
	Price       int    `form:"price" binding:"required,gte=0"`
	Stock       int    `form:"stock" binding:"gte=0"`
	IsAvailable bool   `form:"is_available"`
}

type UpdateMenuRequest struct {
	CategoryID  string `form:"category_id" binding:"required"`
	Name        string `form:"name" binding:"required,min=2,max=100"`
	Price       int    `form:"price" binding:"required,gte=0"`
	Stock       int    `form:"stock" binding:"gte=0"`
	IsAvailable bool   `form:"is_available"`
}

type CategoryResponseMini struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type MenuResponse struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Price       int                   `json:"price"`
	Stock       int                   `json:"stock"`
	ImageURL    *string               `json:"image_url,omitempty"`
	IsAvailable bool                  `json:"is_available"`
	Category    *CategoryResponseMini `json:"category,omitempty"`
}