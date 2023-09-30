package v1

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/fluxx1on_group/event_message_service/internal/entity"
	"gitlab.com/fluxx1on_group/event_message_service/internal/usecase"
)

const mailingPath = basePath + "/mailing"
const mailingStatsPath = mailingPath + "/stats"

type mailingRoutes struct {
	m usecase.Mailing
}

func newMailingRoutes(handler *gin.RouterGroup, m usecase.Mailing) {
	r := &mailingRoutes{m}

	h := handler.Group("/mailing")
	{
		h.GET("/stats", r.GetStats)
		h.POST("/", r.ReadMessages)
		h.PUT("/", r.Add)
		h.PATCH("/", r.Patch)
		h.DELETE("/", r.Delete)
	}
}

// @Summary 	Get MailingStats
// @Description Get MailingStats about all Mailings.
// @ID 			getStats
// @Tags 		mailings
// @Accept 		json
// @Produce 	json
// @Success  	200 {object} []entity.MailingStats "MailingStats received"
// @Failure 	500 {object} errorResponse "Internal server error, failed to receive stats"
// @Router 		/mailing/stats [get]
func (r *mailingRoutes) GetStats(c *gin.Context) {
	stats, err := r.m.GetMailingStats(c.Request.Context())
	if err != nil {
		slog.Info("MailingStats reading failed",
			slog.Int("Status code", http.StatusInternalServerError),
			slog.String("ErrorMsg", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{
			ErrorMsg: "Internal server error, failed to receive stats",
		})
		pushMetric(http.MethodGet, mailingStatsPath, http.StatusInternalServerError)
		return
	}

	slog.Info("MailingStats reading succeeded",
		slog.Int("Status code", http.StatusOK))
	c.JSON(http.StatusOK, stats)
	pushMetric(http.MethodGet, mailingStatsPath, http.StatusOK)
}

// @Summary 	Post with Mailing
// @Description Get Messages by existing mailing.
// @ID 			getMessagesByMailing
// @Tags 		mailings
// @Accept 		json
// @Produce 	json
// @Param 		mailing body entity.Mailing true "Mailing object to select Messages by mailing's filters"
// @Success  	200 {object} entity.Messages "Messages catched successfully"
// @Failure 	400 {object} errorResponse "Bad request, invalid JSON data"
// @Failure 	500 {object} errorResponse "Internal server error, failed to catch messages"
// @Router 		/mailing [post]
func (r *mailingRoutes) ReadMessages(c *gin.Context) {
	var mailing entity.Mailing

	err := c.ShouldBindJSON(&mailing)
	if err != nil {
		slog.Warn("Unexpected request body",
			slog.Int("Status code", http.StatusBadRequest),
			slog.String("ErrorMsg", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{
			ErrorMsg: "Bad request, invalid JSON data",
		})
		pushMetric(http.MethodPost, mailingPath, http.StatusBadRequest)
		return
	}

	msgs, err := r.m.GetMessagesByMailing(c.Request.Context(), &mailing)
	if err != nil {
		slog.Info("Messages reading failed",
			slog.Int("Status code", http.StatusInternalServerError),
			slog.String("ErrorMsg", err.Error()),
			slog.Group("Mailing", mailing.ID, mailing.Tag, mailing.MobileOperator,
				mailing.MessageText[0:40], mailing.DateTimeStart))
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{
			ErrorMsg: "Internal server error, failed to catch messages",
		})
		pushMetric(http.MethodPost, mailingPath, http.StatusInternalServerError)
		return
	}

	slog.Info("Messages reading succeeded",
		slog.Int("Status code", http.StatusOK),
		slog.Group("Mailing", mailing.ID, mailing.Tag, mailing.MobileOperator,
			mailing.MessageText[0:40], mailing.DateTimeStart))
	c.JSON(http.StatusOK, msgs)
	pushMetric(http.MethodPost, mailingPath, http.StatusOK)
}

// @Summary 	Create a mailing
// @Description Create a new mailing.
// @ID 			createMailing
// @Tags 		mailings
// @Accept 		json
// @Produce 	json
// @Param 		mailing body entity.Mailing true "Mailing object to create"
// @Success 	201 "Mailing created successfully"
// @Failure 	400 "Bad request, invalid JSON data"
// @Failure 	500 "Internal server error, failed to create mailing"
// @Router 		/mailing [put]
func (r *mailingRoutes) Add(c *gin.Context) {
	var mailing entity.Mailing

	err := c.ShouldBindJSON(&mailing)
	if err != nil {
		slog.Warn("Unexpected request body",
			slog.Int("Status code", http.StatusBadRequest),
			slog.String("ErrorMsg", err.Error()))
		c.AbortWithStatus(http.StatusBadRequest)
		pushMetric(http.MethodPut, mailingPath, http.StatusBadRequest)
		return
	}

	err = r.m.Add(c.Request.Context(), &mailing)
	if err != nil {
		slog.Info("Mailing creation failed",
			slog.Int("Status code", http.StatusInternalServerError),
			slog.String("ErrorMsg", err.Error()),
			slog.Group("Mailing", mailing.Tag, mailing.MobileOperator,
				mailing.MessageText[0:40], mailing.DateTimeStart))
		c.AbortWithStatus(http.StatusInternalServerError)
		pushMetric(http.MethodPut, mailingPath, http.StatusInternalServerError)
		return
	}

	slog.Info("Mailing creation succeeded",
		slog.Int("Status code", http.StatusOK),
		slog.Group("Mailing", mailing.Tag, mailing.MobileOperator,
			mailing.MessageText[0:40], mailing.DateTimeStart))
	c.Status(http.StatusCreated)
	pushMetric(http.MethodPut, mailingPath, http.StatusCreated)
}

// @Summary 	Update existing mailing
// @Description Update mailing in db.
// @ID 			updateMailing
// @Tags 		mailings
// @Accept 		json
// @Produce 	json
// @Param 		mailing body entity.Mailing true "Mailing object to update"
// @Success 	204 "Mailing updated successfully"
// @Failure 	400 "Bad request, invalid JSON data"
// @Failure 	500 "Internal server error, failed to update mailing"
// @Router 		/mailing [patch]
func (r *mailingRoutes) Patch(c *gin.Context) {
	var mailing entity.Mailing

	err := c.ShouldBindJSON(&mailing)
	if err != nil {
		slog.Warn("Unexpected request body",
			slog.Int("Status code", http.StatusBadRequest),
			slog.String("ErrorMsg", err.Error()))
		c.AbortWithStatus(http.StatusBadRequest)
		pushMetric(http.MethodPatch, mailingPath, http.StatusBadRequest)
		return
	}

	err = r.m.Patch(c.Request.Context(), &mailing)
	if err != nil {
		slog.Info("Mailing updating failed",
			slog.Int("Status code", http.StatusInternalServerError),
			slog.String("ErrorMsg", err.Error()),
			slog.Group("Mailing", mailing.Tag, mailing.MobileOperator,
				mailing.MessageText[0:40], mailing.DateTimeStart))
		c.AbortWithStatus(http.StatusInternalServerError)
		pushMetric(http.MethodPatch, mailingPath, http.StatusInternalServerError)
		return
	}

	slog.Info("Mailing updating succeeded",
		slog.Int("Status code", http.StatusOK),
		slog.Group("Mailing", mailing.Tag, mailing.MobileOperator,
			mailing.MessageText[0:40], mailing.DateTimeStart))
	c.Status(http.StatusNoContent)
	pushMetric(http.MethodPatch, mailingPath, http.StatusNoContent)
}

// @Summary 	Delete existing mailing
// @Description Delete mailing from db.
// @ID 			deleteMailing
// @Tags 		mailings
// @Accept 		json
// @Produce 	json
// @Param 		mailing body entity.Mailing true "Mailing object to delete"
// @Success 	204 "Mailing deleted successfully"
// @Failure 	400 "Bad request, invalid JSON data"
// @Failure 	500 "Internal server error, failed to delete mailing"
// @Router 		/mailing [delete]
func (r *mailingRoutes) Delete(c *gin.Context) {
	var mailing entity.Mailing

	err := c.ShouldBindJSON(&mailing)
	if err != nil {
		slog.Warn("Unexpected request body",
			slog.Int("Status code", http.StatusBadRequest),
			slog.String("ErrorMsg", err.Error()))
		c.AbortWithStatus(http.StatusBadRequest)
		pushMetric(http.MethodDelete, mailingPath, http.StatusBadRequest)
		return
	}

	err = r.m.Delete(c.Request.Context(), &mailing)
	if err != nil {
		slog.Info("Mailing deletion failed",
			slog.Int("Status code", http.StatusInternalServerError),
			slog.String("ErrorMsg", err.Error()),
			slog.Group("Mailing", mailing.Tag, mailing.MobileOperator,
				mailing.MessageText[0:40], mailing.DateTimeStart))
		c.AbortWithStatus(http.StatusInternalServerError)
		pushMetric(http.MethodDelete, mailingPath, http.StatusInternalServerError)
		return
	}

	slog.Info("Mailing deletion succeeded",
		slog.Int("Status code", http.StatusOK),
		slog.Group("Mailing", mailing.Tag, mailing.MobileOperator,
			mailing.MessageText[0:40], mailing.DateTimeStart))
	c.Status(http.StatusNoContent)
	pushMetric(http.MethodDelete, mailingPath, http.StatusNoContent)
}
