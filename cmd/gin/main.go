package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/dalpengida/portfolio-api-go-mysql/cmd/gin/route"
	"github.com/dalpengida/portfolio-api-go-mysql/database/my"
	"github.com/gin-gonic/gin"
)

// https://github.com/awslabs/aws-lambda-go-api-proxy/blob/master/gin/adapterv2.go

var (
	ginLambda *ginadapter.GinLambdaV2

	LOCAL_PORT string = ":5678"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func init() {
	err := my.SetBoilerDatabas()
	if err != nil {
		panic(err)
	}

	app := gin.Default()

	route.SetupV1Route(app)

	if os.Getenv("AWS_EXECUTION_ENV") == "AWS_Lambda_go1.x" {
		ginLambda = ginadapter.NewV2(app)
		lambda.Start(handler)
	} else {
		app.Run(LOCAL_PORT)
	}
}

func main() {}
