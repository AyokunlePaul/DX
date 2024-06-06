package wallet

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type reader interface {
	CreateTransaction(Transaction) error
}

type writer interface {
	GetTransactionsFor(string) ([]Transaction, error)
	GetBalance(string) (int64, error)
}

type Repository interface {
	reader
	writer
}

type repository struct {
	Collection *mongo.Collection
}

func NewWalletRepository(Collection *mongo.Collection) Repository {
	return &repository{
		Collection: Collection,
	}
}

func (r *repository) CreateTransaction(txn Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := r.Collection.InsertOne(ctx, txn); err != nil {
		return err
	}

	return nil
}

func (r *repository) GetTransactionsFor(userId string) (transactions []Transaction, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": userId,
	}

	crs, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = crs.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *repository) GetBalance(userId string) (int64, error) {
	var balances []struct {
		Balance int64 `bson:"balance"`
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"user_id", userId},
		},
		}},
		{{
			"$group", bson.D{
				{"_id", nil},
				{"balance", bson.D{
					{"$sum", bson.D{
						{"$cond", bson.A{
							bson.D{{"$eq", bson.A{"$type", "credit"}}},
							"$amount",
							bson.D{{"$multiply", bson.A{"$amount", -1}}},
						}},
					}},
				}},
			},
		}},
		{{
			"$project", bson.D{
				{"_id", 0},
				{"balance", 1},
			},
		}},
	}

	crs, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	err = crs.All(ctx, &balances)
	if err != nil {
		return 0, err
	}

	if len(balances) == 0 {
		return 0, nil
	}
	return balances[0].Balance, nil
}
