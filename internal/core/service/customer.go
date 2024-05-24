package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	errors "github.com/rotisserie/eris"
	"os"
	"synapsis-challenge/internal/core/domain"
	inservice "synapsis-challenge/internal/core/port/inbound/service"
	"synapsis-challenge/internal/core/port/outbound/registry"
	"synapsis-challenge/shared"
	"time"
)

type CustomerService struct {
	repositoryRegistry registry.RepositoryRegistry
}

func NewAccountService(repositoryRegistry registry.RepositoryRegistry) inservice.CustomerService {
	return &CustomerService{
		repositoryRegistry: repositoryRegistry,
	}
}

func (i *CustomerService) RegisterCustomer(ctx context.Context, customer *domain.Customer) (string, error) {
	var (
		id           string
		err          error
		dataCustomer *domain.Customer

		repo = i.repositoryRegistry.GetCustomerRepository()
	)

	dataCustomer, err = repo.FindByEmail(ctx, customer.Email)
	if err != nil && !errors.Is(err, shared.ErrNotFound) {
		return id, errors.Wrap(err, "RegisterCustomer.FindByEmail")
	}

	if dataCustomer != nil {
		return id, shared.ErrAlreadyExist
	}

	customer.Password, err = shared.EncryptPassword(customer.Password)
	if err != nil {
		return id, errors.Wrap(err, "RegisterCustomer.EncryptPassword")
	}

	customer.ID = uuid.NewString()
	id = customer.ID

	err = repo.Create(ctx, customer)
	if err != nil {
		return id, errors.Wrap(err, "RegisterCustomer.Create")
	}

	return id, nil
}

func (i *CustomerService) SignIn(ctx context.Context, customer *domain.Customer) (*domain.SignIn, error) {
	var (
		err          error
		dataCustomer *domain.Customer
		result       *domain.SignIn

		repo      = i.repositoryRegistry.GetCustomerRepository()
		redisRepo = i.repositoryRegistry.GetRedisRepository()
	)

	dataCustomer, err = repo.FindByEmail(ctx, customer.Email)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return result, shared.ErrNotFound
		}

		return result, errors.Wrap(err, "SignIn.FindByEmail")
	}

	isValid, err := shared.ComparePassword(dataCustomer.Password, customer.Password)
	if err != nil {
		return result, err
	}

	if !isValid {
		return result, shared.ErrInvalidPassword
	}

	accessTokenKey, refreshTokenKey := getAuthRedisKey(dataCustomer.ID)

	isExist, _ := redisRepo.IsExist(accessTokenKey)
	if isExist {
		return result, shared.ErrAlreadyExist
	}

	accessToken, refreshToken, err := generateSignInToken(dataCustomer.ID)
	if err != nil {
		return result, errors.Wrap(err, "SignIn.GenerateSignInToken")
	}

	err = redisRepo.Set(accessTokenKey, accessToken, shared.AccessTokenDuration)
	if err != nil {
		return result, errors.Wrap(err, "SignIn.Set.accessToken")
	}

	err = redisRepo.Set(refreshTokenKey, refreshToken, shared.RefreshTokenDuration)
	if err != nil {
		return result, errors.Wrap(err, "SignIn.Set.refreshToken")
	}

	result = &domain.SignIn{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		Customer:     dataCustomer,
	}

	return result, nil
}

func (i *CustomerService) VerifyToken(token string) error {
	var (
		err       error
		redisRepo = i.repositoryRegistry.GetRedisRepository()
	)

	validAccess, claimsAccess, err := verifyJWT(token)
	if err != nil {
		return err
	}

	if !validAccess {
		return errors.New("JWT access token is invalid")
	}

	accessTokenPayload, err := decodeToken(claimsAccess)
	if err != nil {
		return errors.Wrap(err, "failed decode access token claims")
	}

	if accessTokenPayload.Subject != shared.AccessTokenSubject {
		return errors.New("JWT token is not access token")
	}

	expirationTime := time.Unix(accessTokenPayload.ExpiresAt, 0)
	if time.Now().After(expirationTime) {
		return errors.New("access token is expired")
	}

	accessTokenKey, refreshTokenKey := getAuthRedisKey(accessTokenPayload.ID)

	exist, err := redisRepo.IsExist(accessTokenKey)
	if err != nil {
		return errors.Wrap(err, "error check redis key")
	}

	if !exist {
		return errors.New("user token not found")
	}

	refreshToken, err := redisRepo.GetString(refreshTokenKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.New("refresh token not found")
		}

		return errors.New("error get refresh token from redis")
	}

	validRefresh, _, err := verifyJWT(refreshToken)
	if err != nil {
		return errors.New("failed verify JWT refresh token")
	}

	if !validRefresh {
		return errors.New("JWT token is invalid")
	}

	err = redisRepo.Set(refreshTokenKey, refreshToken, shared.RefreshTokenDuration)
	if err != nil {
		return errors.New("error set jwt refresh_token to redis")
	}

	return nil
}

func (i *CustomerService) SignOut(ctx context.Context, customer *domain.Customer) error {
	return nil
}

func getAuthRedisKey(customerId string) (customerAccessToken, customerRefreshToken string) {
	customerAccessToken = fmt.Sprintf("customerAccessToken:%s", customerId)
	customerRefreshToken = fmt.Sprintf("customerRefreshToken:%s", customerId)
	return
}

func generateSignInToken(id string) (string, string, error) {
	accessToken, err := generateJWT(id, shared.AccessTokenSubject)
	if err != nil {
		return "", "", errors.Wrap(err, "GenerateJW.accessTokenT")
	}

	refreshToken, err := generateJWT(id, shared.RefreshTokenSubject)
	if err != nil {
		return "", "", errors.Wrap(err, "GenerateJWT.refreshToken")
	}

	return accessToken, refreshToken, nil
}

func generateJWT(id, subject string) (string, error) {
	godotenv.Load()
	claims := jwt.MapClaims{
		"customer_id": id,
		"sub":         subject,
	}
	secretKey := os.Getenv("JWT_SECRET_KEY")

	if subject == shared.AccessTokenSubject {
		claims["exp"] = time.Now().Add(shared.AccessTokenDuration).Unix()
	}

	if subject == shared.RefreshTokenSubject {
		claims["exp"] = time.Now().Add(shared.RefreshTokenDuration).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func verifyJWT(tokenStr string) (bool, jwt.MapClaims, error) {
	godotenv.Load()

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}

		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return false, nil, fmt.Errorf("failed to parse JWT token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, claims, nil
	}

	return false, nil, errors.New("invalid JWT token")
}

func decodeToken(claims map[string]any) (tokenPayload domain.JWTCustomer, err error) {
	byteToken, err := json.Marshal(claims)
	if err != nil {
		err = fmt.Errorf("failed marshal claims: %w", err)
		return
	}

	err = json.Unmarshal(byteToken, &tokenPayload)
	if err != nil {
		err = fmt.Errorf("failed unmarshal claims: %w", err)
		return
	}

	return
}
