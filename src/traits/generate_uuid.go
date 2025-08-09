package traits

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GenerateUUIDStruct struct {
	UUID string `gorm:"type:uuid;uniqueIndex" json:"uuid"`
}

func (u *GenerateUUIDStruct) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UUID == "" {
		u.UUID = GenerateUUID()
	}
	return
}

func GenerateUUID() string {
	return uuid.New().String()
}

func SetUUIDForStruct(model interface{}) error {
	val := reflect.ValueOf(model).Elem()
	// Pastikan model adalah pointer ke struct
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got %s", val.Kind())
	}

	// Cari field dengan nama "UUID"
	uidField := val.FieldByName("UUID")
	if uidField.IsValid() && uidField.CanSet() && uidField.Kind() == reflect.String {
		// Set UUID jika field UUID ditemukan
		uidField.SetString(GenerateUUID())
	}

	return nil
}
