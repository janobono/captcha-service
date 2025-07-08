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

```bash
docker build -t captcha-service:latest .
```

## run

```bash
docker compose up
```

## stop

```bash
docker compose down
```

## make

If you have `make`, `Go` an `node.js` installed, you can use these prepared targets:

- `tools` - to install all tools and modules
- `clean` - to delete all generated sources
- `generate-openapi` - generate source files from openapi (http) into `generated/openapi`
- `generate-proto` - generate source files from proto (gRPC) into `generated/proto`
- `generate` - to generate all sources
- `build` - default target will call generate and build everything
- `fmt` - format code
- `test` - run tests
- `vet` - check code with vet

## local run

To start the service locally, follow these steps:

1. **Set Environment Variables**  
   Create a `.env.local` file with the required environment variables, or set them manually in your system.

2. **Set `CAPTCHA_FONT`**  
   The `CAPTCHA_FONT` variable must point to a valid font file on your system. Here are typical paths by OS:

    - **Linux**: `/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf`
    - **macOS**: `/Library/Fonts/Arial.ttf`
    - **Windows**: `C:\\Windows\\Fonts\\arialbd.ttf`

> 💡 **Tip**: You can copy the content of the `captcha-service.env` file and simply update the `CAPTCHA_FONT` value as
> needed.

## environment variables

| Name                         | Example                                              | Description                                                     |
|:-----------------------------|:-----------------------------------------------------|:----------------------------------------------------------------|
| PROD                         | false                                                | Production mode flag - log level is switched from debug to info |
| GRPC_ADDRESS                 | :50052                                               | Service gRPC port                                               |
| HTTP_ADDRESS                 | :8080                                                | Service HTTP port                                               |
| CONTEXT_PATH                 | /api                                                 | REST API context path                                           |
|                              |                                                      |                                                                 |
| CAPTCHA_CHARACTERS           | abcdefghijklmnopqrstuvwxyz0123456789                 | Characters used in captchas                                     |
| CAPTCHA_TEXT_LENGTH          | 8                                                    | Captcha text length                                             |
| CAPTCHA_IMAGE_WIDTH          | 200                                                  | Captcha image width                                             |
| CAPTCHA_IMAGE_HEIGHT         | 70                                                   | Captcha image height                                            |
| CAPTCHA_NOISE_LINES          | 8                                                    | Number of noise lines in the captcha                            |
| CAPTCHA_FONT                 | /usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf | Font path used for captcha generation                           |
| CAPTCHA_FONT_SIZE            | 32                                                   | Font size for captcha                                           |
| CAPTCHA_TOKEN_ISSUER         | captcha                                              | Token issuer for captcha                                        |
| CAPTCHA_TOKEN_EXPIRES_IN     | 30                                                   | Captcha token expiration time in minutes                        |
| CAPTCHA_TOKEN_JWK_EXPIRES_IN | 720                                                  | Captcha token JWK expiration time in minutes                    |
|                              |                                                      |                                                                 |
| CORS_ALLOWED_ORIGINS         | http://localhost:5173,http://localhost:3000          | Allowed origins for CORS                                        |
| CORS_ALLOWED_METHODS         | GET,POST,PUT,PATCH,DELETE                            | Allowed HTTP methods                                            |
| CORS_ALLOWED_HEADERS         | Origin,Content-Type,Accept,Authorization             | Allowed HTTP headers                                            |
| CORS_EXPOSED_HEADERS         | Content-length                                       | Exposed headers in CORS                                         |
| CORS_ALLOW_CREDENTIALS       | true                                                 | Whether credentials are allowed in CORS                         |
| CORS_MAX_AGE                 | 12                                                   | Max age (in hours) for CORS preflight response caching          |
