package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/app"
	grpcPort "homework9/internal/ports/grpc"
	"net"
	"testing"
	"time"
)

func TestGRPCChangeStatusAdOfAnotherUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user1, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	assert.NoError(t, err)

	user2, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	assert.NoError(t, err)

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "GOMD", Text: "Role Modelz", UserId: user1.Id})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: ad.Id, UserId: user2.Id, Published: true})
	assert.Error(t, err)

	assert.Equal(t, ErrGRPCForbidden.Error(), err.Error())
}

func TestGRPCUpdateAdOfAnotherUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	_, err = client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	assert.NoError(t, err)

	_, err = client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Pimp", Text: "A Butterfly", UserId: 1})
	assert.NoError(t, err)

	_, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{AdId: 0, UserId: 0, Title: "Mr. Morale", Text: "The Big Steppers"})
	assert.Error(t, err)

	assert.Equal(t, ErrGRPCForbidden.Error(), err.Error())
}

func TestGRPCCreateAd_ID(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	_, err = client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	assert.NoError(t, err)

	_, err = client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	assert.NoError(t, err)

	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Pimp", Text: "A Butterfly", UserId: 1})
	assert.NoError(t, err)
	assert.Equal(t, res.Id, int64(0))

	res, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 1, Title: "Mr. Morale", Text: "The Big Steppers"})
	assert.NoError(t, err)
	assert.Equal(t, res.Id, int64(1))

	res, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 0, Title: "Cole World", Text: "Born Sinner"})
	assert.NoError(t, err)
	assert.Equal(t, res.Id, int64(2))
}

func TestGRPCDeleteAdOfAnotherUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	_, err = client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	assert.NoError(t, err)

	_, err = client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: 0})
	assert.NoError(t, err)

	_, err = client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AdId: 0, AuthorId: 1})
	assert.Error(t, err)

	assert.Equal(t, ErrGRPCForbidden.Error(), err.Error())
}

func TestGRPCGetUser_NonExistentID(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	_, err = client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "MacMiller", Email: "blue_slide_park@hotmail.com"})
	assert.NoError(t, err)

	_, err = client.GetUser(ctx, &grpcPort.GetUserRequest{Id: 1})
	assert.Error(t, err)

	assert.Equal(t, ErrUserNotFound.Error(), err.Error())
}
