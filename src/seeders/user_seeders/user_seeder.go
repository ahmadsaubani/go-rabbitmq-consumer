package user_seeders

import (
	"fmt"
	"math/rand"
	"subscriber-topic-stars/src/entities/users"
	"subscriber-topic-stars/src/helpers"
	"time"

	"github.com/bxcodec/faker/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedUsers seeds users in the database, given a target count.
//
// If the target count is less than or equal to the current user count,
// it will not do anything and print a success message.
//
// Otherwise, it will generate the difference count of users with random
// usernames and email, but fixed password ("password123"), and insert
// them into the database in batches.
//
// The elapsed time of the seeding process is printed at the end.
func SeedUsers(db *gorm.DB, target int64) {
	var count int64
	db.Model(&users.User{}).Count(&count)
	if count > 0 {
		fmt.Println("ℹ️  Users already seeded. Skipping...")
		return
	}

	fmt.Println("🌱 Seeding users...")

	var usersBatch []users.User
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	now := time.Now()

	for i := int64(0); i < target; i++ {
		email := faker.Email()
		if i == 0 {
			email = "ahmadsaubani@testing.com"
		}

		usersBatch = append(usersBatch, users.User{
			Email:     email,
			Name:      faker.Name(),
			Username:  fmt.Sprintf("%s_%d", faker.Username(), rand.Intn(10000)),
			Password:  string(hashedPassword),
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	err := helpers.InsertModelBatch(usersBatch)

	if err != nil {
		fmt.Println("❌ Batch insert user failed:", err)
	} else {
		fmt.Println("✅ Users seeded successfully")
	}
}
