//go:build go1.20

package utils

import (
	"github.com/gofiber/fiber/v2/utils"
)

// StringToBytes converts string to byte slice without a memory allocation.
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
// copy from https://github.com/gin-gonic/gin/blob/master/internal/bytesconv/bytesconv_1.20.go
// package of https://github.com/gofiber/fiber/blob/master/utils/convert_s2b_new.go
func StringToBytes(s string) []byte {
	if s == "" {
		return nil
	}
	return utils.UnsafeBytes(s)
}

// BytesToString converts byte slice to string without a memory allocation.
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
// copy from https://github.com/gin-gonic/gin/blob/master/internal/bytesconv/bytesconv_1.20.go
// package of https://github.com/gofiber/fiber/blob/master/utils/convert_b2s_new.go
func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return utils.UnsafeString(b)
}
