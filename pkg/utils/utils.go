package utils

import (
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/mail"
	"reflect"
	"strconv"
	"time"

	"github.com/tandy9527/js-slot/pkg/consts"

	"golang.org/x/crypto/bcrypt"
)

// GenerateCode generates a random numeric code of specified length
func GenerateCode(length int) string {
	code := ""
	for range length {
		n := rand.Intn(10) // 0-9
		code += fmt.Sprintf("%d", n)
	}
	return code
}

// 获取当前10位时间戳（秒级）
func CurrentTimestamp() int64 {
	return time.Now().Unix()
}

// HashPassword 使用 bcrypt 加密密码
func HashPassword(password string) (string, error) {
	// bcrypt 默认会生成随机盐
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

// VerifyPassword 校验密码
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// StructToMap 把任意 struct 转成 map[string]interface{}，使用 json tag
func StructToMap(s any) map[string]any {
	result := make(map[string]any)
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		tag := fieldType.Tag.Get("json")
		if tag == "" || tag == "-" {
			tag = fieldType.Name
		}

		result[tag] = field.Interface()
	}

	return result
}

// GetUserRedisKey 根据用户 UID 生成 Redis Key
func GetUserRedisKey(uid int64) string {

	return fmt.Sprintf("%s%d", consts.REDIS_USER_KEY, uid)
}

// GenerateJTI 生成安全随机 JTI（32 位十六进制字符串）
func GenerateJTI() string {
	b := make([]byte, 16) // 16 字节 = 128 位
	crand.Read(b)
	return hex.EncodeToString(b)
}

// NewTraceID 生成短小唯一 TraceID
func NewTraceID() string {

	// 当前时间戳（毫秒）
	t := time.Now().UnixNano() / 1e6

	// 随机数（0~999999）
	r := rand.Intn(1000000)

	str := make([]byte, 0, 32)
	str = strconv.AppendInt(str, t, 10)
	str = append(str, '-')
	str = strconv.AppendInt(str, int64(r), 10)

	return string(str)
}

// IsValidEmail 校验邮箱格式
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// IsEmpty 判断字符串是否为空或仅包含空白字符
func IsEmpty(s string) bool {
	return len(s) == 0
}
func NotEmpty(s string) bool {
	return !IsEmpty(s)
}

// IsBlank 判断字符串是否为空，或者只包含空格、换行、制表符等
func IsBlank(s string) bool {
	for _, r := range s {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			return false
		}
	}
	return true
}

// ConvertStruct 自动将 src 结构体的同名字段复制到 dst
// dst 必须是指针类型
func ConvertStruct(src interface{}, dst any) {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	if dstVal.Kind() != reflect.Ptr || dstVal.Elem().Kind() != reflect.Struct {
		panic("dst must be a pointer to struct")
	}

	srcVal = reflect.Indirect(srcVal)
	dstVal = reflect.Indirect(dstVal)

	for i := 0; i < dstVal.NumField(); i++ {
		dstField := dstVal.Type().Field(i)
		dstFieldVal := dstVal.Field(i)

		srcFieldVal := srcVal.FieldByName(dstField.Name)
		if srcFieldVal.IsValid() && dstFieldVal.CanSet() {
			dstFieldVal.Set(srcFieldVal)
		}
	}
}
