package repository

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/eko/gocache/lib/v4/store"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
	"gorm.io/gorm/clause"

	"github.com/opentreehole/backend/internal/model"
	"github.com/opentreehole/backend/internal/schema"
	"github.com/opentreehole/backend/pkg/utils"
)

type AccountRepository interface {
	Repository

	// FindUserByID find a user by id
	FindUserByID(ctx context.Context, id int) (user *model.User, err error)

	// GetUserByEmail get a user by email
	GetUserByEmail(ctx context.Context, email string) (user *model.User, err error)

	// CheckIfUserExists check whether user exists
	CheckIfUserExists(ctx context.Context, email string) (bool, error)

	// CheckIfUserDeleted check whether user is deleted
	CheckIfUserDeleted(ctx context.Context, email string) (bool, error)

	// CreateUser create a user
	CreateUser(ctx context.Context, email, password string) (user *model.User, err error)

	// AddDeletedIdentifier add deleted identifier to database
	AddDeletedIdentifier(ctx context.Context, userID int, identifier string) error

	// MakeIdentifier make user identifier from email
	MakeIdentifier(ctx context.Context, email string) string

	// MakePassword make password from rawPassword
	// using pbkdf2_sha256 algorithm
	MakePassword(ctx context.Context, rawPassword string) (string, error)

	// CheckPassword check whether rawPassword matches encryptedPassword
	CheckPassword(ctx context.Context, rawPassword, encryptedPassword string) error

	// CreateJWTToken create jwt token for user
	CreateJWTToken(ctx context.Context, user *model.User) (access, refresh string, err error)

	// CheckVerificationCode check whether verificationCode matches email
	CheckVerificationCode(ctx context.Context, scope, email, verificationCode string) error

	// SetVerificationCode set verificationCode to cache
	SetVerificationCode(ctx context.Context, email, scope string) (string, error)

	// DeleteVerificationCode delete verificationCode from cache
	DeleteVerificationCode(ctx context.Context, email, scope string) error

	/* 用于鉴权 */

	// GetCurrentUser get current user
	GetCurrentUser(ctx context.Context) (user *model.User, err error)
}

type accountRepository struct {
	Repository
}

func NewAccountRepository(repository Repository) AccountRepository {
	return &accountRepository{Repository: repository}
}

/* 接口实现 */

func (a *accountRepository) FindUserByID(ctx context.Context, id int) (user *model.User, err error) {
	var u model.User
	// TODO: 缓存
	return &u, a.GetDB(ctx).First(&u, id).Error
}

func (a *accountRepository) GetUserByEmail(ctx context.Context, email string) (user *model.User, err error) {
	var u model.User
	return &u, a.GetDB(ctx).Where("identifier = ?", a.MakeIdentifier(ctx, email)).First(&u).Error
}

func (a *accountRepository) CheckIfUserExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := a.GetDB(ctx).Raw("SELECT EXISTS (SELECT 1 FROM users WHERE identifier = ?)", a.MakeIdentifier(ctx, email)).Scan(&exists).Error
	return exists, err
}

func (a *accountRepository) CheckIfUserDeleted(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := a.GetDB(ctx).Raw("SELECT EXISTS (SELECT 1 FROM delete_identifier WHERE identifier = ?)", a.MakeIdentifier(ctx, email)).Scan(&exists).Error
	return exists, err
}

func (a *accountRepository) CreateUser(ctx context.Context, email, password string) (user *model.User, err error) {
	user = &model.User{
		Nickname: "user",
		Identifier: sql.NullString{
			String: a.MakeIdentifier(ctx, email),
			Valid:  true,
		},
		Password:      utils.Must(a.MakePassword(ctx, password)),
		UserJwtSecret: randstr.Base62(32),
		IsActive:      true,
	}

	return user, a.GetDB(ctx).Create(user).Error
}

func (a *accountRepository) AddDeletedIdentifier(ctx context.Context, userID int, identifier string) error {
	deleteIdentifier := model.DeleteIdentifier{UserID: userID, Identifier: identifier}
	return a.GetDB(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&deleteIdentifier).Error
}

func (a *accountRepository) MakeIdentifier(ctx context.Context, email string) string {
	decryptedIdentifierSalt := a.GetConfig(ctx).DecryptedIdentifierSalt
	return hex.EncodeToString(
		pbkdf2.Key([]byte(email), decryptedIdentifierSalt, 1, 64, sha3.New512),
	)
}

func (a *accountRepository) MakePassword(_ context.Context, rawPassword string) (string, error) {
	const (
		algorithm  = "sha256"
		iterations = 216000
	)
	salt, err := saltGenerator(12)
	if err != nil {
		return "", err
	}
	hashBase64 := passwordHash([]byte(rawPassword), salt, iterations, 32, sha256.New)
	return fmt.Sprintf("pbkdf2_%v$%v$%v$%v", algorithm, iterations, string(salt), hashBase64), nil
}

func (a *accountRepository) CheckPassword(_ context.Context, rawPassword, encryptedPassword string) error {
	splitEncryptedPassword := strings.Split(encryptedPassword, "$")
	if len(splitEncryptedPassword) != 4 {
		return fmt.Errorf("parse encryptedPassword error: %v", encryptedPassword)
	}
	algorithm := splitEncryptedPassword[0]
	splitAlgorithm := strings.Split(algorithm, "_")
	if len(splitAlgorithm) != 2 {
		return fmt.Errorf("parse encryptedPassword algorithm error: %v", encryptedPassword)
	}

	var hashOutputSize int
	var hashFactory func() hash.Hash
	if splitAlgorithm[1] == "sha256" {
		hashOutputSize = 32
		hashFactory = sha256.New
	} else {
		return fmt.Errorf("invalid sum algorithm: %v", splitAlgorithm[1])
	}

	iterations, err := strconv.Atoi(splitEncryptedPassword[1])
	if err != nil {
		return err
	}

	salt := splitEncryptedPassword[2]
	hashBase64 := passwordHash([]byte(rawPassword), []byte(salt), iterations, hashOutputSize, hashFactory)

	if hashBase64 != splitEncryptedPassword[3] {
		return fmt.Errorf("密码错误")
	}
	return nil
}

func (a *accountRepository) CreateJWTToken(ctx context.Context, user *model.User) (access, refresh string, err error) {
	var (
		key          = fmt.Sprintf("user_%d", user.ID)
		secret       = user.UserJwtSecret
		claim        = (&schema.UserClaims{}).FromModel(user)
		accessToken  string
		refreshToken string
	)

	if a.GetConfig(ctx).Features.ExternalGateway {
		// TODO: get key from kong or other api gateway
	}

	if !a.GetConfig(ctx).Features.RegistrationTest {
		claim.HasAnsweredQuestions = true
	}

	if user.UserJwtSecret == "" {
		// generate jwt secret
		user.UserJwtSecret = randstr.Base62(32)
		err = a.GetDB(ctx).Model(user).Update("user_jwt_secret", user.UserJwtSecret).Error
		if err != nil {
			return "", "", err
		}

		secret = user.UserJwtSecret
	}

	claim.Issuer = key

	claim.Type = schema.JWTTypeAccess
	claim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(30 * time.Minute)) // 30 minutes
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	claim.Type = schema.JWTTypeRefresh
	claim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)) // 30 days
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secret))

	return accessToken, refreshToken, nil
}

func (a *accountRepository) CheckVerificationCode(ctx context.Context, scope, email, verificationCode string) error {
	var storedCode string
	_, err := a.GetCache(ctx).Get(ctx, fmt.Sprintf("%v-%v", scope, a.MakeIdentifier(ctx, email)), &storedCode)
	if err != nil {
		return err
	}

	if storedCode != verificationCode {
		return fmt.Errorf("验证码错误") // TODO i18n
	}
	return nil
}

func (a *accountRepository) SetVerificationCode(ctx context.Context, email, scope string) (string, error) {
	codeInt, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	code := fmt.Sprintf("%06d", codeInt.Uint64())

	return code, a.GetCache(ctx).Set(
		ctx,
		fmt.Sprintf("%v-%v", scope, a.MakeIdentifier(ctx, email)),
		code,
		store.WithExpiration(time.Second*time.Duration(a.GetConfig(ctx).Features.VerificationCodeExpires)),
	)
}

func (a *accountRepository) DeleteVerificationCode(ctx context.Context, email, scope string) error {
	return a.GetCache(ctx).Delete(
		ctx,
		fmt.Sprintf("%v-%v", scope, a.MakeIdentifier(ctx, email)),
	)
}

func (a *accountRepository) GetCurrentUser(ctx context.Context) (user *model.User, err error) {
	// get fiber context
	c := a.GetFiberCtx(ctx)

	// 只有在开启 auth 模块的情况下才能从数据库里读 user
	if a.GetConfig(ctx).Modules.Auth {
		// get user id from header: X-Consumer-Username if through Kong
		username := c.Get("X-Consumer-Username")
		if username != "" {
			id, err := strconv.Atoi(username)
			if err == nil {
				return a.FindUserByID(ctx, id)
			}
		}
	}

	// get user id from jwt
	// ID and UserID are both valid
	var userClaims schema.UserClaims

	token := GetJWTToken(c)
	if token == "" {
		return nil, schema.Unauthorized("Unauthorized")
	}

	err = ParseJWTToken(token, &userClaims)
	if err != nil {
		return nil, schema.Unauthorized("Unauthorized")
	}

	if userClaims.ID == 0 && userClaims.UserID == 0 && userClaims.UID == 0 {
		return nil, schema.Unauthorized("Unauthorized")
	}

	if userClaims.ExpiresAt != nil && userClaims.ExpiresAt.Time.Before(time.Now()) {
		return nil, schema.Unauthorized("token expired")
	}

	user = userClaims.ToModel()

	if a.GetConfig(ctx).Modules.Auth {
		return a.FindUserByID(ctx, user.ID)
	}

	return user, nil
}

/* 工具函数，非导出函数 */

func passwordHash(bytePassword, salt []byte, iterations, KeyLen int, hash func() hash.Hash) string {
	return base64.StdEncoding.EncodeToString(pbkdf2.Key(bytePassword, salt, iterations, KeyLen, hash))
}

func saltGenerator(stringLen int) ([]byte, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsLength := len(chars)
	var builder bytes.Buffer
	for i := 0; i < stringLen; i++ {
		choiceIndex, err := rand.Int(rand.Reader, big.NewInt(int64(charsLength)))
		if err != nil {
			return nil, err
		}
		err = builder.WriteByte(chars[choiceIndex.Int64()])
		if err != nil {
			return nil, err
		}
	}
	return builder.Bytes(), nil
}

// GetJWTToken extracts token from header or cookie
// return empty string if not found
func GetJWTToken(c *fiber.Ctx) string {
	tokenString := c.Get("Authorization") // token in header
	if tokenString == "" {
		tokenString = c.Cookies("access") // token in cookie
	}
	return tokenString
}

var (
	ErrJWTTokenRequired = schema.Unauthorized("jwt token required")
	ErrInvalidJWTToken  = schema.Unauthorized("invalid jwt token")
)

// ParseJWTToken extracts and parse token, whatever start with "Bearer " or not
func ParseJWTToken(token string, user any) error {
	// remove "Bearer " prefix if exists
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}
	token = strings.TrimSpace(token)
	payloads := strings.SplitN(token, ".", 3) // extract "Bearer "
	if len(payloads) < 3 {
		return ErrJWTTokenRequired
	}

	payloadString := payloads[1]

	// jwt encoding ignores padding, so RawStdEncoding should be used instead of StdEncoding
	// jwt encoding uses url safe base64 encoding, so RawURLEncoding should be used instead of RawStdEncoding
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadString) // the middle one is payload
	if err != nil {
		log.Err(err).Str("payload_string", payloadString).Msg("jwt parse error")
		return ErrInvalidJWTToken
	}

	err = json.Unmarshal(payloadBytes, user)
	if err != nil {
		log.Err(err).Str("payload_string", payloadString).Msg("jwt parse error")
		return ErrInvalidJWTToken
	}

	return nil
}
