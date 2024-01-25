package env

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/betam/glb/lib/try"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Environment interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | bool | string
}

func init() {
	_ = godotenv.Load()
}

func Value[T Environment](key string, fallback T) (result T) {
	try.Catch(
		func() {
			result = NeedValue[T](key)
		},
		func(throwable error) {
			logrus.Tracef("%s, fallback to %v...", throwable, fallback)
			result = fallback
		},
	)

	return result
}

func Array[T Environment](key string, fallback []T) (result []T) {
	try.Catch(
		func() {
			result = NeedArray[T](key)
		},
		func(throwable error) {
			logrus.Tracef("%s, fallback to %v...", throwable, fallback)
			result = fallback
		},
	)

	return result
}

func NeedValue[T Environment](key string) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(errors.New(fmt.Sprintf("no env found: %s", key)))
	}

	var result T
	try.Catch(
		func() {
			result = convert[T](value)
		},
		func(throwable error) {
			panic(errors.New(fmt.Sprintf("error during parse environment: %s: %v", key, throwable)))
		},
	)

	return result
}

func NeedArray[T Environment](key string) []T {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(errors.New(fmt.Sprintf("no env found: %s", key)))
	}

	var result []T
	try.Catch(
		func() {
			re := regexp.MustCompile(" *, *")
			list := re.Split(value, -1)
			for _, element := range list {
				result = append(result, convert[T](element))
			}
		},
		func(throwable error) {
			panic(errors.New(fmt.Sprintf("error during parse array environment: %s: %v", key, throwable)))
		},
	)
	return result
}

func convert[T Environment](value string) T {
	var result T
	switch any(result).(type) {
	case string:
		return any(value).(T)
	case int:
		return any(int(readInt(value))).(T)
	case int64:
		return any(readInt(value)).(T)
	case int32:
		return any(int32(readInt(value))).(T)
	case int16:
		return any(int16(readInt(value))).(T)
	case int8:
		return any(int8(readInt(value))).(T)
	case uint:
		return any(uint(readUint(value))).(T)
	case uint64:
		return any(readUint(value)).(T)
	case uint32:
		return any(uint32(readUint(value))).(T)
	case uint16:
		return any(uint16(readUint(value))).(T)
	case uint8:
		return any(uint8(readUint(value))).(T)
	case float64:
		return any(readFloat(value)).(T)
	case float32:
		return any(float32(readFloat(value))).(T)
	case bool:
		return any(readBool(value)).(T)
	default:
		panic(errors.New("unexpected type"))
	}
}

func readInt(value string) int64 {
	return try.Throw(strconv.ParseInt(value, 10, 0))
}

func readUint(value string) uint64 {
	return try.Throw(strconv.ParseUint(value, 10, 0))
}

func readFloat(value string) float64 {
	return try.Throw(strconv.ParseFloat(value, 64))
}

func readBool(value string) bool {
	switch value {
	case "1", "true":
		return true
	case "0", "false":
		return false
	default:
		panic(errors.New(fmt.Sprintf(`non-bool expression: "%s"`, value)))
	}
}
