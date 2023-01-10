package order

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
)

var (
	ErrGettingOrder     = errors.New("could not get order by ID")
	ErrGettingAllOrders = errors.New("could not get all orders")
	ErrPostingOrder     = errors.New("could not post order")
	ErrUpdatingOrder    = errors.New("could not update order")
	ErrDeletingOrder    = errors.New("could not delete order")
	ErrNoOrderFound     = errors.New("no orders found")
)

// Order - defines order structure
type Order struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	State       int    `json:"state"`
	CreateAt    string `json:"create_at"`
}

// CommentStore - defines the interface, that order storage layer
// needs to implement
type OrderStore interface {
	GetOrder(context.Context, int) (Order, error)
	GetAllOrders(context.Context) ([]Order, error)
	PostOrder(context.Context, Order) (Order, error)
	UpdateOrder(context.Context, int, Order) (Order, error)
	DeleteOrder(context.Context, int) error
}

// Service - struct for the order service
type Service struct {
	Store OrderStore
}

// NewService - returnes a new order service
func NewService(store OrderStore) *Service {
	return &Service{
		Store: store,
	}
}

// GetOrder - retrieves order by it's ID from the database
func (s *Service) GetOrder(ctx context.Context, ID int) (Order, error) {
	ord, err := s.Store.GetOrder(ctx, ID)
	if err != nil {
		if errors.Is(err, ErrNoOrderFound) {
			return Order{}, ErrNoOrderFound
		}
		log.Errorf("get order: %w", err)
		return Order{}, ErrGettingOrder
	}
	return ord, nil
}

// GetAllOrders - retrieves all orders from the database
func (s *Service) GetAllOrders(ctx context.Context) ([]Order, error) {
	ords, err := s.Store.GetAllOrders(ctx)
	if err != nil {
		log.Errorf("get all orders: %w", err)
		return nil, ErrGettingAllOrders
	}
	return ords, nil
}

// PostOrder - adds a new order to the database
func (s *Service) PostOrder(ctx context.Context, ord Order) (Order, error) {
	ord, err := s.Store.PostOrder(ctx, ord)
	if err != nil {
		log.Errorf("post order: %w", err)
		return Order{}, ErrPostingOrder
	}
	return ord, nil
}

// UpdateOrder - update an order by ID with new order info
func (s *Service) UpdateOrder(ctx context.Context, ID int, newOrd Order) (Order, error) {
	ord, err := s.Store.UpdateOrder(ctx, ID, newOrd)
	if err != nil {
		if errors.Is(err, ErrNoOrderFound) {
			return Order{}, ErrNoOrderFound
		}
		log.Errorf("update order: %w", err)
		return Order{}, ErrUpdatingOrder
	}
	return ord, nil
}

// DeleteOrder - deletes an order from the database by ID
func (s *Service) DeleteOrder(ctx context.Context, ID int) error {
	err := s.Store.DeleteOrder(ctx, ID)
	if err != nil {
		if errors.Is(err, ErrNoOrderFound) {
			return err
		}
		log.Errorf("delete order: %w", err)
		return ErrDeletingOrder
	}
	return nil
}
