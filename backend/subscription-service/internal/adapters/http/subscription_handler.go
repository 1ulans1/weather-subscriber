package http

import (
	"context"
	"subscription-service/internal/core/services"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc *services.SubscriptionService
}

func NewHandler(svc *services.SubscriptionService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/weather", h.GetWeather)
	api.Post("/subscribe", h.Subscribe)
	api.Get("/confirm/:token", h.Confirm)
	api.Get("/unsubscribe/:token", h.Unsubscribe)
}

func (h *Handler) GetWeather(c *fiber.Ctx) error {
	city := c.Query("city")
	if city == "" {
		return c.Status(400).JSON(fiber.Map{"error": "city is required"})
	}
	w, err := h.svc.GetWeather(context.Background(), city)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "city not found"})
	}
	return c.JSON(fiber.Map{
		"temperature": w.Temperature,
		"description": w.Condition,
	})
}

func (h *Handler) Subscribe(c *fiber.Ctx) error {
	var req SubscribeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}
	err := h.svc.Subscribe(c.Context(), req.Email, req.City, req.Frequency)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Subscription successful. Confirmation email sent."})
}

func (h *Handler) Confirm(c *fiber.Ctx) error {
	tok := c.Params("token")
	if err := h.svc.Confirm(c.Context(), tok); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Subscription confirmed successfully"})
}

func (h *Handler) Unsubscribe(c *fiber.Ctx) error {
	tok := c.Params("token")
	if err := h.svc.Unsubscribe(c.Context(), tok); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Unsubscribed successfully"})
}
