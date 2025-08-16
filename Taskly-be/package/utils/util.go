package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func PtrIntIfValid(n sql.NullInt32) *int32 {
	if n.Valid {
		return &n.Int32
	}
	return nil
}
func ToNullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{
			Time:  *t,
			Valid: true,
		}
	}
	return sql.NullTime{
		Valid: false,
	}
}
func ToNullInt32(i *int32) sql.NullInt32 {
	if i != nil {
		return sql.NullInt32{Int32: *i, Valid: true}
	}
	return sql.NullInt32{Valid: false}
}

func ToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func ToSliceNullString(arr *[]string) []sql.NullString {
	if arr == nil {
		return nil
	}
	result := make([]sql.NullString, len(*arr))
	for i, v := range *arr {
		if v != "" {
			result[i] = sql.NullString{String: v, Valid: true}
		} else {
			result[i] = sql.NullString{Valid: false}
		}
	}
	return result
}

func PtrTimeIfValid(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

func PtrSliceIfValid(slice []string) *[]string {
	if len(slice) > 0 {
		return &slice
	}
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // DefaultCost = 10
	return string(bytes), err
}

// CheckPasswordHash kiểm tra mật khẩu gốc với chuỗi hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserKey(hashKey string) string {
	return fmt.Sprintf("u:%s:otp", hashKey)
}
func GenerateCliTokenUUID(userId int) string {
	newUUID := uuid.New()
	uuidString := strings.ReplaceAll((newUUID).String(), "", "")
	return strconv.Itoa(userId) + "clitoken" + uuidString
}

func PtrStringIfValid(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
