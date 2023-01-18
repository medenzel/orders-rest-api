package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/medenzel/orders-rest-api/internal/order"
	log "github.com/sirupsen/logrus"
)

func sendResponse(w http.ResponseWriter, data any) error {
	switch d := data.(type) {
	case string:
		if err := json.NewEncoder(w).Encode(Response{Message: d}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
	default:
		if err := json.NewEncoder(w).Encode(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
	}
	return nil
}

// OrderService - define the interface that the concrete implementation
// has to adhere to
type OrderService interface {
	GetOrder(context.Context, int) (order.Order, error)
	GetAllOrders(context.Context, int, int) ([]order.Order, error)
	PostOrder(context.Context, order.Order) (order.Order, error)
	UpdateOrder(context.Context, int, order.Order) (order.Order, error)
	DeleteOrder(context.Context, int) error
}

// GetOrder - retrieves an order by id and returns response
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	stringOrderID := chi.URLParam(r, "id")
	if stringOrderID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	intOrderID, err := strconv.Atoi(stringOrderID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ord, err := h.Service.GetOrder(r.Context(), intOrderID)
	if err != nil {
		if errors.Is(err, order.ErrNoOrderFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(ord); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type GetAllRequestParams struct {
	PageID   int32 `validate:"required,min=1"`
	PageSize int32 `validate:"required,min=5,max=10"`
}

// GetAllOrder - retrieves all orders and returns response with pagination
func (h *Handler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	pageID, err := strconv.Atoi(r.URL.Query().Get("page_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := "page_id param doesn't set"
		if err := sendResponse(w, msg); err != nil {
			log.Errorf("sending response: %w", err)
		}
		return
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := "page_size param doesn't set"
		if err := sendResponse(w, msg); err != nil {
			log.Errorf("sending response: %w", err)
		}
		return
	}

	if pageID < 1 || pageSize < 5 || pageSize > 10 {
		w.WriteHeader(http.StatusBadRequest)
		msg := "Bad params: page_id min = 1, page_size min = 5, max = 10"
		if err := sendResponse(w, msg); err != nil {
			log.Errorf("sending response: %w", err)
		}
		return
	}

	ords, err := h.Service.GetAllOrders(r.Context(), pageSize, (pageID-1)*pageSize)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(ords); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// PostOrderRequest - clone of order struct, helps to validate fields
type PostOrderRequest struct {
	Description string `json:"description" validatate:"required"`
	State       int    `json:"state" validate:"required,oneof=1 2 3 4"`
	CreateAt    string `json:"create_at" validate:"omitempty,datetime=02/01/2006 15:04:05"`
}

// orderFromPostOrderRequest - converts the validated struct into order
func orderFromPostOrderRequest(por PostOrderRequest) order.Order {
	return order.Order{
		Description: por.Description,
		State:       por.State,
		CreateAt:    por.CreateAt,
	}
}

// PostOrder - adds a new order
func (h *Handler) PostOrder(w http.ResponseWriter, r *http.Request) {
	var postOrderReq PostOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&postOrderReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(postOrderReq)
	if err != nil {
		log.Info(fmt.Errorf("post order validate: %w", err))
		validationErrors := err.(validator.ValidationErrors)
		errMsg := "Incorrect fields: "
		for _, err := range validationErrors {
			errMsg += err.StructField() + "|"
		}
		w.WriteHeader(http.StatusBadRequest)
		if err := sendResponse(w, errMsg); err != nil {
			log.Errorf("sending response: %w", err)
		}
		return
	}

	ord := orderFromPostOrderRequest(postOrderReq)
	postedOrd, err := h.Service.PostOrder(r.Context(), ord)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(postedOrd); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateOrderRequest - clone of order struct, helps to validate fields
type UpdateOrderRequest struct {
	Description string `json:"description" validatate:"required"`
	State       int    `json:"state" validate:"required,oneof=1 2 3 4"`
}

// orderFromUpdateOrderRequest - converts the validated struct into order
func orderFromUpdateOrderRequest(uor UpdateOrderRequest) order.Order {
	return order.Order{
		Description: uor.Description,
		State:       uor.State,
	}
}

// UpdateOrder - updates an order by ID
func (h *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	stringOrderID := chi.URLParam(r, "id")
	if stringOrderID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	intOrderID, err := strconv.Atoi(stringOrderID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var updateOrderReq UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&updateOrderReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(updateOrderReq)
	if err != nil {
		log.Info(fmt.Errorf("update order validate: %w", err))
		validationErrors := err.(validator.ValidationErrors)
		errMsg := "Incorrect fields: "
		for _, err := range validationErrors {
			errMsg += err.StructField() + "|"
		}
		w.WriteHeader(http.StatusBadRequest)
		if err := sendResponse(w, errMsg); err != nil {
			log.Errorf("sending response: %w", err)
		}
		return
	}

	ord := orderFromUpdateOrderRequest(updateOrderReq)
	updatedOrd, err := h.Service.UpdateOrder(r.Context(), intOrderID, ord)
	if err != nil {
		if errors.Is(err, order.ErrNoOrderFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := sendResponse(w, updatedOrd); err != nil {
		log.Errorf("sending response: %w", err)
	}
}

// DeleteOrder - deletes order by ID
func (h *Handler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	stringOrderID := chi.URLParam(r, "id")
	if stringOrderID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	intOrderID, err := strconv.Atoi(stringOrderID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.Service.DeleteOrder(r.Context(), intOrderID)
	if err != nil {
		if errors.Is(err, order.ErrNoOrderFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := sendResponse(w, "Successfully deleted"); err != nil {
		log.Errorf("sending response: %w", err)
	}
}
