package wallet

import (
	"DX/src/domain/entity"
	"time"
)

type Wallet struct {
	Balance      int64         `json:"balance"`
	Escrow       int64         `json:"escrow"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Id              entity.DatabaseId `json:"id" bson:"_id"`
	UserId          string            `json:"user_id" bson:"user_id"`
	TransactionType Type              `json:"transaction_type" bson:"transaction_type"`
	Type            string            `json:"type" bson:"type"`
	ItemId          string            `json:"-" json:"item_id"`
	Amount          int64             `json:"amount" bson:"amount"`
	Description     string            `json:"description" bson:"description"`
	CreatedAt       time.Time         `json:"created_at" bson:"created_at"`
}

type Type int

const (
	Debit Type = iota
	Credit
)

func (t Type) String() string {
	if t == Debit {
		return "debit"
	}
	if t == Credit {
		return "credit"
	}
	return ""
}

func NewDebitTransaction(userId, description, itemId string, amount int64) Transaction {
	return Transaction{
		Id:              entity.NewDatabaseId(),
		UserId:          userId,
		Type:            Debit.String(),
		TransactionType: Debit,
		CreatedAt:       time.Now(),
		Description:     description,
		Amount:          amount,
		ItemId:          itemId,
	}
}

func NewCreditTransaction(userId, description, itemId string, amount int64) Transaction {
	return Transaction{
		Id:              entity.NewDatabaseId(),
		UserId:          userId,
		Type:            Credit.String(),
		TransactionType: Credit,
		CreatedAt:       time.Now(),
		Description:     description,
		Amount:          amount,
		ItemId:          itemId,
	}
}
