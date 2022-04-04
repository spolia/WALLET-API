package internal

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spolia/wallet-api/internal/wallet/movement"
	"github.com/spolia/wallet-api/internal/wallet/user"
)

//signIn and signOut user

//create user
func createUser(service Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userRequest user.User
		if err := ctx.ShouldBindJSON(&userRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		userID, err := service.CreateUser(ctx, userRequest.FirstName, userRequest.LastName, userRequest.Alias, userRequest.Email)
		if err != nil {
			if err == user.ErrorAlreadyExist {
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusCreated, userID)
	}
}

func getBalance(service Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		alias := ctx.Param("alias")
		if alias == "" {
			ctx.JSON(http.StatusBadRequest, "alias required")
			return
		}

		userBalance, err := service.GetBalance(ctx, alias)
		if err != nil {
			if err == user.ErrorUserNotFound {
				ctx.JSON(http.StatusNotFound, err.Error())
				return
			}

			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, userBalance)
	}
}

func createMovement(service Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var movementRequest movement.Movement
		if err := ctx.ShouldBindJSON(&movementRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		movementID, err := service.CreateMovement(ctx, movementRequest)
		if err != nil {
			if err == movement.ErrorWrongCurrency || err == movement.ErrorWrongUser || err == movement.ErrorInsufficientBalance {
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}

			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, movementID)
	}
}

func searchMovement(service Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := strconv.ParseInt(ctx.Query("userid"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		limit, _ := strconv.ParseUint(ctx.DefaultQuery("limit", "0"), 10, 64)
		offset, _ := strconv.ParseUint(ctx.DefaultQuery("offset", "0"), 10, 0)
		typeMov := ctx.Query("type")
		currency := ctx.Query("currencyname")

		var searchRequest = struct {
			Limit        uint64
			Offset       uint64
			Type         string
			CurrencyName string
			UserID       int64
		}{
			Limit:        limit,
			Offset:       offset,
			Type:         typeMov,
			CurrencyName: currency,
			UserID:       userID,
		}

		movementsResult, err := service.SearchMovement(ctx, searchRequest.UserID, searchRequest.Limit, searchRequest.Offset, searchRequest.Type,
			searchRequest.CurrencyName)
		if err != nil {
			if err == movement.ErrorNoMovements {
				ctx.JSON(http.StatusNotFound, err.Error())
				return
			}

			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, movementsResult)
	}
}
