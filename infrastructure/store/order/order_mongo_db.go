package order

import (
	"context"
	"time"

	"gopkg.in/mgo.v2/bson"

	"ordercontext/domain/customer"
	"ordercontext/domain/order"
	"ordercontext/domain/product"
	"ordercontext/infrastructure"
	"ordercontext/shared/aggregate"
)

const collectionName = "orders"

type orderBson struct {
	ID          string    `bson:"id"`
	CustomerID  string    `bson:"customerId"`
	ProductID   string    `bson:"productId"`
	CreatedTime time.Time `bson:"createdTime"`
	Status      int       `bson:"status"`
	Version     string    `bson:"version"`
}

func FromOrder(o *order.Order) *orderBson {
	return &orderBson{
		ID:         o.ID(),
		Status:     order.FromStatus(o.Status()),
		CustomerID: o.CustomerID(),
		ProductID:  o.ProductID(),
		Version:    o.Version(),
	}
}

func FromBson(o *orderBson) *order.Order {
	ord, _ := order.NewOrder(order.OrderID(o.ID),
		customer.ID(o.CustomerID),
		product.ID(o.ProductID),
		func() time.Time { return time.Now() },
		order.ToStatus(o.Status),
		aggregate.ToVersion(o.Version))

	ord.Clear()
	return ord
}

type OrderMongoRepository struct {
	mStore *infrastructure.MongoStore
}

func NewMongoRepository(mongoStore *infrastructure.MongoStore) *OrderMongoRepository {
	return &OrderMongoRepository{mStore: mongoStore}
}

func (r *OrderMongoRepository) GetOrders(ctx context.Context) ([]*order.Order, error) {

	var result []*orderBson

	if err := r.mStore.FindAll(ctx, collectionName, bson.M{}, &result); err != nil {
		return nil, err
	}

	var orders []*order.Order

	for _, o := range result {
		orders = append(orders, FromBson(o))
	}

	return orders, nil
}

func (r *OrderMongoRepository) Get(ctx context.Context, id string) (*order.Order, error) {
	var (
		bsonResult = &orderBson{}
		query      = bson.M{"id": id}
	)

	if err := r.mStore.FindOne(ctx, collectionName, query, nil, bsonResult); err != nil {
		return nil, err
	}

	return FromBson(bsonResult), nil
}

func (r *OrderMongoRepository) Update(ctx context.Context, o *order.Order) error {
	var (
		query  = bson.M{"id": o.ID(), "version": o.Version()}
		update = bson.M{"$set": bson.M{"status": o.Status(), "version": aggregate.NewVersion().String()}}
	)

	if err := r.mStore.Update(ctx, collectionName, query, update); err != nil {
		return err
	}
	return nil
}

func (r *OrderMongoRepository) Create(ctx context.Context, o *order.Order) error {
	bson := FromOrder(o)
	if bson.Version == "" {
		bson.Version = aggregate.NewVersion().String()
	}
	if err := r.mStore.Store(ctx, collectionName, bson); err != nil {
		return err
	}
	return nil
}