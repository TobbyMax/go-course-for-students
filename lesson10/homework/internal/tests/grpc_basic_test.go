package tests

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/app"
	grpcPort "homework10/internal/ports/grpc"
)

func TestGRRPCCreateUser(t *testing.T) {
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

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Name)
	assert.Equal(t, "ivanov@yandex.ru", res.Email)
}

func TestGRRPCGetUser(t *testing.T) {
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

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	assert.NoError(t, err)

	res, err := client.GetUser(ctx, &grpcPort.GetUserRequest{Id: &user.Id})
	assert.NoError(t, err)

	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Oleg", res.Name)
	assert.Equal(t, "ivanov@yandex.ru", res.Email)
}

func TestGRRPCDeleteUser(t *testing.T) {
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

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	assert.NoError(t, err)

	_, err = client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{Id: &user.Id})

	assert.NoError(t, err)

	_, err = client.GetUser(ctx, &grpcPort.GetUserRequest{Id: &user.Id})
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound.Error(), err.Error())
}

func TestGRRPCUpdateUser(t *testing.T) {
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

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	assert.NoError(t, err)

	res, err := client.UpdateUser(ctx, &grpcPort.UpdateUserRequest{Id: &user.Id, Name: "Kanye", Email: "graduation@west.com"})
	assert.NoError(t, err)
	assert.Equal(t, "Kanye", res.Name)
	assert.Equal(t, "graduation@west.com", res.Email)
}

func TestGRRPCCreateAd(t *testing.T) {
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

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	assert.NoError(t, err)

	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: &user.Id})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Forest", res.Title)
	assert.Equal(t, "Hill Drive", res.Text)
	assert.Equal(t, int64(0), res.AuthorId)
	assert.Equal(t, false, res.Published)
}

func TestGRRPCChangeAdStatus(t *testing.T) {
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

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	assert.NoError(t, err)

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: &user.Id})
	assert.NoError(t, err)

	res, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &ad.Id, UserId: &user.Id, Published: true})
	assert.NoError(t, err)

	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Forest", res.Title)
	assert.Equal(t, "Hill Drive", res.Text)
	assert.Equal(t, int64(0), res.AuthorId)
	assert.Equal(t, true, res.Published)
}

func TestGRRPCUpdateAd(t *testing.T) {
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

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	assert.NoError(t, err)

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: &user.Id})
	assert.NoError(t, err)

	res, err := client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{AdId: &ad.Id, UserId: &user.Id, Title: "Corny", Text: "Low Key"})
	assert.NoError(t, err)

	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Corny", res.Title)
	assert.Equal(t, "Low Key", res.Text)
	assert.Equal(t, int64(0), res.AuthorId)
	assert.Equal(t, false, res.Published)
}

func TestGRRPCGetAd(t *testing.T) {
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

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	assert.NoError(t, err)

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: &user.Id})
	assert.NoError(t, err)

	_, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{AdId: &ad.Id, UserId: &user.Id, Title: "Corny", Text: "Low Key"})
	assert.NoError(t, err)
	res, err := client.GetAd(ctx, &grpcPort.GetAdRequest{AdId: &ad.Id})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Corny", res.Title)
	assert.Equal(t, "Low Key", res.Text)
	assert.Equal(t, int64(0), res.AuthorId)
	assert.Equal(t, false, res.Published)
}

func TestGRRPCDeleteAd(t *testing.T) {
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

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	user, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	assert.NoError(t, err)

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: &user.Id})
	assert.NoError(t, err)

	_, err = client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AdId: &ad.Id, AuthorId: &user.Id})
	assert.NoError(t, err)

	_, err = client.GetAd(ctx, &grpcPort.GetAdRequest{AdId: &ad.Id})
	assert.Error(t, err)

	assert.Equal(t, ErrAdNotFound.Error(), err.Error())
}
