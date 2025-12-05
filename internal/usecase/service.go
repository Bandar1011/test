package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"Aicon-assignment/internal/domain/entity"
	domainErrors "Aicon-assignment/internal/domain/errors"
)

const (
	maxNameLength  = 100
	maxBrandLength = 100
	minPrice       = 0
)

type ItemUsecase interface {
	GetAllItems(ctx context.Context) ([]*entity.Item, error)
	GetItemByID(ctx context.Context, id int64) (*entity.Item, error)
	CreateItem(ctx context.Context, input CreateItemInput) (*entity.Item, error)
	DeleteItem(ctx context.Context, id int64) error
	PatchItem(ctx context.Context, id int64, req *UpdateItemRequest) (*entity.Item, error)
	GetCategorySummary(ctx context.Context) (*CategorySummary, error)
}

type CreateItemInput struct {
	Name          string `json:"name"`
	Category      string `json:"category"`
	Brand         string `json:"brand"`
	PurchasePrice int    `json:"purchase_price"`
	PurchaseDate  string `json:"purchase_date"`
}

type UpdateItemRequest struct {
	Name          *string `json:"name,omitempty"`
	Brand         *string `json:"brand,omitempty"`
	PurchasePrice *int    `json:"purchase_price,omitempty"`
}

type CategorySummary struct {
	Categories map[string]int `json:"categories"`
	Total      int            `json:"total"`
}

type itemUsecase struct {
	itemRepo ItemRepository
}

func NewItemUsecase(itemRepo ItemRepository) ItemUsecase {
	return &itemUsecase{
		itemRepo: itemRepo,
	}
}

func (u *itemUsecase) GetAllItems(ctx context.Context) ([]*entity.Item, error) {
	items, err := u.itemRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve items: %w", err)
	}

	return items, nil
}

func (u *itemUsecase) GetItemByID(ctx context.Context, id int64) (*entity.Item, error) {
	if id <= 0 {
		return nil, domainErrors.ErrInvalidInput
	}

	item, err := u.itemRepo.FindByID(ctx, id)
	if err != nil {
		if domainErrors.IsNotFoundError(err) {
			return nil, domainErrors.ErrItemNotFound
		}
		return nil, fmt.Errorf("failed to retrieve item: %w", err)
	}

	return item, nil
}

func (u *itemUsecase) CreateItem(ctx context.Context, input CreateItemInput) (*entity.Item, error) {
	// バリデーションして、新しいエンティティを作成
	item, err := entity.NewItem(
		input.Name,
		input.Category,
		input.Brand,
		input.PurchasePrice,
		input.PurchaseDate,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", domainErrors.ErrInvalidInput, err.Error())
	}

	createdItem, err := u.itemRepo.Create(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return createdItem, nil
}

func (u *itemUsecase) DeleteItem(ctx context.Context, id int64) error {
	if id <= 0 {
		return domainErrors.ErrInvalidInput
	}

	_, err := u.itemRepo.FindByID(ctx, id)
	if err != nil {
		if domainErrors.IsNotFoundError(err) {
			return domainErrors.ErrItemNotFound
		}
		return fmt.Errorf("failed to check item existence: %w", err)
	}

	err = u.itemRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

func (u *itemUsecase) PatchItem(ctx context.Context, id int64, req *UpdateItemRequest) (*entity.Item, error) {
	if id <= 0 {
		return nil, domainErrors.ErrInvalidInput
	}

	// Fetch existing item
	item, err := u.itemRepo.FindByID(ctx, id)
	if err != nil {
		if domainErrors.IsNotFoundError(err) {
			return nil, domainErrors.ErrItemNotFound
		}
		return nil, fmt.Errorf("failed to retrieve item: %w", err)
	}

	// Apply partial updates
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Brand != nil {
		item.Brand = *req.Brand
	}
	if req.PurchasePrice != nil {
		item.PurchasePrice = *req.PurchasePrice
	}

	// Update timestamp
	item.UpdatedAt = time.Now()

	// Validate updated fields
	if validationErrors := validateUpdateRequest(req, item); len(validationErrors) > 0 {
		return nil, fmt.Errorf("%w: %s", domainErrors.ErrInvalidInput, strings.Join(validationErrors, ", "))
	}

	// Save updated item
	updatedItem, err := u.itemRepo.Update(ctx, item)
	if err != nil {
		if domainErrors.IsNotFoundError(err) {
			return nil, domainErrors.ErrItemNotFound
		}
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return updatedItem, nil
}

func (u *itemUsecase) GetCategorySummary(ctx context.Context) (*CategorySummary, error) {
	categoryCounts, err := u.itemRepo.GetSummaryByCategory(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get category summary: %w", err)
	}

	// 合計計算
	total := 0
	for _, count := range categoryCounts {
		total += count
	}

	summary := make(map[string]int)
	for _, category := range entity.GetValidCategories() {
		if count, exists := categoryCounts[category]; exists {
			summary[category] = count
		} else {
			summary[category] = 0
		}
	}

	return &CategorySummary{
		Categories: summary,
		Total:      total,
	}, nil
}

// validateUpdateRequest validates the fields being updated in a PATCH request
func validateUpdateRequest(req *UpdateItemRequest, item *entity.Item) []string {
	var validationErrors []string

	if req.Name != nil {
		if item.Name == "" {
			validationErrors = append(validationErrors, "name is required")
		} else if len(item.Name) > maxNameLength {
			validationErrors = append(validationErrors, fmt.Sprintf("name must be %d characters or less", maxNameLength))
		}
	}

	if req.Brand != nil {
		if item.Brand == "" {
			validationErrors = append(validationErrors, "brand is required")
		} else if len(item.Brand) > maxBrandLength {
			validationErrors = append(validationErrors, fmt.Sprintf("brand must be %d characters or less", maxBrandLength))
		}
	}

	if req.PurchasePrice != nil {
		if item.PurchasePrice < minPrice {
			validationErrors = append(validationErrors, fmt.Sprintf("purchase_price must be >= %d", minPrice))
		}
	}

	return validationErrors
}
