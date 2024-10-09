package repository

import "context"

type ImageRepository interface {
	GetImage(ctx context.Context, id string) (*MediaFile, error)
	DeleteImage(ctx context.Context, userId string) error
}
