package domain

type TweetStore interface {
	GetAll() (*[]Tweet, error)
	AddProduct(t *Tweet)
	PutProduct(t *Tweet, id int) error
	DeleteProduct(id int) error
}
