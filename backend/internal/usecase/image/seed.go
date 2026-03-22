package image

import (
	"context"

	"example.com/m/internal/domain/image"
)

// VMのイメージ初期データ
var seedImages = []*image.Image{
	image.NewImage("img-ubuntu-2404", "ubuntu/24.04", "", "https://images.linuxcontainers.org", "simplestreams", true),
	// 他のイメージ...
}

type SeedImageUseCase interface {
	Execute(ctx context.Context) error
}

type seedImageInteractor struct {
	imageRepo image.Repository
}

func NewSeedImageUseCase(imageRepo image.Repository) SeedImageUseCase {
	return &seedImageInteractor{
		imageRepo: imageRepo,
	}
}

func (uc *seedImageInteractor) Execute(ctx context.Context) error {
	// べき等性担保のための存在確認
	existingImages, err := uc.imageRepo.FindAll(ctx)
	if err != nil {
		return err
	}

	// 既存のイメージと比較して、存在しない場合は作成
	// NOTE: パフォーマンスを考えるとBulkがいいけど、初期Initなので、あんまり気にしない
	for _, seedImage := range seedImages {
		found := false
		for _, existingImage := range existingImages {
			if existingImage.ID() == seedImage.ID() {
				found = true
				break
			}
		}
		if !found {
			if err := uc.imageRepo.Save(ctx, seedImage); err != nil {
				return err
			}
		}
	}
	return nil
}
