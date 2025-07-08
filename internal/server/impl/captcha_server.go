package impl

import (
	"context"
	"github.com/janobono/captcha-service/generated/openapi"
	"github.com/janobono/captcha-service/generated/proto"
	"github.com/janobono/captcha-service/internal/service"
	"github.com/janobono/go-util/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	"log/slog"
)

type captchaServer struct {
	proto.UnimplementedCaptchaServer
	captchaService *service.CaptchaService
}

func NewCaptchaServer(captchaService *service.CaptchaService) proto.CaptchaServer {
	return &captchaServer{captchaService: captchaService}
}

func (c *captchaServer) Create(ctx context.Context, empty *emptypb.Empty) (*proto.CaptchaDetail, error) {
	result, err := c.captchaService.Create(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create CAPTCHA: %v", err)
	}
	return &proto.CaptchaDetail{
		Token: result.CaptchaToken,
		Image: result.CaptchaImage,
	}, nil
}

func (c *captchaServer) Validate(ctx context.Context, data *proto.CaptchaData) (*wrapperspb.BoolValue, error) {
	if data == nil || common.IsBlank(data.Token) || common.IsBlank(data.Text) {
		var token, text string
		if data != nil {
			token = data.Token
			text = data.Text
		} else {
			token = "<nil>"
			text = "<nil>"
		}
		slog.Warn("Invalid CAPTCHA input", "token", token, "text", text)
		return nil, status.Error(codes.InvalidArgument, "captcha token and text must be provided")
	}

	result := c.captchaService.Validate(ctx, &openapi.CaptchaData{
		CaptchaToken: data.Token,
		CaptchaText:  data.Text,
	})
	return &wrapperspb.BoolValue{Value: result.Value}, nil
}
