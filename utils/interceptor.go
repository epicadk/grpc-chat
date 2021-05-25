package utils

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Interceptor struct {
	jwtManager       *JWTManager
	protectedMethods []string
}

func NewInterceptor(jwtManager *JWTManager, roles []string) *Interceptor {
	return &Interceptor{jwtManager: jwtManager, protectedMethods: roles}
}

func Roles() []string {
	const servicePath = "/chat.ChatService/"

	return []string{
		servicePath + "Connect",
		servicePath + "SendChat",
	}
}

func (interceptor *Interceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (interceptor *Interceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Println("--> stream interceptor: ", info.FullMethod)

		err := interceptor.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func (interceptor *Interceptor) authorize(ctx context.Context, rpcMethod string) error {
	methods, isProtected := interceptor.protectedMethods, false
	for _, method := range methods {
		if method == rpcMethod {
			isProtected = true
		}
	}
	if !isProtected {
		// open method
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid")
	}
	md.Append("user", claims.PhoneNumber)
	log.Println(claims)
	return nil
}
