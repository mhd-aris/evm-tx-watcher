package service

import (
	"context"

	"evm-tx-watcher/internal/domain"
	"evm-tx-watcher/internal/dto"
	"evm-tx-watcher/internal/repository"
	"evm-tx-watcher/internal/validator"

	v10 "github.com/go-playground/validator/v10"
)

type AddressService interface {
	Register(ctx context.Context, address *dto.RegisterRequest) (domain.Address, error)
	GetAll(ctx context.Context) ([]*domain.Address, error)
}

type addressService struct {
	repo      repository.AddressRepository
	validator *v10.Validate
}

func NewAddressService(repo repository.AddressRepository) AddressService {
	return &addressService{repo: repo,
		validator: validator.New(),
	}
}

func (s *addressService) Register(ctx context.Context, address *dto.RegisterRequest) (domain.Address, error) {

	if err := s.validator.Struct(address); err != nil {
		return domain.Address{}, err
	}

	domainAddr := &domain.Address{
		Address: address.Address,
		ChainID: address.ChainID,
	}

	createdAddress, err := s.repo.Create(ctx, domainAddr)
	if err != nil {
		return domain.Address{}, err
	}
	domainAddr.ID = createdAddress.ID
	return *domainAddr, nil
}

func (s *addressService) GetAll(ctx context.Context) ([]*domain.Address, error) {
	return s.repo.GetAll(ctx)
}
