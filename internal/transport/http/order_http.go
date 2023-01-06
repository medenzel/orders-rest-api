package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/medenzel/orders-rest-api/internal/order"
)

type OrderService interface {
	GetOrder(context.Context, int) (order.Order, error)
	GetAllOrders(context.Context) ([]order.Order, error)
	PostOrder(context.Context, order.Order) (order.Order, error)
	UpdateOrder(context.Context, int, order.Order) (order.Order, error)
	DeleteOrder(context.Context, int) error
}

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

func (h *Handler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	ords, err := h.Service.GetAllOrders(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(ords); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PostOrder(w http.ResponseWriter, r *http.Request) {
	var ord order.Order
	if err := json.NewDecoder(r.Body).Decode(&ord); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	postedOrd, err := h.Service.PostOrder(r.Context(), ord)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(postedOrd); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

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

	var ord order.Order
	if err := json.NewDecoder(r.Body).Decode(&ord); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedOrd, err := h.Service.UpdateOrder(r.Context(), intOrderID, ord)
	if err != nil {
		if errors.Is(err, order.ErrNoOrderFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(updatedOrd); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

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

	if err := json.NewEncoder(w).Encode(Response{Message: "Successfully Deleted"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
