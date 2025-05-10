package usecase

type UseCases struct {
	Auth     Auth
	Cart     Cart
	Order    Order
	Producer Producer
	Product  Product
	User     User
}

func NewUseCases(a Auth, c Cart, o Order, p Producer, pr Product, u User) *UseCases {
	return &UseCases{
		Auth:     a,
		Cart:     c,
		Order:    o,
		Producer: p,
		Product:  pr,
		User:     u,
	}
}
