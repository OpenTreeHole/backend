package accountV3

import (
	"github.com/gofiber/fiber/v2"

	accountV2 "github.com/opentreehole/backend/internal/auth/handler/v2/account"
)

// Register godoc
//
// @Summary Register
// @Description Register with email, password and optional verification code if enabled
// @Tags Account V3
// @Accept json
// @Produce json
// @Router /api/v3/register [post]
// @Param json body accountV2.RegisterRequest true "RegisterRequest"
// @Success 201 {object} accountV2.TokenResponse
// @Failure 400 {object} utils.MessageResponse "验证码错误，用户已存在"
func Register(c *fiber.Ctx) error {
	return accountV2.Register(c)
}

// RegisterDebug godoc
//
// @Summary register in batch, debug only
// @Description register with email, password, not need verification code
// @Tags Account V3
// @Accept json
// @Produce json
// @Router /api/v3/debug/register [post]
// @Param json body accountV2.RegisterInBatchRequest true "RegisterInBatchRequest"
// @Success 201 {object} accountV2.TokenResponse
// @Failure 400 {object} utils.MessageResponse "用户已注册"
// @Failure 500 {object} utils.MessageResponse
// @Security ApiKeyAuth
func RegisterDebug(c *fiber.Ctx) error {
	return accountV2.RegisterDebugInBatch(c)
}

// Login godoc
//
// @Summary Login
// @Description Login with email and password, return jwt token, no need jwt
// @Tags Account V3
// @Accept json
// @Produce json
// @Router /api/v3/login [post]
// @Param json body accountV2.LoginRequest true "LoginRequest"
// @Success 200 {object} accountV2.TokenResponse
// @Failure 400 {object} utils.MessageResponse
// @Failure 404 {object} utils.MessageResponse
// @Failure 500 {object} utils.MessageResponse
func Login(c *fiber.Ctx) error {
	return accountV2.Login(c)
}

// Logout godoc
//
// @Summary Logout
// @Description Logout, need jwt
// @Tags Account V3
// @Produce json
// @Router /api/v3/logout [post]
// @Success 200 {object} utils.MessageResponse
// @Failure 400 {object} utils.MessageResponse
// @Failure 404 {object} utils.MessageResponse
// @Failure 500 {object} utils.MessageResponse
func Logout(c *fiber.Ctx) error {
	return accountV2.Logout(c)
}

// Refresh godoc
//
// @Summary Refresh
// @Description Refresh jwt token, need jwt
// @Tags Account V3
// @Accept json
// @Produce json
// @Router /api/v3/refresh [post]
// @Success 200 {object} accountV2.TokenResponse
func Refresh(c *fiber.Ctx) error {
	return accountV2.Refresh(c)
}

// ResetPassword godoc
//
// @Summary reset password
// @Description reset password and jwt credential
// @Tags Account V3
// @Accept json
// @Produce json
// @Router /api/v3/register [put]
// @Param json body accountV2.RegisterRequest true "RegisterRequest"
// @Success 200 {object} accountV2.TokenResponse
// @Failure 400 {object} utils.MessageResponse "验证码错误"
// @Failure 500 {object} utils.MessageResponse
func ResetPassword(c *fiber.Ctx) error {
	return accountV2.ResetPassword(c)
}

// VerifyWithEmail godoc
//
// @Summary verify with email in query
// @Description verify with email in query, Send verification email
// @Tags Account V3
// @Produce json
// @Router /api/v3/verify/email [get]
// @Param email query string true "email"
// @Param scope query string false "scope"
// @Success 200 {object} accountV2.EmailVerifyResponse
// @Failure 400 {object} utils.MessageResponse
// @Failure 403 {object} utils.MessageResponse "email不在白名单中"
// @Failure 500 {object} utils.MessageResponse
func VerifyWithEmail(c *fiber.Ctx) error {
	return accountV2.VerifyWithEmail(c)
}

// DeleteUserByMe
//
// @Summary delete user by me
// @Description delete user account and related jwt credentials
// @Tags Account V3
// @Router /api/v3/users/me [delete]
// @Param json body accountV2.LoginRequest true "email, password"
// @Success 204
// @Failure 400 {object} utils.MessageResponse "密码错误"
// @Failure 404 {object} utils.MessageResponse "用户不存在"
// @Failure 500 {object} utils.MessageResponse
func DeleteUserByMe(c *fiber.Ctx) error {
	return accountV2.DeleteUserByMe(c)
}

// DeleteUserByID
//
// @Summary delete user by id
// @Description delete user account, admin only
// @Tags Account V3
// @Router /api/v3/users/{id} [delete]
// @Param id path string true "user id"
// @Success 204
// @Failure 404 {object} utils.MessageResponse "用户不存在"
// @Failure 500 {object} utils.MessageResponse
func DeleteUserByID(c *fiber.Ctx) error {
	return c.JSON(nil)
}
