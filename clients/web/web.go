package web

import (
	"Sudoku-Solver/clients/web/views"
	"Sudoku-Solver/internals"
	"Sudoku-Solver/pkg/logger"
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var Log = logger.Get()

type ClientWebOpts struct {
	Host    string
	Port    int
	Timeout time.Duration
}

type ClientWeb struct {
	internals.Client
	Router  *gin.Engine
	Host    string
	Port    int
	Timeout time.Duration
}

type WebResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
	Status  int    `json:"status"`
	Success bool   `json:"success"`
}

func NewClientWeb(client internals.Client, clientOptions ClientWebOpts) *ClientWeb {

	router := gin.Default()

	return &ClientWeb{
		Client:  client,
		Router:  router,
		Host:    clientOptions.Host,
		Port:    clientOptions.Port,
		Timeout: clientOptions.Timeout,
	}
}

func (c *ClientWeb) StartClient() {

	c.Router.GET("/", c.handleIndex())
	c.Router.Use(c.handleError)
	c.Router.StaticFS("/static", http.Dir("clients/web/static"))

	Log.Info().Str("host", c.Host).Int("port", c.Port).Msg("Starting web client")

	err := c.Router.Run(fmt.Sprintf("%s:%d", c.Host, c.Port))

	if err != nil {
		Log.Fatal().Err(err).Msg("Failed to start web client")
	}
}

func (c *ClientWeb) handleIndex() gin.HandlerFunc {
	return func(apiContext *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout))
		defer cancel()

		err := c.render(apiContext, http.StatusOK, views.Layout(views.Index()))

		if err != nil {
			err = apiContext.Error(err)

			if err != nil {
				Log.Error().Err(err).Msg("Failed to render error")
				panic(err)
			}
		}
	}
}

func (c *ClientWeb) handleError(ctx *gin.Context) {
	ctx.Next()

	if len(ctx.Errors) > 0 {
		ctx.JSON(http.StatusInternalServerError, c.buildResponse(ctx.Errors[0].Error(), nil, http.StatusInternalServerError, false))
	}
}

func (c *ClientWeb) render(ctx *gin.Context, status int, template templ.Component) error {
	ctx.Status(status)
	return template.Render(ctx.Request.Context(), ctx.Writer)
}

func (c *ClientWeb) buildResponse(message string, data interface{}, status int, success bool) WebResponse[interface{}] {
	return WebResponse[interface{}]{
		Message: message,
		Data:    data,
		Status:  status,
		Success: success,
	}
}
