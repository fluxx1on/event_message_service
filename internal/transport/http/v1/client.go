package v1

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/fluxx1on_group/event_message_service/internal/entity"
	"gitlab.com/fluxx1on_group/event_message_service/internal/usecase"
)

const clientPath = basePath + "/client"

type clientRoutes struct {
	c usecase.Client
}

func newClientRoutes(handler *gin.RouterGroup, t usecase.Client) {
	r := &clientRoutes{t}

	h := handler.Group("/client")
	{
		h.PUT("/", r.Add)
		h.PATCH("/", r.Patch)
		h.DELETE("/", r.Delete)
	}
}

// @Summary 	Create a new client
// @Description Create new client entity.
// @ID 			createClient
// @Tags 		clients
// @Accept 		json
// @Produce 	json
// @Param 		client body entity.Client true "Client object to create"
// @Success 	201 "Client created successfully"
// @Failure 	400 "Bad request, invalid JSON data"
// @Failure 	500 "Internal server error, failed to create a client"
// @Router 		/client [put]
func (r *clientRoutes) Add(c *gin.Context) {
	var client entity.Client

	err := c.ShouldBindJSON(&client)
	if err != nil {
		slog.Warn("Unexpected request body",
			slog.Int("Status code", http.StatusBadRequest),
			slog.String("ErrorMsg", err.Error()))
		c.AbortWithStatus(http.StatusBadRequest)
		pushMetric(http.MethodPut, clientPath, http.StatusBadRequest)
		return
	}

	err = r.c.Add(c.Request.Context(), &client)
	if err != nil {
		slog.Info("Client creation failed",
			slog.Int("Status code", http.StatusInternalServerError),
			slog.String("ErrorMsg", err.Error()),
			slog.Group("Client", client.PhoneNumber, client.MobileOperator,
				client.Tag, client.TimeZone))
		c.AbortWithStatus(http.StatusInternalServerError)
		pushMetric(http.MethodPut, clientPath, http.StatusInternalServerError)
		return
	}

	slog.Info("Client creation succeeded",
		slog.Int("Status code", http.StatusOK),
		slog.Group("Client", client.PhoneNumber, client.MobileOperator,
			client.Tag, client.TimeZone))
	c.Status(http.StatusCreated)
	pushMetric(http.MethodPut, clientPath, http.StatusCreated)
}

// @Summary 	Update existing client
// @Description Update client in db.
// @ID 			updateClient
// @Tags 		clients
// @Accept 		json
// @Produce 	json
// @Param 		client body entity.Client true "Client object to update"
// @Success 	204 "Client updated successfully"
// @Failure 	400 "Bad request, invalid JSON data"
// @Failure 	500 "Internal server error, failed to update client"
// @Router 		/client [patch]
func (r *clientRoutes) Patch(c *gin.Context) {
	var client entity.Client

	err := c.ShouldBindJSON(&client)
	if err != nil {
		slog.Warn("Unexpected request body",
			slog.Int("Status code", http.StatusBadRequest),
			slog.String("ErrorMsg", err.Error()))
		c.AbortWithStatus(http.StatusBadRequest)
		pushMetric(http.MethodPatch, clientPath, http.StatusBadRequest)
		return
	}

	err = r.c.Patch(c.Request.Context(), &client)
	if err != nil {
		slog.Info("Client updating failed",
			slog.Int("Status code", http.StatusInternalServerError),
			slog.String("ErrorMsg", err.Error()),
			slog.Group("Client", client.PhoneNumber, client.MobileOperator,
				client.Tag, client.TimeZone))
		c.AbortWithStatus(http.StatusInternalServerError)
		pushMetric(http.MethodPatch, clientPath, http.StatusInternalServerError)
		return
	}

	slog.Info("Client updating succeeded",
		slog.Int("Status code", http.StatusOK),
		slog.Group("Client", client.PhoneNumber, client.MobileOperator,
			client.Tag, client.TimeZone))
	c.Status(http.StatusNoContent)
	pushMetric(http.MethodPatch, clientPath, http.StatusNoContent)
}

// @Summary 	Delete existing client
// @Description Delete client from db.
// @ID 			deleteClient
// @Tags 		clients
// @Accept 		json
// @Produce 	json
// @Param 		client body entity.Client true "Client object to delete"
// @Success 	204 "Client deleted successfully"
// @Failure 	400 "Bad request, invalid JSON data"
// @Failure 	500 "Internal server error, failed to delete client"
// @Router 		/client [delete]
func (r *clientRoutes) Delete(c *gin.Context) {
	var client entity.Client

	err := c.ShouldBindJSON(&client)
	if err != nil {
		slog.Warn("Unexpected request body",
			slog.Int("Status code", http.StatusBadRequest),
			slog.String("ErrorMsg", err.Error()))
		c.AbortWithStatus(http.StatusBadRequest)
		pushMetric(http.MethodDelete, clientPath, http.StatusBadRequest)
		return
	}

	err = r.c.Delete(c.Request.Context(), &client)
	if err != nil {
		slog.Info("Client deletion failed",
			slog.Int("Status code", http.StatusInternalServerError),
			slog.String("ErrorMsg", err.Error()),
			slog.Group("Client", client.PhoneNumber, client.MobileOperator,
				client.Tag, client.TimeZone))
		c.AbortWithStatus(http.StatusInternalServerError)
		pushMetric(http.MethodDelete, clientPath, http.StatusInternalServerError)
		return
	}

	slog.Info("Client deletion succeeded",
		slog.Int("Status code", http.StatusOK),
		slog.Group("Client", client.PhoneNumber, client.MobileOperator,
			client.Tag, client.TimeZone))
	c.Status(http.StatusNoContent)
	pushMetric(http.MethodDelete, clientPath, http.StatusNoContent)
}
