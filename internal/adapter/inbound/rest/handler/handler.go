package handler

import "synapsis-challenge/internal/core/port/inbound/registry"

type Handler struct {
	serviceReg registry.ServiceRegistry
}

func New(reg registry.ServiceRegistry) *Handler {
	return &Handler{
		serviceReg: reg,
	}
}

func (h *Handler) GetServiceRegistry() registry.ServiceRegistry {
	return h.serviceReg
}
