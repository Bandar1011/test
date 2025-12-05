package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Aicon-assignment/internal/domain/entity"
	domainErrors "Aicon-assignment/internal/domain/errors"
	"Aicon-assignment/internal/usecase"
)

// MockItemUsecase is a mock implementation of ItemUsecase for testing
type MockItemUsecase struct {
	mock.Mock
}

func (m *MockItemUsecase) GetAllItems(ctx context.Context) ([]*entity.Item, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Item), args.Error(1)
}

func (m *MockItemUsecase) GetItemByID(ctx context.Context, id int64) (*entity.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockItemUsecase) CreateItem(ctx context.Context, input usecase.CreateItemInput) (*entity.Item, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockItemUsecase) DeleteItem(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockItemUsecase) PatchItem(ctx context.Context, id int64, req *usecase.UpdateItemRequest) (*entity.Item, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockItemUsecase) GetCategorySummary(ctx context.Context) (*usecase.CategorySummary, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.CategorySummary), args.Error(1)
}

func TestItemHandler_PatchItem(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name           string
		itemID         string
		requestBody    map[string]interface{}
		setupMock      func(*MockItemUsecase)
		expectedStatus int
		expectedError  string
		expectedDetails []string
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "Success - update name",
			itemID: "1",
			requestBody: map[string]interface{}{
				"name": "Updated Item Name",
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				updatedItem, _ := entity.NewItem("Updated Item Name", "時計", "ROLEX", 1000000, "2023-01-01")
				updatedItem.ID = 1
				updatedItem.CreatedAt = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
				updatedItem.UpdatedAt = time.Now()

				req := &usecase.UpdateItemRequest{
					Name: stringPtr("Updated Item Name"),
				}
				mockUsecase.On("PatchItem", mock.Anything, int64(1), req).Return(updatedItem, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var item entity.Item
				err := json.Unmarshal(rec.Body.Bytes(), &item)
				assert.NoError(t, err)
				assert.Equal(t, "Updated Item Name", item.Name)
				assert.Equal(t, int64(1), item.ID)
			},
		},
		{
			name:   "Success - update purchase_price",
			itemID: "1",
			requestBody: map[string]interface{}{
				"purchase_price": 2000000,
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				updatedItem, _ := entity.NewItem("時計1", "時計", "ROLEX", 2000000, "2023-01-01")
				updatedItem.ID = 1
				updatedItem.CreatedAt = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
				updatedItem.UpdatedAt = time.Now()

				req := &usecase.UpdateItemRequest{
					PurchasePrice: intPtr(2000000),
				}
				mockUsecase.On("PatchItem", mock.Anything, int64(1), req).Return(updatedItem, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var item entity.Item
				err := json.Unmarshal(rec.Body.Bytes(), &item)
				assert.NoError(t, err)
				assert.Equal(t, 2000000, item.PurchasePrice)
			},
		},
		{
			name:   "Success - update brand only",
			itemID: "1",
			requestBody: map[string]interface{}{
				"brand": "Updated Brand Name",
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				updatedItem, _ := entity.NewItem("時計1", "時計", "Updated Brand Name", 1000000, "2023-01-01")
				updatedItem.ID = 1
				updatedItem.CreatedAt = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
				updatedItem.UpdatedAt = time.Now()

				req := &usecase.UpdateItemRequest{
					Brand: stringPtr("Updated Brand Name"),
				}
				mockUsecase.On("PatchItem", mock.Anything, int64(1), req).Return(updatedItem, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var item entity.Item
				err := json.Unmarshal(rec.Body.Bytes(), &item)
				assert.NoError(t, err)
				assert.Equal(t, "Updated Brand Name", item.Brand)
				assert.Equal(t, int64(1), item.ID)
			},
		},
		{
			name:   "Success - update multiple fields",
			itemID: "1",
			requestBody: map[string]interface{}{
				"name":           "New Name",
				"brand":          "New Brand",
				"purchase_price": 1500000,
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				updatedItem, _ := entity.NewItem("New Name", "時計", "New Brand", 1500000, "2023-01-01")
				updatedItem.ID = 1
				updatedItem.CreatedAt = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
				updatedItem.UpdatedAt = time.Now()

				req := &usecase.UpdateItemRequest{
					Name:          stringPtr("New Name"),
					Brand:         stringPtr("New Brand"),
					PurchasePrice: intPtr(1500000),
				}
				mockUsecase.On("PatchItem", mock.Anything, int64(1), req).Return(updatedItem, nil)
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var item entity.Item
				err := json.Unmarshal(rec.Body.Bytes(), &item)
				assert.NoError(t, err)
				assert.Equal(t, "New Name", item.Name)
				assert.Equal(t, "New Brand", item.Brand)
				assert.Equal(t, 1500000, item.PurchasePrice)
			},
		},
		{
			name:   "404 - item not found",
			itemID: "9999",
			requestBody: map[string]interface{}{
				"name": "Updated Name",
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				req := &usecase.UpdateItemRequest{
					Name: stringPtr("Updated Name"),
				}
				mockUsecase.On("PatchItem", mock.Anything, int64(9999), req).Return((*entity.Item)(nil), domainErrors.ErrItemNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "item not found",
		},
		{
			name:   "400 - invalid price (negative)",
			itemID: "1",
			requestBody: map[string]interface{}{
				"purchase_price": -100,
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				req := &usecase.UpdateItemRequest{
					PurchasePrice: intPtr(-100),
				}
				err := fmt.Errorf("%w: %s", domainErrors.ErrInvalidInput, "purchase_price must be >= 0")
				mockUsecase.On("PatchItem", mock.Anything, int64(1), req).Return((*entity.Item)(nil), err)
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "validation failed",
			expectedDetails: []string{"purchase_price must be >= 0"},
		},
		{
			name:   "400 - immutable field (id)",
			itemID: "1",
			requestBody: map[string]interface{}{
				"id":   999,
				"name": "Updated Name",
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				// Mock should not be called when immutable field is present
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "validation failed",
			expectedDetails: []string{"id is immutable"},
		},
		{
			name:   "400 - immutable field (created_at)",
			itemID: "1",
			requestBody: map[string]interface{}{
				"created_at": "2023-01-01T00:00:00Z",
				"name":       "Updated Name",
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				// Mock should not be called when immutable field is present
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "validation failed",
			expectedDetails: []string{"created_at is immutable"},
		},
		{
			name:   "400 - immutable field (updated_at)",
			itemID: "1",
			requestBody: map[string]interface{}{
				"updated_at": "2023-01-01T00:00:00Z",
				"name":       "Updated Name",
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				// Mock should not be called when immutable field is present
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "validation failed",
			expectedDetails: []string{"updated_at is immutable"},
		},
		{
			name:   "400 - multiple immutable fields",
			itemID: "1",
			requestBody: map[string]interface{}{
				"id":         999,
				"created_at": "2023-01-01T00:00:00Z",
				"name":       "Updated Name",
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				// Mock should not be called when immutable fields are present
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "validation failed",
			expectedDetails: []string{"id is immutable", "created_at is immutable"},
		},
		{
			name:   "400 - invalid item ID",
			itemID: "invalid",
			requestBody: map[string]interface{}{
				"name": "Updated Name",
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				// Mock should not be called when ID is invalid
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid item ID",
		},
		{
			name:   "400 - name too long",
			itemID: "1",
			requestBody: map[string]interface{}{
				"name": string(make([]byte, 101)), // 101 characters
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				longName := string(make([]byte, 101))
				req := &usecase.UpdateItemRequest{
					Name: &longName,
				}
				err := fmt.Errorf("%w: %s", domainErrors.ErrInvalidInput, "name must be 100 characters or less")
				mockUsecase.On("PatchItem", mock.Anything, int64(1), req).Return((*entity.Item)(nil), err)
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "validation failed",
			expectedDetails: []string{"name must be 100 characters or less"},
		},
		{
			name:   "400 - brand too long",
			itemID: "1",
			requestBody: map[string]interface{}{
				"brand": string(make([]byte, 101)), // 101 characters
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				longBrand := string(make([]byte, 101))
				req := &usecase.UpdateItemRequest{
					Brand: &longBrand,
				}
				err := fmt.Errorf("%w: %s", domainErrors.ErrInvalidInput, "brand must be 100 characters or less")
				mockUsecase.On("PatchItem", mock.Anything, int64(1), req).Return((*entity.Item)(nil), err)
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "validation failed",
			expectedDetails: []string{"brand must be 100 characters or less"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockItemUsecase)
			tt.setupMock(mockUsecase)

			handler := &ItemHandler{
				itemUsecase: mockUsecase,
			}

			bodyBytes, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/items/"+tt.itemID, bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/items/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.itemID)

			err = handler.PatchItem(c)

			// Echo handles errors automatically, so we check the response code
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedError != "" {
				var errorResp ErrorResponse
				err := json.Unmarshal(rec.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResp.Error)
				if len(tt.expectedDetails) > 0 {
					assert.Equal(t, tt.expectedDetails, errorResp.Details)
				}
			}

			if tt.validateResponse != nil {
				tt.validateResponse(t, rec)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

