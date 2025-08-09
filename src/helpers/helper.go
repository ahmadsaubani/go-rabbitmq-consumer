package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	database "subscriber-topic-stars/src/configs"
	"subscriber-topic-stars/src/traits"
	"subscriber-topic-stars/src/utils/filters"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const maxBatchSize = 500
const layout = "2006-01-02"

func StringToDateOnly(input string) time.Time {
	t, _ := time.Parse(layout, input)
	return t
}
func InsertModelBatch[T any](models []T) error {
	if len(models) == 0 {
		return nil
	}

	if database.GormDB == nil {
		return fmt.Errorf("❌ Database connection is not initialized")
	}

	for start := 0; start < len(models); start += maxBatchSize {
		end := start + maxBatchSize
		if end > len(models) {
			end = len(models)
		}
		batch := models[start:end]

		// Set UUID untuk setiap item dalam batch
		for i := range batch {
			err := traits.SetUUIDForStruct(&batch[i])
			if err != nil {
				return fmt.Errorf("❌ Error setting UUID: %w", err)
			}
		}

		err := database.GormDB.Transaction(func(tx *gorm.DB) error {

			if err := tx.Create(&batch).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("❌ GORM transaction failed: %w", err)
		}
	}

	return nil
}

func InsertModel[T any](model *T) error {
	if err := traits.SetUUIDForStruct(model); err != nil {
		return fmt.Errorf("❌ Error setting UUID: %w", err)
	}
	return database.GormDB.Create(model).Error
}

func GetAllModelsWithDB[T any](ctx *gin.Context, db *gorm.DB, models *[]T) error {
	orderBy := ctx.DefaultQuery("order_by", "")

	if orderBy != "" {
		orderParts := strings.Split(orderBy, ",")
		if len(orderParts) == 2 {
			orderBy = fmt.Sprintf("%s %s", orderParts[0], orderParts[1])
		}
	}

	if db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	query := db

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			query = query.Limit(limit)
		}
	}

	if offsetStr := ctx.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			query = query.Offset(offset)
		}
	}

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	whereClause, args, err := filters.BuildFilters(ctx)
	if err != nil {
		return err
	}
	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}
	return db.Find(models).Error
}

func GetAllModels[T any](ctx *gin.Context, models *[]T) error {
	orderBy := ctx.DefaultQuery("order_by", "")

	if orderBy != "" {
		orderParts := strings.Split(orderBy, ",")
		if len(orderParts) == 2 {
			orderBy = fmt.Sprintf("%s %s", orderParts[0], orderParts[1])
		}
	}

	if database.GormDB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	query := database.GormDB

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			query = query.Limit(limit)
		}
	}

	if offsetStr := ctx.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			query = query.Offset(offset)
		}
	}

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	whereClause, args, err := filters.BuildFilters(ctx)
	if err != nil {
		return err
	}
	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	return query.Find(models).Error
}
func GettingAllModels[T any](models *[]T, preloads []string, conditions ...any) error {

	if database.GormDB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	query := database.GormDB

	if len(conditions)%2 != 0 {
		return fmt.Errorf("conditions must be in key-value pairs")
	}
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	for i := 0; i < len(conditions); i += 2 {
		field := conditions[i].(string)
		value := conditions[i+1]
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	return query.Find(models).Error
}

func GetModelByID[T any](model *T, id any) error {
	return database.GormDB.First(model, id).Error
}

func UpdateModelByIDWithMap[T any](updatedFields map[string]interface{}, id any) error {
	return database.GormDB.Model(new(T)).Where("id = ?", id).Updates(updatedFields).Error
}

func UpdateModelByID[T any](model *T, id any) error {
	return database.GormDB.Model(model).Where("id = ?", id).Updates(model).Error

}

func DeleteModelByID[T any](model *T, id any) error {
	return database.GormDB.Delete(model, id).Error
}

func FindOneByField[T any](model *T, conditions ...any) error {
	return FindOneByFieldWithPreload(model, nil, conditions...)
}

func FindOneByFieldWithPreload[T any](model *T, preloads []string, conditions ...any) error {
	if len(conditions)%2 != 0 {
		return fmt.Errorf("conditions must be in key-value pairs")
	}

	query := database.GormDB
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	for i := 0; i < len(conditions); i += 2 {
		field := conditions[i].(string)
		value := conditions[i+1]
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	return query.First(model).Error
}

func CountModel[T any]() (int64, error) {
	var total int64
	err := database.GormDB.Model(new(T)).Count(&total).Error
	return total, err

}

func CountModelWithFilter[T any](filters func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	err := filters(database.GormDB.Model(new(T))).Count(&total).Error
	return total, err
}
func GenerateRandomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate random bytes")
	}
	return hex.EncodeToString(b)
}

func ParseAndValidateToken(tokenStr string) (jwt.MapClaims, error) {
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return nil, fmt.Errorf("token expired")
		}
	}

	return claims, nil
}

func ConvertTokenToUserId(token map[string]interface{}) (uint64, error) {
	tokenRaw, ok := token["token"]
	if !ok {
		return 0, fmt.Errorf("missing token")
	}

	tokenStr := fmt.Sprintf("%v", tokenRaw)

	// Validasi token
	claims, err := ParseAndValidateToken(tokenStr)
	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user_id in token : %v", claims["user_id"])
	}
	userID := uint64(userIDFloat)

	return userID, nil
}
