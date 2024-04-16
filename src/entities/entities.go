package entities

type User struct {
	ID        uint    `gorm:"autoIncrement;primaryKey"`
	Nickname  string  `json:"nickname" validate:"required,gt=0" binding:"required,gt=0"`
	Email     string  `json:"email" validate:"required,gt=0,email" binding:"required,gt=0,email"`
	Phone     string  `json:"phone" validate:"required,gt=0" binding:"required,gt=0"`
	Password  string  `json:"password" validate:"required,gt=0" binding:"required,gt=0"`
	FirstName *string `json:"firstName"`
	Surname   *string `json:"surname"`
	LastName  *string `json:"lastName"`
	Work      *string `json:"work"`
	Study     *string `json:"study"`
	Telegram  *string `json:"telegram"`
}