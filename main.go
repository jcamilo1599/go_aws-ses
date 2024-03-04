package main

import (
  "context"
  "fmt"

  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/config"
  "github.com/aws/aws-sdk-go-v2/credentials"
  "github.com/aws/aws-sdk-go-v2/service/ses"
  "github.com/aws/aws-sdk-go-v2/service/ses/types"
)

func main() {
  // Configuración manual de las credenciales y la región
  customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
    if service == ses.ServiceID {
      return aws.Endpoint{
        PartitionID:   "aws",
        URL:           "https://email." + region + ".amazonaws.com",
        SigningRegion: region,
      }, nil
    }
    // Fallback to the default resolver
    return aws.Endpoint{}, &aws.EndpointNotFoundError{}
  })

  cfg, err := config.LoadDefaultConfig(context.TODO(),
    config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("KEY", "SECRET-KEY", "")),
    config.WithRegion("us-east-2"),
    config.WithEndpointResolverWithOptions(customResolver),
  )
  if err != nil {
    fmt.Println("Error al cargar la configuración de AWS:", err)
    return
  }

  // Crear un cliente SES
  client := ses.NewFromConfig(cfg)

  // Definir los parámetros del correo electrónico
  input := &ses.SendEmailInput{
    Destination: &types.Destination{
      ToAddresses: []string{
        "example@email.com", // Destinatario(s)
      },
    },
    Message: &types.Message{
      Body: &types.Body{
        Text: &types.Content{
          Data: aws.String("Hola, este es el cuerpo del correo."),
        },
      },
      Subject: &types.Content{
        Data: aws.String("Asunto del correo"),
      },
    },
    Source: aws.String("example@email.com"), // Remitente (Identidad de SES)
  }

  // Enviar el correo electrónico
  result, err := client.SendEmail(context.TODO(), input)
  if err != nil {
    fmt.Printf("Error al enviar el correo electrónico: %s\n", err)
    return
  }

  fmt.Printf("Correo enviado, ID: %s\n", *result.MessageId)
}
