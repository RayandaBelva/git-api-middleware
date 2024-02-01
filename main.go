package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Account struct {
	AccountID   int    `json:"accountid" binding:"required"`
	AccountName string `json:"accountname" binding:"required"`
	Membership  string `json:"membership" binding:"required"`
	Status      string `json:"status" binding:"required"`
	Durasi      string `json:"durasi" binding:"required"`
}

var data = []Account{
	{1, "Rayanda", "Platinum", "Active", "1 Month"},
	{2, "Redo", "Gold", "Active", "6 Month"},
	{3, "Ismail", "Bronze", "Active", "12 Month"},
	{4, "Belva", "Non Member", "Expired", "0 Month"},
}

type AuthHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}

type LogRequest struct {
	Latency    time.Duration
	StatusCode int
	ClientIP   string
	Method     string
	Path       string
	Message    string
}

func main() {
	router := gin.Default()

	router.Use(AuthMiddleware())
	router.Use(Logger())
	accountRouter := router.Group("/account")
	accountRouter.POST("/registration", CreateAccount)
	accountRouter.GET("/id/:id", getAccountByID)
	accountRouter.GET("/name/:accountname", getAccountByAccountName)
	accountRouter.GET("/member/:membership", getAccountByMembership)
	accountRouter.GET("/status/:status", getAccountByStatus)
	accountRouter.GET("/", getAllAccount)
	accountRouter.PUT("/:id", updateAccountByID)
	accountRouter.DELETE("/:id", deleteAccountByID)

	err := router.Run()
	if err != nil {
		panic(err)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/account/create" {
			c.Next()
		} else {
			h := AuthHeader{}
			if err := c.ShouldBindHeader(&h); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				c.Abort()
			}

			if h.AuthorizationHeader == "Token Aktivasi" {
				c.Next()
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "token invalid",
				})
				c.Abort()
			}
		}
	}
}

func Logger() gin.HandlerFunc {
	// Open log file
	logFile, err := os.OpenFile("E:/membership-streaming/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// create new logger instance
	logger := log.New(logFile, "", 0)

	return func(c *gin.Context) {
		// request start time
		startTime := time.Now()

		// process request
		c.Next()

		// request end time
		endTime := time.Now()

		// log request details
		logRequest := LogRequest{
			Latency:    endTime.Sub(startTime),
			StatusCode: c.Writer.Status(),
			ClientIP:   c.ClientIP(),
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			Message:    c.Errors.ByType(gin.ErrorTypePrivate).String(),
		}

		logString := "\n" +
			"[GIN] " + endTime.Format("2006/01/02 - 15:04:05") + " " +
			strconv.Itoa(logRequest.StatusCode) + " " +
			logRequest.Latency.String() + " " +
			logRequest.ClientIP + " " +
			logRequest.Method + " " +
			logRequest.Path + " " +
			logRequest.Message
		logger.Println(logString)

		_, err := gin.DefaultWriter.Write([]byte(logString))

		if err != nil {
			fmt.Println(err)
		}
	}
}

func CreateAccount(c *gin.Context) {
	var dataAccount Account

	if err := c.ShouldBind(&dataAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create",
		})
		return
	}

	data = append(data, dataAccount)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Account Has Been Created",
		"Token Code": "Token Aktivasi",
	})
}

func getAccountByID(c *gin.Context) {
	accountID := c.Param("id")

	for _, acc := range data {
		if strconv.Itoa(acc.AccountID) == accountID {
			c.JSON(http.StatusOK, gin.H{
				"data": acc,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Account Id Not Found",
	})
}

func getAccountByAccountName(c *gin.Context) {
	accountName := c.Param("accountname")

	var matchingAccounts []Account

	for _, acc := range data {
		if accountName == acc.AccountName {
			matchingAccounts = append(matchingAccounts, acc)
		}
	}

	if len(matchingAccounts) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": matchingAccounts,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No accounts found",
		})
	}
}

func getAccountByMembership(c *gin.Context) {
	memberShip := c.Param("membership")

	var matchingAccounts []Account

	for _, acc := range data {
		if memberShip == acc.Membership {
			matchingAccounts = append(matchingAccounts, acc)
		}
	}

	if len(matchingAccounts) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": matchingAccounts,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No accounts found with the specified membership",
		})
	}
}

func getAccountByStatus(c *gin.Context) {
	status := c.Param("status")

	var matchingAccounts []Account

	for _, acc := range data {
		if status == acc.Status {
			matchingAccounts = append(matchingAccounts, acc)
		}
	}

	if len(matchingAccounts) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": matchingAccounts,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No accounts found with the specified status",
		})
	}
}

func getAllAccount(c *gin.Context) {
	if len(data) == 0 {
		data = []Account{}
	}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func updateAccountByID(c *gin.Context) {
	accountIDs := c.Param("id")
	accountID, err := strconv.Atoi(accountIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	var user Account
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	for i, acc := range data {
		if acc.AccountID == accountID {
			if user.AccountName != "" {
				data[i].AccountName = user.AccountName
			}

			if user.Membership != "" {
				data[i].Membership = user.Membership
			}

			if user.Durasi != "" {
				data[i].Durasi = user.Durasi
			}

			if user.Status != "" {
				data[i].Status = user.Status
			}
			c.JSON(http.StatusOK, gin.H{"message": "Account Updated"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Account Not Found"})
}

func deleteAccountByID(c *gin.Context) {
	accountID := c.Param("id")

	for i, acc := range data {
		if strconv.Itoa(acc.AccountID) == accountID {
			data = append(data[:i], data[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "Account deleted successfully",
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"message": "Account not found",
	})
}
