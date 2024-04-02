package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var (
	ginLambda *ginadapter.GinLambdaV2

	LOCAL_PORT string = ":5678"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func init() {
	app := gin.Default()

	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if os.Getenv("AWS_EXECUTION_ENV") == "AWS_Lambda_go1.x" {
		ginLambda = ginadapter.NewV2(app)
		lambda.Start(handler)
	} else {
		app.Run(LOCAL_PORT)
	}
}

func main() {}
