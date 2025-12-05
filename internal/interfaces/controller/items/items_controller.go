package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	domainErrors "Aicon-assignment/internal/domain/errors"
	"Aicon-assignment/internal/usecase"

	"github.com/labstack/echo/v4"
)

const (
	// Immutable field names that cannot be updated via PATCH
	fieldID        = "id"
	fieldCreatedAt = "created_at"
	fieldUpdatedAt = "updated_at"
)

type ItemHandler struct {
	itemUsecase usecase.ItemUsecase
}

func NewItemHandler(itemUsecase usecase.ItemUsecase) *ItemHandler {
	return &ItemHandler{
		itemUsecase: itemUsecase,
	}
}

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Error   string   `json:"error"`
	Details []string `json:"details,omitempty"`
}

// parseItemID extracts and validates the item ID from the URL parameter
func parseItemID(idStr string) (int64, error) {
	if idStr == "" {
		return 0, strconv.ErrSyntax
	}
	return strconv.ParseInt(idStr, 10, 64)
}

// parseValidationErrorDetails extracts validation error details from a wrapped error
func parseValidationErrorDetails(err error) []string {
	details := []string{err.Error()}
	if strings.Contains(err.Error(), ": ") {
		parts := strings.SplitN(err.Error(), ": ", 2)
		if len(parts) == 2 {
			details = strings.Split(parts[1], ", ")
		}
	}
	return details
}

// checkImmutableFields validates that the request body doesn't contain immutable fields
func checkImmutableFields(requestBody map[string]interface{}) []string {
	var errors []string
	immutableFields := map[string]string{
		fieldID:        "id is immutable",
		fieldCreatedAt: "created_at is immutable",
		fieldUpdatedAt: "updated_at is immutable",
	}

	for field, message := range immutableFields {
		if _, exists := requestBody[field]; exists {
			errors = append(errors, message)
		}
	}

	return errors
}

func (h *ItemHandler) GetItems(c echo.Context) error {
	items, err := h.itemUsecase.GetAllItems(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to retrieve items",
		})
	}

	return c.JSON(http.StatusOK, items)
}

func (h *ItemHandler) GetItem(c echo.Context) error {
	id, err := parseItemID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid item ID",
		})
	}

	item, err := h.itemUsecase.GetItemByID(c.Request().Context(), id)
	if err != nil {
		if domainErrors.IsNotFoundError(err) {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "item not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to retrieve item",
		})
	}

	return c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) CreateItem(c echo.Context) error {
	var input usecase.CreateItemInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request format",
		})
	}

	// バリデーション
	if validationErrors := validateCreateItemInput(input); len(validationErrors) > 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation failed",
			Details: validationErrors,
		})
	}

	item, err := h.itemUsecase.CreateItem(c.Request().Context(), input)
	if err != nil {
		if domainErrors.IsValidationError(err) {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation failed",
				Details: []string{err.Error()},
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to create item",
		})
	}

	return c.JSON(http.StatusCreated, item)
}

func (h *ItemHandler) DeleteItem(c echo.Context) error {
	id, err := parseItemID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid item ID",
		})
	}

	err = h.itemUsecase.DeleteItem(c.Request().Context(), id)
	if err != nil {
		if domainErrors.IsNotFoundError(err) {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "item not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to delete item",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ItemHandler) GetSummary(c echo.Context) error {
	summary, err := h.itemUsecase.GetCategorySummary(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to retrieve summary",
		})
	}

	return c.JSON(http.StatusOK, summary)
}

func (h *ItemHandler) PatchItem(c echo.Context) error {
	id, err := parseItemID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid item ID",
		})
	}

	// Read and parse request body into a map first to check for immutable fields
	var requestBody map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request format",
		})
	}

	// Check for immutable fields
	if immutableErrors := checkImmutableFields(requestBody); len(immutableErrors) > 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation failed",
			Details: immutableErrors,
		})
	}

	// Parse into UpdateItemRequest struct
	var req usecase.UpdateItemRequest
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request format",
		})
	}

	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request format",
		})
	}

	item, err := h.itemUsecase.PatchItem(c.Request().Context(), id, &req)
	if err != nil {
		if domainErrors.IsNotFoundError(err) {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "item not found",
			})
		}
		if domainErrors.IsValidationError(err) {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation failed",
				Details: parseValidationErrorDetails(err),
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to update item",
		})
	}

	return c.JSON(http.StatusOK, item)
}

func validateCreateItemInput(input usecase.CreateItemInput) []string {
	var errs []string

	// Basic required field validation
	if input.Name == "" {
		errs = append(errs, "name is required")
	}
	if input.Category == "" {
		errs = append(errs, "category is required")
	}
	if input.Brand == "" {
		errs = append(errs, "brand is required")
	}
	if input.PurchaseDate == "" {
		errs = append(errs, "purchase_date is required")
	}
	if input.PurchasePrice < 0 {
		errs = append(errs, "purchase_price must be 0 or greater")
	}

	return errs
}
