package seeders

import (
	"subscriber-topic-stars/src/configs"
	"subscriber-topic-stars/src/seeders/user_seeders"
)

func Run(db *configs.DBConnection) {
	user_seeders.SeedUsers(db.Gorm, 5000)
}
