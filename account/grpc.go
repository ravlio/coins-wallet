package account

import (
	"context"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/ravlio/wallet/account/pb"
	"github.com/ravlio/wallet/pkg/grpcutil"
	"google.golang.org/grpc"
)

// some mappers. Without errors, for simplicity
func accountFromProtobuf(req *pb.Account) *Account {
	ret := &Account{}
	ret.fromProtobuf(req)

	return ret
}

func (a *Account) fromProtobuf(req *pb.Account) {
	a.ID = req.ID
	a.Name = req.Name
	a.Email = req.Email
	a.CreatedAt = req.CreatedAt
	a.UpdatedAt = req.UpdatedAt
}

func (a *Account) toProtobuf() *pb.Account {
	ret := &pb.Account{
		ID:        a.ID,
		Name:      a.Name,
		Email:     a.Email,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}

	return ret
}

type grpcServer struct {
	cl Client
}

var _ pb.AccountServiceServer = &grpcServer{}

func (g *grpcServer) CreateAccount(ctx context.Context, req *pb.Account) (resp *pb.Account, err error) {
	defer grpcutil.LogRequest("account", "CreateAccount", req).LogResponse(ctx, resp, err)

	clresp, err := g.cl.CreateAccount(ctx, accountFromProtobuf(req))
	if err != nil {
		return nil, err
	}

	return clresp.toProtobuf(), nil
}

func (g *grpcServer) GetAccount(ctx context.Context, id *wrappers.UInt32Value) (*pb.Account, error) {
	resp, err := g.cl.GetAccount(ctx, id.Value)
	if err != nil {
		return nil, err
	}
	return resp.toProtobuf(), nil
}

func (g *grpcServer) DeleteAccount(ctx context.Context, req *wrappers.UInt32Value) (*empty.Empty, error) {
	err := g.cl.DeleteAccount(ctx, req.Value)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (g *grpcServer) UpdateAccount(ctx context.Context, req *pb.Account) (*pb.Account, error) {
	resp, err := g.cl.UpdateAccount(ctx, accountFromProtobuf(req))
	if err != nil {
		return nil, err
	}

	return resp.toProtobuf(), nil
}

func (g *grpcServer) ListAccount(_ *empty.Empty, ret pb.AccountService_ListAccountServer) error {
	resp, err := g.cl.ListAccounts(context.Background())
	if err != nil {
		return err
	}

	for _, v := range resp {
		err = ret.Send(v.toProtobuf())
		if err != nil {
			return err
		}
	}

	return nil
}

func NewGRPCServer(svc Client) *grpc.Server {
	s := &grpcServer{cl: svc}
	ret := grpc.NewServer()
	pb.RegisterAccountServiceServer(ret, s)

	return ret
}

type grpcClient struct {
	cl pb.AccountServiceClient
}

var _ Client = &grpcClient{}

func NewClient(cl pb.AccountServiceClient) Client {
	return &grpcClient{cl: cl}
}

func (g *grpcClient) CreateAccount(ctx context.Context, req *Account) (resp *Account, err error) {
	gresp, err := g.cl.CreateAccount(ctx, req.toProtobuf())
	if err != nil {
		return nil, err
	}

	return accountFromProtobuf(gresp), nil
}

func (g *grpcClient) GetAccount(ctx context.Context, id uint32) (*Account, error) {
	resp, err := g.cl.GetAccount(ctx, &wrappers.UInt32Value{Value: id})

	if err != nil {
		return nil, err
	}

	return accountFromProtobuf(resp), nil
}

func (g *grpcClient) DeleteAccount(ctx context.Context, id uint32) error {
	_, err := g.cl.DeleteAccount(ctx, &wrappers.UInt32Value{Value: id})

	return err
}

func (g *grpcClient) UpdateAccount(ctx context.Context, req *Account) (*Account, error) {
	resp, err := g.cl.UpdateAccount(ctx, req.toProtobuf())
	if err != nil {
		return nil, err
	}

	return accountFromProtobuf(resp), nil
}

func (g *grpcClient) ListAccounts(ctx context.Context) ([]*Account, error) {
	stream, err := g.cl.ListAccount(ctx, nil)
	if err != nil {
		return nil, err
	}

	ret := make([]*Account, 0)

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return ret, nil
		}
		if err != nil {
			return nil, err
		}

		ret = append(ret, accountFromProtobuf(in))
	}
}
