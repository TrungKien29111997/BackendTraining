package repository

type UserRepository interface {
	Create()
	FindByUUID()
	Update()
	Delete()
	FillAll()
}
