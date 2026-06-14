package user

type Repository interface {
	Create(entity *User) error
	FindByID(id uint64) (*User, error)
}
