package service

import (
	"context"
	"time"

	"evm-tx-watcher/internal/domain"
	"evm-tx-watcher/internal/dto"
	"evm-tx-watcher/internal/errors"
	"evm-tx-watcher/internal/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AddressService interface {
	Register(ctx context.Context, address *dto.RegisterAddressRequest) (*dto.AddressResponse, *errors.AppError)
	GetAll(ctx context.Context) ([]*dto.AddressResponse, *errors.AppError)
}

type addressService struct {
	unitOfWork  repository.UnitOfWork
	addressRepo repository.AddressRepository
	webhookRepo repository.WebhookRepository
}

func NewAddressService(unitOfWork repository.UnitOfWork, repo repository.AddressRepository, webhookRepo repository.WebhookRepository) AddressService {
	return &addressService{unitOfWork: unitOfWork, addressRepo: repo, webhookRepo: webhookRepo}
}
func (s *addressService) Register(ctx context.Context, address *dto.RegisterAddressRequest) (*dto.AddressResponse, *errors.AppError) {
	existingAddress, err := s.addressRepo.FindByAddress(ctx, address.Address)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeDatabase, "failed to check existing address", err)
	}

	if existingAddress != nil {
		return nil, errors.AlreadyExists("Address already exists")
	}

	newAddress := &domain.Address{
		ID:          uuid.New(),
		Address:     address.Address,
		ChainID:     address.ChainID,
		IsContract:  false, // Default value, can be updated later
		IsActive:    true,  // Default to active
		Label:       address.Label,
		Description: address.Description,
		UserID:      nil, // Assuming no user association for now
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	newWebhook := &domain.Webhook{
		ID:        uuid.New(),
		AddressID: newAddress.ID,
		URL:       address.WebhookURL,
		Secret:    address.Secret,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var createdAddress domain.Address

	err = s.unitOfWork.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		addr, err := s.addressRepo.Create(ctx, tx, newAddress)
		if err != nil {
			return err
		}
		createdAddress = addr
		_, err = s.webhookRepo.Create(ctx, tx, newWebhook)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeDatabase, "failed to register address", err)
	}

	// Map UserID to string pointer
	var userID *string
	if createdAddress.UserID != nil {
		s := createdAddress.UserID.String()
		userID = &s
	}

	return &dto.AddressResponse{
		ID:          createdAddress.ID.String(),
		Address:     createdAddress.Address,
		ChainID:     createdAddress.ChainID,
		IsContract:  createdAddress.IsContract,
		IsActive:    createdAddress.IsActive,
		Label:       createdAddress.Label,
		WebhookURL:  newWebhook.URL,
		Description: createdAddress.Description,
		UserID:      userID,
		CreatedAt:   createdAddress.CreatedAt,
		UpdatedAt:   createdAddress.UpdatedAt,
	}, nil
}

func (s *addressService) GetAll(ctx context.Context) ([]*dto.AddressResponse, *errors.AppError) {
	addresses, err := s.addressRepo.FindAll(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeDatabase, "failed to get all addresses", err)
	}

	var responses []*dto.AddressResponse
	for _, addr := range addresses {
		var userID *string
		if addr.UserID != nil {
			s := addr.UserID.String()
			userID = &s
		}
		responses = append(responses, &dto.AddressResponse{
			ID:          addr.ID.String(),
			Address:     addr.Address,
			ChainID:     addr.ChainID,
			IsContract:  addr.IsContract,
			IsActive:    addr.IsActive,
			Label:       addr.Label,
			Description: addr.Description,
			UserID:      userID,
			CreatedAt:   addr.CreatedAt,
			UpdatedAt:   addr.UpdatedAt,
		})
	}

	return responses, nil
}
