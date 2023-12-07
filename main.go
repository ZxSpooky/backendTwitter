package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/ZxSpooky/backendTwitter/awsgo"
	"github.com/ZxSpooky/backendTwitter/bd"
	"github.com/ZxSpooky/backendTwitter/handlers"
	"github.com/ZxSpooky/backendTwitter/models"
	"github.com/ZxSpooky/backendTwitter/secretmanager"
	
)

func main() {
	fmt.Println("Iniciando TwittorAWS")

	lambda.Start(EjecutoLambda)

}

/*EjecutoLambda funcion que permite realizar el llamado a la api de lamda*/
func EjecutoLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse

	awsgo.InicializoAWS()

	if !ValidoParametros() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400, //Error por falta de paramtro
			Body:       "Error en las variables de entorno. deben incluir 'SecretName','BucketName','UrlPrefix' ",
			Headers: map[string]string{
				"Content-Type": "appication/json",
			},
		}

		return res, nil
	}

	SecretModel, err := secretmanager.GetSecret(os.Getenv("SecretName"))
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400, //Error por falta de paramtro
			Body:       "Error en la lectura de Secret " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "appication/json",
			},
		}
		return res, nil
	}

	path := strings.Replace(request.PathParameters["twittergo"], os.Getenv("UrlPrefix"), "", -1)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), SecretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), SecretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), SecretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), SecretModel.Database)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("jwtSign"), SecretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketname"), os.Getenv("BucketName"))

	//Chequeo conexion a la base de datos
	err = bd.ConectarBD(awsgo.Ctx)
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 500, //Error por falta de paramtro
			Body:       "Error conectado en la base de datos " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "appication/json",
			},
		}
		return res, nil
	}

	respAPI := handlers.Manejadores(awsgo.Ctx, request)
	if respAPI.CustomResp == nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: respAPI.Status, //Error por falta de paramtro
			Body:       respAPI.Message,
			Headers: map[string]string{
				"Content-Type": "appication/json",
			},
		}
		return res, nil
	} else {
		return respAPI.CustomResp, nil
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

	//lcano

}