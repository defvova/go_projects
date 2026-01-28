package handler

import (
	"shortener/db"
	"shortener/internal/service"
)

type Handler struct {
	Q         *db.Queries
	NewRender *service.Renderer
}
