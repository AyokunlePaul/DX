package wallet

import (
	"DX/src/domain/entity/auth"
	"DX/src/domain/entity/wallet"
	"DX/src/pkg/error_service"
	"errors"
)

type UseCase interface {
	AddMoney(string, int64) error
	RemoveMoney(string, int64) error
	GetWalletFor(string) (*wallet.Wallet, error)
	GetBalance() (int64, error)
}

type impl struct {
	auth.Manager
	wallet.Repository
	error_service.Service
}

func NewUseCase(repo wallet.Repository, service error_service.Service, authManager auth.Manager) UseCase {
	return &impl{
		Repository: repo,
		Service:    service,
		Manager:    authManager,
	}
}

func (i *impl) AddMoney(token string, amount int64) error {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	txn := wallet.NewCreditTransaction(*userId, "Credit", "", amount)
	err := i.Repository.CreateTransaction(txn)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("transaction", err).Message)
	}

	return nil
}

func (i *impl) RemoveMoney(token string, amount int64) error {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	txn := wallet.NewDebitTransaction(*userId, "Debit", "", amount)
	err := i.Repository.CreateTransaction(txn)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("transaction", err).Message)
	}

	return nil
}

func (i *impl) GetBalance() (int64, error) {
	if balance, err := i.Repository.GetBalance("64117a0cd472c91f3fe4834c"); err != nil {
		return 0, err
	} else {
		return balance, nil
	}
}

func (i *impl) GetWalletFor(token string) (*wallet.Wallet, error) {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return nil, errors.New(resp.Message)
	}

	txns, err := i.GetTransactionsFor(*userId)
	if err != nil {
		return nil, errors.New(i.HandleMongoDbError("transaction", err).Message)
	}

	balance, err := i.Repository.GetBalance(*userId)
	if err != nil {
		return nil, errors.New(i.HandleMongoDbError("balance", err).Message)
	}

	return &wallet.Wallet{
		Transactions: txns,
		Balance:      balance,
		Escrow:       0,
	}, nil
}
