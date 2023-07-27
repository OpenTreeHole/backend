package accountV2

import (
	"github.com/gofiber/fiber/v2"

	_ "github.com/opentreehole/backend/pkg/utils"
)

// Register godoc
//
// @Summary Register
// @Description Register with email, password and optional verification code if enabled
// @Tags Account
// @Accept json
// @Produce json
// @Router /api/register [post]
// @Param json body RegisterRequest true "RegisterRequest"
// @Success 201 {object} TokenResponse
// @Failure 400 {object} utils.MessageResponse "验证码错误，用户已存在"
func Register(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// RegisterDebug godoc
//
// @Summary register, debug only
// @Description register with email, password, not need verification code
// @Tags Account
// @Accept json
// @Produce json
// @Router /api/debug/register [post]
// @Param json body LoginRequest true "json"
// @Success 201 {object} TokenResponse
// @Failure 400 {object} utils.MessageResponse "用户已注册"
// @Failure 500 {object} utils.MessageResponse
// @Security ApiKeyAuth
func RegisterDebug(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// RegisterDebugInBatch godoc
//
// @Summary register in batch, debug only
// @Description register with email, password, not need verification code
// @Tags Account
// @Accept json
// @Produce json
// @Router /api/debug/register/_batch [post]
// @Param json body RegisterInBatchRequest true "RegisterInBatchRequest"
// @Success 201 {object} TokenResponse
// @Failure 400 {object} utils.MessageResponse "用户已注册"
// @Failure 500 {object} utils.MessageResponse
// @Security ApiKeyAuth
func RegisterDebugInBatch(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// Login godoc
//
// @Summary Login
// @Description Login with email and password, return jwt token, no need jwt
// @Tags Account
// @Accept json
// @Produce json
// @Router /api/login [post]
// @Param json body LoginRequest true "LoginRequest"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} utils.MessageResponse
// @Failure 404 {object} utils.MessageResponse
// @Failure 500 {object} utils.MessageResponse
func Login(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// Logout godoc
//
// @Summary Logout
// @Description Logout, need jwt
// @Tags Account
// @Produce json
// @Router /api/logout [post]
// @Success 200 {object} utils.MessageResponse
// @Failure 400 {object} utils.MessageResponse
// @Failure 404 {object} utils.MessageResponse
// @Failure 500 {object} utils.MessageResponse
func Logout(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// Refresh godoc
//
// @Summary Refresh
// @Description Refresh jwt token, need jwt
// @Tags Account
// @Accept json
// @Produce json
// @Router /api/refresh [post]
// @Success 200 {object} TokenResponse
func Refresh(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// ResetPassword godoc
//
// @Summary reset password
// @Description reset user password and jwt credential
// @Tags Account
// @Accept json
// @Produce json
// @Router /api/register [put]
// @Param json body RegisterRequest true "RegisterRequest"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} utils.MessageResponse "验证码错误"
// @Failure 500 {object} utils.MessageResponse
func ResetPassword(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// VerifyWithEmailOld godoc
//
// @Summary verify with email in path
// @Description verify with email in path, send verification email
// @Deprecated
// @Tags Account
// @Produce json
// @Router /api/verify/email/{email} [get]
// @Param email path string true "email"
// @Param scope query string false "scope"
// @Success 200 {object} EmailVerifyResponse
// @Failure 400 {object} utils.MessageResponse “email不在白名单中”
// @Failure 500 {object} utils.MessageResponse
func VerifyWithEmailOld(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// VerifyWithEmail godoc
//
// @Summary verify with email in query
// @Description verify with email in query, Send verification email
// @Tags Account
// @Produce json
// @Router /api/verify/email [get]
// @Param email query string true "email"
// @Param scope query string false "scope"
// @Success 200 {object} EmailVerifyResponse
// @Failure 400 {object} utils.MessageResponse
// @Failure 403 {object} utils.MessageResponse “email不在白名单中”
// @Failure 500 {object} utils.MessageResponse
func VerifyWithEmail(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// VerifyWithApikey godoc
//
// @Summary verify with email in query and apikey
// @Description verify with email in query, return verification code
// @Deprecated
// @Tags Account
// @Produce json
// @Router /api/verify/apikey [get]
// @Param email query ApikeyRequest true "apikey, email"
// @Success 200 {object} ApikeyResponse
// @Success 200 {object} utils.MessageResponse "用户未注册"
// @Failure 403 {object} utils.MessageResponse "apikey不正确"
// @Failure 409 {object} utils.MessageResponse "用户已注册"
// @Failure 500 {object} utils.MessageResponse
func VerifyWithApikey(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// DeleteUserByMe
//
// @Summary delete user
// @Description delete user account and related jwt credentials
// @Tags Account
// @Router /api/users/me [delete]
// @Param json body LoginRequest true "email, password"
// @Success 204
// @Failure 400 {object} utils.MessageResponse "密码错误“
// @Failure 404 {object} utils.MessageResponse "用户不存在“
// @Failure 500 {object} utils.MessageResponse
func DeleteUserByMe(c *fiber.Ctx) error {
	return c.JSON(nil)
}
