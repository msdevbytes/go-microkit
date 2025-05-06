package handler

import "github.com/gofiber/fiber/v2"

type DefaultHandlerStruct struct {
}

func NewDefaultHandler() *DefaultHandlerStruct {
	return &DefaultHandlerStruct{}
}

func (h *DefaultHandlerStruct) Register(router fiber.Router) {
	router.Get("/", h.Default)
}

func (h *DefaultHandlerStruct) Default(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
