# captcha-service

Captcha service written in Golang

![example](./img.png)

- [openapi contract](contract/openapi/captcha-service.yaml)
- [proto contract](contract/proto/captcha-service.proto)
- [License](./LICENSE)

## build

```bash
./build.sh
```

or

```shell
docker build -t captcha-service:latest .
```

## run

```shell
docker compose up
```

## stop

```shell
docker compose down
```

## make

If you have make and golang installed, you can use prepared targets.

- clean - to delete all generated sources
- generate-openapi - generate source files from openapi (http) into `generated/openapi`
- generate-proto - generate source files from proto (gRPC) into `generated/proto`
- generate - to generate all sources
- build - default target will call generate and build everything
- fmt - format code
- test - run tests
- vet - check code with vet

## environment variables

| Name                         | Example                                              | Description                                                     |
|:-----------------------------|:-----------------------------------------------------|:----------------------------------------------------------------|
| PROD                         | false                                                | Production mode flag - log level is switched from debug to info |
| GRPC_ADDRESS                 | :50052                                               | Service gRPC port                                               |
| HTTP_ADDRESS                 | :8080                                                | Service http port                                               |
| CONTEXT_PATH                 | /api                                                 | Rest api context path                                           |
|                              |                                                      |                                                                 |
| CAPTCHA_CHARACTERS           | abcdefghijklmnopqrstuvwxyz0123456789                 | Captcha characters                                              |
| CAPTCHA_TEXT_LENGTH          | 8                                                    | Captcha text length                                             |
| CAPTCHA_IMAGE_WIDTH          | 200                                                  | Captcha image width                                             |
| CAPTCHA_IMAGE_HEIGHT         | 70                                                   | Captcha image height                                            |
| CAPTCHA_NOISE_LINES          | 8                                                    | Captcha noise lines                                             |
| CAPTCHA_FONT                 | /usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf | Captcha font path                                               |
| CAPTCHA_FONT_SIZE            | 32                                                   | Captcha font size                                               |
| CAPTCHA_TOKEN_ISSUER         | captcha                                              | Captcha token issuer                                            |
| CAPTCHA_TOKEN_EXPIRES_IN     | 30                                                   | Captcha token expiration in minutes                             |
| CAPTCHA_TOKEN_JWK_EXPIRES_IN | 720                                                  | Captcha token jwk expiration in minutes                         |
|                              |                                                      |                                                                 |
| CORS_ALLOWED_ORIGINS         | http://localhost:5173,http://localhost:3000          | Allowed origins                                                 |
| CORS_ALLOWED_METHODS         | GET,POST,PUT,PATCH,DELETE                            | Allowed methods                                                 |
| CORS_ALLOWED_HEADERS         | Origin,Content-Type,Accept,Authorization             | Allowed headers                                                 |
| CORS_EXPOSED_HEADERS         | Content-length                                       | Exposed headers                                                 |
| CORS_ALLOW_CREDENTIALS       | true                                                 | Allow credentials                                               |
| CORS_MAX_AGE                 | 12                                                   | Max age in hours                                                |
