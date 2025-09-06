package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	microservice "github.com/ikotiki/go-rest-api-service-subscriptions"

	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/service"
	"github.com/ikotiki/go-rest-api-service-subscriptions/pkg/e"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	sub service.Subscriptions
}

func NewSubscriptionHandler(g *gin.RouterGroup, service service.Subscriptions) *SubscriptionHandler {
	a := &SubscriptionHandler{
		sub: service,
	}
	a.registerRoutes(g)
	return a
}

func (a *SubscriptionHandler) registerRoutes(g *gin.RouterGroup) {
	sub := g.Group("/subscription")
	{
		sub.GET("/:id", a.getSubscriptionByID)
		sub.POST("/", a.createSubscription)
		sub.PUT("/:id", a.updateSubscription)
		sub.DELETE("/:id", a.deleteSubscription)

		sub.GET("/query", a.querySubscriptions)
		sub.GET("/sum", a.sumSubscriptions)
	}

}

// getSubscriptionByID godoc
// @Summary      Subscription By ID
// @Description  Get subscription by its id
// @Tags         subscriptions
// @Produce      json
// @Param        id    path     microservice.SubscriptionID true  "id of the subscription"  minimum(1)    maximum(10)
// @Success      200  {object}  respSuc{obj=microservice.Subscription}
// @Failure      400  {object}  respErr
// @Failure      404  {object}  respErr
// @Failure      500  {object}  respErr
// @Router       /subscription/{id}	 [get]
func (a *SubscriptionHandler) getSubscriptionByID(c *gin.Context) {
	const op = "handler.getSubscriptionByID"
	log, ctx := prepareTools(c, op)

	// Parse ID
	id, err := a.parseSubscriptionID(c)
	if err != nil {
		log.Debug().Err(err).Str("id", c.Param("id")).Msg("can't parse subscription id")
		writeBadRequest(c, err.Error())
		return
	}

	log.Debug().Int("id", int(id)).Msg("subscription id")

	// Get subscription
	sub, err := a.sub.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrNoSuchSubscription) {
			writeNotFound(c, "no such subscription")
			return
		}
		// logger.Error("error getting subscription", err)
		log.Error().Err(err).Msg("error getting subscription")
		writeServerInternal(c, "error getting subscription on the server")
		return
	}

	log.Info().Interface("subscription", sub).Msg("got subscription")

	writeObj(c, sub)
}

// createSubscription godoc
// @Summary      Create Subscription
// @Description  Create a new subscription
// @Tags         subscriptions
// @Produce      json
// @Param        subscription  body     microservice.Subscription  true  "subscription object"
// @Success      201  {object}  respSuc{obj=microservice.SubscriptionID}
// @Failure      400  {object}  respErr
// @Failure      422  {object}  respErr
// @Failure      500  {object}  respErr
// @Router       /subscription/	 [post]
func (a *SubscriptionHandler) createSubscription(c *gin.Context) {
	const op = "handler.createSubscription"
	log, ctx := prepareTools(c, op)

	sub := &microservice.Subscription{}
	if err := c.ShouldBindJSON(sub); err != nil {
		writeBadRequest(c, "error binding json: "+err.Error())
		return
	}

	id, err := a.sub.Create(ctx, sub)
	if err != nil {
		if errors.Is(err, service.ErrUserSubscriptionPairAlreadyExists) {
			writeFailure(c, http.StatusUnprocessableEntity, "user-subscription pair already exists", nil)
			return
		}
		log.Error().Err(err).Msg("error creating subscription")
		writeServerInternal(c, "error creating subscription on the server")
		return
	}

	log.Info().Int("id", int(id)).Msg("subscription created")

	writeSuccess(c, http.StatusCreated, msgSuccess, id)
}

// updateSubscription godoc
// @Summary      Update Subscription
// @Description  Update the subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id    path     microservice.SubscriptionID  true  "id of the subscription"  minimum(1)    maximum(10)
// @Param        subscription  body     microservice.Subscription  true  "subscription object"
// @Success      200  {object}  respSucNoObj
// @Failure      400  {object}  respErr
// @Failure      422  {object}  respErr
// @Failure      500  {object}  respErr
// @Router       /subscription/{id}	 [put]
func (a *SubscriptionHandler) updateSubscription(c *gin.Context) {
	const op = "handler.updateSubscription"
	log, ctx := prepareTools(c, op)

	sub := &microservice.Subscription{}
	if err := c.ShouldBindJSON(sub); err != nil {
		writeBadRequest(c, "error binding json: "+err.Error())
		return
	}

	err := a.sub.Update(ctx, sub)
	if err != nil {
		if errors.Is(err, service.ErrNoSuchSubscription) {
			writeNotFound(c, "no such subscription")
		}
		log.Error().Err(err).Msg("error updating subscription")
		writeServerInternal(c, "error updating subscription")
		return
	}

	writeOK(c)
}

// deleteSubscription godoc
// @Summary      Delete Subscription
// @Description  Delete a new subscription
// @Tags         subscriptions
// @Produce      json
// @Param        id    path     microservice.SubscriptionID  true  "id of the subscription"  minimum(1)    maximum(10)
// @Success      204  {object}  respSucNoObj
// @Failure      400  {object}  respErr
// @Failure      404  {object}  respErr
// @Failure      500  {object}  respErr
// @Router       /subscription/{id}	 [delete]
func (a *SubscriptionHandler) deleteSubscription(c *gin.Context) {
	const op = "handler.deleteSubscription"
	log, ctx := prepareTools(c, op)

	// Parse ID
	id, err := a.parseSubscriptionID(c)
	if err != nil {
		log.Debug().Err(err).Str("id", c.Param("id")).Msg("can't parse subscription id")
		writeBadRequest(c, "can't parse subscription id"+err.Error())
		return
	}

	log.Debug().Int("id", int(id)).Msg("subscription id")

	err = a.sub.DeleteByID(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrNoSuchSubscription) {
			log.Debug().Err(err).Msg("no such subscription")
			writeNotFound(c, "no such subscription")
			return
		}
		log.Error().Err(err).Msg("error deleting subscription")
		writeServerInternal(c, "error deleting subscription on the server")
		return
	}

	log.Info().Int("id", int(id)).Msg("subscription deleted")

	writeSuccess(c, http.StatusNoContent, msgSuccess, nil)
}

// querySubscriptions godoc
// @Summary      Get Subscriptions
// @Description  Get Subscriptions by a query
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        query  body     service.SubscriptionQueryArgs  true  "query arguments"
// @Success      200  {object}  respSuc{obj=[]microservice.Subscription}
// @Failure      400  {object}  respErr
// @Failure      404  {object}  respErr
// @Failure      500  {object}  respErr
// @Router       /subscription/query	 [get]
func (a *SubscriptionHandler) querySubscriptions(c *gin.Context) {
	const op = "handler.querySubscriptions"
	log, ctx := prepareTools(c, op)

	args := &service.SubscriptionQueryArgs{}
	if err := c.ShouldBindJSON(args); err != nil {
		if !errors.Is(err, io.EOF) {
			log.Debug().Err(err).Msg("error binding json")
			writeBadRequest(c, "error binding json: "+err.Error())
			return
		}
		log.Debug().Msg("empty body")
	}

	subs, err := a.sub.Query(ctx, args)
	if err != nil {
		if errors.Is(err, service.ErrNoSuchSubscription) {
			log.Debug().Err(err).Msg("no subscriptions founds")
			writeNotFound(c, "no subscriptions founds")
			return
		}
		log.Error().Err(err).Msg("error getting subscriptions")
		writeServerInternal(c, "error getting subscriptions")
		return
	}
	if len(subs) == 0 {
		log.Debug().Msg("no subscriptions founds")
		writeNotFound(c, "no subscriptions founds")
		return
	}

	writeObj(c, subs)
}

// sumSubscriptions godoc
// @Summary      Sum Subscriptions price
// @Description  Sum Subscriptions by a query
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        query  body    service.SubscriptionQueryArgs  true  "query arguments"
// @Success      200  {object}  resp{obj=int}
// @Failure      400  {object}  respErr
// @Failure      404  {object}  respErr
// @Failure      500  {object}  respErr
// @Router       /subscription/sum	 [get]
func (a *SubscriptionHandler) sumSubscriptions(c *gin.Context) {
	const op = "handler.getSubscriptionsSum"
	log, ctx := prepareTools(c, op)

	args := &service.SubscriptionQueryArgs{}
	if err := c.ShouldBindJSON(args); err != nil {
		if !errors.Is(err, io.EOF) {
			log.Debug().Err(err).Msg("error binding json")
			writeBadRequest(c, "error binding json: "+err.Error())
			return
		}
		log.Debug().Msg("empty body")
	}

	sum, err := a.sub.Sum(ctx, args)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoUserID) || errors.Is(err, service.ErrNoSubscriptionID):
			log.Debug().Err(err)
			writeBadRequest(c, err.Error())
			return
		case errors.Is(err, service.ErrNoSuchSubscription):
			log.Debug().Err(err).Msg("no subscriptions founds")
			writeSuccess(c, http.StatusNotFound, "not found subscriptions for given args", 0)
			// writeNotFound(c, "no subscriptions founds")
			return
		default:
			log.Error().Err(err).Msg("error getting subscriptions sum")
			writeServerInternal(c, "error getting subscriptions sum")
			return
		}
	}

	log.Info().Interface("sum", int(sum)).Msg("subscriptions sum")

	writeObj(c, sum)
}

func (a *SubscriptionHandler) parseSubscriptionID(c *gin.Context) (int64, error) {
	// Parse ID
	idStr := c.Param("id")
	if idStr == "" {
		return 0, errors.New("can't parse subscription id: empty")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, e.Wrap(fmt.Sprintf("can't parse subscription id by string %s", idStr), err)
	}

	return id, nil
}
