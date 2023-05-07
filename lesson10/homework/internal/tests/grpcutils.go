package tests

import (
	"context"
	"errors"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/app"
	grpcPort "homework10/internal/ports/grpc"
	"log"
	"net"
	"time"
)

var (
	ErrUserNotFound    = errors.New("rpc error: code = NotFound desc = user with such id does not exist")
	ErrAdNotFound      = errors.New("rpc error: code = NotFound desc = ad with such id does not exist")
	ErrGRPCForbidden   = errors.New("rpc error: code = PermissionDenied desc = forbidden")
	ErrInvalidEmail    = errors.New("rpc error: code = InvalidArgument desc = mail: missing '@' or angle-addr")
	ErrMissingArgument = errors.New("rpc error: code = InvalidArgument desc = required argument is missing")
	ErrMockInternal    = errors.New("rpc error: code = Internal desc = mock error")
	ErrValidationMock  = errors.New("rpc error: code = InvalidArgument desc = ")
	ErrDateMock        = errors.New("rpc error: code = InvalidArgument desc = parsing time \"20/02/2022\" as \"2006-01-02\": cannot parse \"2/2022\" as \"2006\"")
)

type GRPCSuite struct {
	suite.Suite
	Client  grpcPort.AdServiceClient
	Conn    *grpc.ClientConn
	Context context.Context
	Cancel  context.CancelFunc
	Server  *grpc.Server
	Lis     *bufconn.Listener
}

func (suite *GRPCSuite) SetupTest() {
	log.Println("Setting Up Test")

	suite.Lis = bufconn.Listen(1024 * 1024)
	suite.Server = grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpcPort.UnaryLoggerInterceptor,
		grpcPort.UnaryRecoveryInterceptor(),
	))

	svc := grpcPort.NewService(app.NewApp(adrepo.New()))
	grpcPort.RegisterAdServiceServer(suite.Server, svc)
	go func() {
		suite.NoError(suite.Server.Serve(suite.Lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return suite.Lis.Dial()
	}

	suite.Context, suite.Cancel = context.WithTimeout(context.Background(), 30*time.Second)

	conn, err := grpc.DialContext(suite.Context, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	suite.NoError(err, "grpc.DialContext")
	suite.Conn = conn

	suite.Client = grpcPort.NewAdServiceClient(suite.Conn)
}

func (suite *GRPCSuite) TearDownTest() {
	log.Println("Tearing Down Test")

	err := suite.Conn.Close()
	if err != nil {
		log.Println("Error closing connection")
	}
	suite.Cancel()
	suite.Server.Stop()
	err = suite.Lis.Close()
	if err != nil {
		log.Println("Error closing listener")
	}
}
