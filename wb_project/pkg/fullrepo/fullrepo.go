package fullrepo

import (
	"wb_project/pkg/delivery"
	"wb_project/pkg/items"
	"wb_project/pkg/pay"
	"wb_project/pkg/user"
)

type FullRepo struct {
	RepoU user.Repository
	RepoP pay.Repository
	RepoI items.Repository
	RepoD delivery.Repository
}
