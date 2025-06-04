package entity

type Cart struct {
	Product *Product `json:"product"`
	Count   int      `json:"count"`
}
