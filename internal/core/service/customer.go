package service

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
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
