package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/app"
	grpcPort "homework9/internal/ports/grpc"
	"net"
	"testing"
	"time"
)

func TestGRPCListAds(t *testing.T) {
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
	user1, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	assert.NoError(t, err)

	_, err = client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	assert.NoError(t, err)

	ad1, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	assert.NoError(t, err)

	publishedAd, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &ad1.Id, UserId: &user1.Id, Published: true})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Cole World"})
	assert.NoError(t, err)

	ads, err := client.ListAds(ctx, &grpcPort.ListAdRequest{})

	assert.NoError(t, err)
	assert.Len(t, ads.List, 1)
	assert.Equal(t, ads.List[0].Id, publishedAd.Id)
	assert.Equal(t, ads.List[0].Title, publishedAd.Title)
	assert.Equal(t, ads.List[0].Text, publishedAd.Text)
	assert.Equal(t, ads.List[0].AuthorId, publishedAd.AuthorId)
	assert.True(t, ads.List[0].Published)
}

func TestGRPCListAdsPublished(t *testing.T) {
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
	user1, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	assert.NoError(t, err)

	_, err = client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	assert.NoError(t, err)

	ad1, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	assert.NoError(t, err)

	publishedAd, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &ad1.Id, UserId: &user1.Id, Published: true})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Cole World"})
	assert.NoError(t, err)

	published := true
	ads, err := client.ListAds(ctx, &grpcPort.ListAdRequest{Published: &published})

	assert.NoError(t, err)
	assert.Len(t, ads.List, 1)
	assert.Equal(t, ads.List[0].Id, publishedAd.Id)
	assert.Equal(t, ads.List[0].Title, publishedAd.Title)
	assert.Equal(t, ads.List[0].Text, publishedAd.Text)
	assert.Equal(t, ads.List[0].AuthorId, publishedAd.AuthorId)
	assert.True(t, ads.List[0].Published)
}

func TestGRPCListAdsNotPublished(t *testing.T) {
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
	user1, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	assert.NoError(t, err)

	_, err = client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	assert.NoError(t, err)

	ad1, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &ad1.Id, UserId: &user1.Id, Published: true})
	assert.NoError(t, err)

	notPublishedAd, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Cole World"})
	assert.NoError(t, err)

	published := false
	ads, err := client.ListAds(ctx, &grpcPort.ListAdRequest{Published: &published})

	assert.NoError(t, err)
	assert.Len(t, ads.List, 1)
	assert.Equal(t, ads.List[0].Id, notPublishedAd.Id)
	assert.Equal(t, ads.List[0].Title, notPublishedAd.Title)
	assert.Equal(t, ads.List[0].Text, notPublishedAd.Text)
	assert.Equal(t, ads.List[0].AuthorId, notPublishedAd.AuthorId)
	assert.False(t, ads.List[0].Published)
}

func TestGRPCListAdsByUser(t *testing.T) {
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
	user1, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	assert.NoError(t, err)

	user2, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	assert.NoError(t, err)

	adByUser1, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser1.Id, UserId: &user1.Id, Published: true})
	assert.NoError(t, err)

	adByUser2, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "GOMD", Text: "Cole World"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser2.Id, UserId: &user2.Id, Published: true})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Born Sinner", Text: "Cole World"})
	assert.NoError(t, err)

	ads, err := client.ListAds(ctx, &grpcPort.ListAdRequest{UserId: &user1.Id})

	assert.NoError(t, err)
	assert.Len(t, ads.List, 1)
	assert.Equal(t, ads.List[0].Id, adByUser1.Id)
	assert.Equal(t, ads.List[0].Title, adByUser1.Title)
	assert.Equal(t, ads.List[0].Text, adByUser1.Text)
	assert.Equal(t, ads.List[0].AuthorId, adByUser1.AuthorId)
	assert.True(t, ads.List[0].Published)
}

func TestGRPCListAdsByDate(t *testing.T) {
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
	user1, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	assert.NoError(t, err)

	user2, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	assert.NoError(t, err)

	adByUser1, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser1.Id, UserId: &user1.Id, Published: true})
	assert.NoError(t, err)

	adByUser2, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "GOMD", Text: "Cole World"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser2.Id, UserId: &user2.Id, Published: true})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Born Sinner", Text: "Cole World"})
	assert.NoError(t, err)

	today := time.Now().UTC().Format(DateLayout)
	ads, err := client.ListAds(ctx, &grpcPort.ListAdRequest{Date: &today})
	assert.NoError(t, err)

	assert.Len(t, ads.List, 3)
}

func TestGRPCListAdsByDate_Yesterday(t *testing.T) {
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
	user1, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	assert.NoError(t, err)

	user2, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	assert.NoError(t, err)

	adByUser1, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser1.Id, UserId: &user1.Id, Published: true})
	assert.NoError(t, err)

	adByUser2, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "GOMD", Text: "Cole World"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser2.Id, UserId: &user2.Id, Published: true})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Born Sinner", Text: "Cole World"})
	assert.NoError(t, err)

	yesterday := time.Now().UTC().Add(time.Duration(-24) * time.Hour).Format(DateLayout)
	ads, err := client.ListAds(ctx, &grpcPort.ListAdRequest{Date: &yesterday})
	assert.NoError(t, err)

	assert.Len(t, ads.List, 0)
}

func TestGRPCListAdsByUserAndDate(t *testing.T) {
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
	user1, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	assert.NoError(t, err)

	user2, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	assert.NoError(t, err)

	adByUser1, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser1.Id, UserId: &user1.Id, Published: true})
	assert.NoError(t, err)

	adByUser2, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "GOMD", Text: "Cole World"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser2.Id, UserId: &user2.Id, Published: true})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Born Sinner", Text: "Cole World"})
	assert.NoError(t, err)

	today := time.Now().UTC().Format(DateLayout)
	ads, err := client.ListAds(ctx, &grpcPort.ListAdRequest{Date: &today, UserId: &user2.Id})
	assert.NoError(t, err)

	assert.Len(t, ads.List, 2)
}

func TestGRPCListAdsByTitle(t *testing.T) {
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
	user1, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	assert.NoError(t, err)

	user2, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	assert.NoError(t, err)

	gomd, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &gomd.Id, UserId: &user1.Id, Published: true})
	assert.NoError(t, err)

	adByUser2, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Fire Squad", Text: "Cole World"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser2.Id, UserId: &user2.Id, Published: true})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Born Sinner", Text: "Cole World"})
	assert.NoError(t, err)

	title := "GOMD"
	ads, err := client.ListAds(ctx, &grpcPort.ListAdRequest{Title: &title})
	assert.NoError(t, err)

	assert.Len(t, ads.List, 1)
	assert.Equal(t, ads.List[0].Id, gomd.Id)
	assert.Equal(t, ads.List[0].Title, gomd.Title)
	assert.Equal(t, ads.List[0].Text, gomd.Text)
	assert.Equal(t, ads.List[0].AuthorId, gomd.AuthorId)
	assert.True(t, ads.List[0].Published)
}

func TestGRPCListAdsByTitle_Multiple(t *testing.T) {
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
	user1, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	assert.NoError(t, err)

	user2, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	assert.NoError(t, err)

	adByUser1, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser1.Id, UserId: &user1.Id, Published: true})
	assert.NoError(t, err)

	adByUser2, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "GOMD", Text: "Cole World"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser2.Id, UserId: &user2.Id, Published: true})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Born Sinner", Text: "Cole World"})
	assert.NoError(t, err)

	title := "GOMD"
	ads, err := client.ListAds(ctx, &grpcPort.ListAdRequest{Title: &title})
	assert.NoError(t, err)

	assert.Len(t, ads.List, 2)
}

func TestGRPCListAdsByOptions(t *testing.T) {
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
	user1, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	assert.NoError(t, err)

	user2, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	assert.NoError(t, err)

	adByUser1, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	assert.NoError(t, err)

	target, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser1.Id, UserId: &user1.Id, Published: true})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	assert.NoError(t, err)

	adByUser2, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "GOMD", Text: "Cole World"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser2.Id, UserId: &user2.Id, Published: true})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Born Sinner", Text: "Cole World"})
	assert.NoError(t, err)

	today := time.Now().UTC().Format(DateLayout)
	published := true
	ads, err := client.ListAds(ctx, &grpcPort.ListAdRequest{UserId: &user1.Id, Date: &today, Title: &target.Title, Published: &published})
	assert.NoError(t, err)
	assert.Len(t, ads.List, 1)
	assert.Equal(t, ads.List[0].Id, target.Id)
	assert.Equal(t, ads.List[0].Title, target.Title)
	assert.Equal(t, ads.List[0].Text, target.Text)
	assert.Equal(t, ads.List[0].AuthorId, target.AuthorId)
	assert.True(t, ads.List[0].Published)
}
