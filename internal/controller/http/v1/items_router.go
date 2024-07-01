package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"warehousesAPI/internal/entity"
	"warehousesAPI/internal/usecase"
	"warehousesAPI/pkg/logger"
)

type itemsAPIRouter struct {
	items usecase.ItemsUseCase
	l     logger.Interface
}

func newItemsAPIRoutes(handler *gin.RouterGroup, i usecase.ItemsUseCase, l logger.Interface) {
	items := &itemsAPIRouter{
		items: i,
		l:     l,
	}

	h := handler.Group("/items")
	{
		h.GET("/:warehouse_id/quantity", items.getItemsQuantity)
		h.PUT("", items.createItems)
	}
}

// @Summary     Get items quantity
// @Description Count items in warehouse
// @ID          getItemsQuantity
// @Tags  	    itmes
// @Accept      json
// @Produce     json
// @Param 		warehouse_id	path 		int			true 	"warehouse_id"
// @Success     201 			{object} 	response
// @Failure		400				{object}	response
// @Failure     500 			{object} 	response
// @Router      /items/{warehouse_id}/quantity 			[get]
func (r *itemsAPIRouter) getItemsQuantity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("warehouse_id"))
	if err != nil {
		r.l.Error(err, "error converting warehouse_id to int")
		errorResponse(c, http.StatusBadRequest, "error converting warehouse_id to int")
	}

	q, err := r.items.Quantity(c.Request.Context(), id)
	if err != nil {
		r.l.Error(err, "error getting items quantity")
		errorResponse(c, http.StatusInternalServerError, "error getting items quantity")
	}

	successResponse(c, http.StatusOK, q)
}

type itemsCreateRequest struct {
	Items []entity.Item `json:"items"`
}

// @Summary     Create items
// @Description Create items in warehouse
// @ID          createItem
// @Tags  	    itmes
// @Accept      json
// @Produce     json
// @Param 		item	 		body 		itemsCreateRequest		true 	"items"
// @Success     201 			{object} 	response
// @Failure		400				{object}	response
// @Failure     500 			{object} 	response
// @Router      /items 			[put]
func (r *itemsAPIRouter) createItems(c *gin.Context) {
	var itemsReq itemsCreateRequest
	if err := c.BindJSON(&itemsReq); err != nil {
		r.l.Error(err, "error binding JSON")
		errorResponse(c, http.StatusBadRequest, "provided data is invalid")
	}

	if err := r.items.CreateItems(c.Request.Context(), itemsReq.Items); err != nil {
		r.l.Error(err, "failed to create item")
		errorResponse(c, http.StatusInternalServerError, "items service problems")

		return
	}

	successResponse(c, http.StatusCreated, "items successfully created")
}