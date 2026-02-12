package services

import menuv1 "github.com/defvova/go_projects/mini_food_delivery/menu/pkg/menu/v1"

type BaseService struct {
	MenuClient menuv1.MenuServiceClient
}
