package main

import (
	"context"
	"os"

	"github.com/ZxSpooky/backendTwitter/awsgo"
	"github.com/ZxSpooky/backendTwitter/models"
	"github.com/ZxSpooky/backendTwitter/secretmanager"
	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(EjecutoLambda)
}

func EjecutoLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse

	awsgo.InicializoAWS()
	if !ValidoParametros() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en las variables de entorno, deben incluir 'SecretName', 'BucketName', 'UrlPrefix'",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}

		SecretModel, err:=secretmanager.GetSecret(os.Getenv("SecretName"))
		if err != nil {
			res = &events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Error en la lectura de Secret"+err.Error(),
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
		}
}
}

func ValidoParametros() bool {
	_, traeParametro := os.LookupEnv("SecretName")
	if !traeParametro {
		return traeParametro
	}

	_, traeParametro = os.LookupEnv("BucketName")
	if !traeParametro {
		return traeParametro
	}

	_, traeParametro = os.LookupEnv("UrlPrefix")
	if !traeParametro {
		return traeParametro
	}
	return traeParametro
}
