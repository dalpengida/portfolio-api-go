package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
)

var (
	fiberLambda *fiberadapter.FiberLambda
	app         *fiber.App

	LOCAL_PORT string = ":5678"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return fiberLambda.ProxyWithContextV2(ctx, req)
}

func init() {
	app = fiber.New(
		fiber.Config{
			UnescapePath: true,
		},
	)
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON("pong")
	})

	if os.Getenv("AWS_EXECUTION_ENV") == "AWS_Lambda_go1.x" {
		fiberLambda = fiberadapter.New(app)
		lambda.Start(handler)
	} else {
		app.Listen(LOCAL_PORT)
	}

}

func main() {}
