package grpcAPI

import (
	"context"
	"errors"
	collectorV1 "github.com/Skvidich/CollectorProto/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ControlAPI struct {
	collectorV1.UnimplementedControlServer
	cord Coordinator
}

func NewControlAPI(cord Coordinator) *ControlAPI {
	return &ControlAPI{
		UnimplementedControlServer: collectorV1.UnimplementedControlServer{},
		cord:                       cord,
	}
}

func (c *ControlAPI) GetList(ctx context.Context, request *collectorV1.ListRequest) (*collectorV1.ListResponse, error) {

	select {
	case <-ctx.Done():
		err := ctx.Err()
		switch {
		case errors.Is(err, context.Canceled):
			return nil, status.Error(codes.Canceled, "context canceled")
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.Canceled, "context canceled")
		default:
			return nil, err
		}

	default:
		list := c.cord.GetterList()

		resList := make([]*collectorV1.Getter, len(list))

		for i := range list {
			resList[i] = &collectorV1.Getter{Name: list[i].Name, Status: list[i].State}
		}
		resp := &collectorV1.ListResponse{List: resList}
		return resp, nil
	}

}

func (c *ControlAPI) GetInfo(ctx context.Context, request *collectorV1.InfoRequest) (*collectorV1.InfoResponse, error) {
	select {
	case <-ctx.Done():
		err := ctx.Err()
		switch {
		case errors.Is(err, context.Canceled):
			return nil, status.Error(codes.Canceled, "context canceled")
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.Canceled, "context canceled")
		default:
			return nil, err
		}

	default:
		info, err := c.cord.Getter(request.GetName())
		if err != nil {
			return nil, status.Error(codes.NotFound, "such getter don't exist")
		}

		return &collectorV1.InfoResponse{Status: info.State}, nil
	}
}

func (c *ControlAPI) Start(ctx context.Context, request *collectorV1.StartRequest) (*collectorV1.StartResponse, error) {
	select {
	case <-ctx.Done():
		err := ctx.Err()
		switch {
		case errors.Is(err, context.Canceled):
			return nil, status.Error(codes.Canceled, "context canceled")
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.Canceled, "context canceled")
		default:
			return nil, err
		}

	default:
		err := c.cord.Start(request.GetName())
		if err != nil {
			return &collectorV1.StartResponse{Result: false}, status.Error(codes.NotFound, err.Error())
		}
		return &collectorV1.StartResponse{Result: true}, nil
	}
}

func (c *ControlAPI) Stop(ctx context.Context, request *collectorV1.StopRequest) (*collectorV1.StopResponse, error) {
	select {
	case <-ctx.Done():
		err := ctx.Err()
		switch {
		case errors.Is(err, context.Canceled):
			return nil, status.Error(codes.Canceled, "context canceled")
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.Canceled, "context canceled")
		default:
			return nil, err
		}

	default:
		err := c.cord.Stop(request.GetName())
		if err != nil {
			return &collectorV1.StopResponse{Result: false}, status.Error(codes.NotFound, err.Error())
		}
		return &collectorV1.StopResponse{Result: true}, nil
	}
}

func (c *ControlAPI) Shutdown(ctx context.Context, request *collectorV1.ShutdownRequest) (*collectorV1.ShutdownResponse, error) {
	select {
	case <-ctx.Done():
		err := ctx.Err()
		switch {
		case errors.Is(err, context.Canceled):
			return nil, status.Error(codes.Canceled, "context canceled")
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.Canceled, "context canceled")
		default:
			return nil, err
		}

	default:
		c.cord.Shutdown()
		return &collectorV1.ShutdownResponse{Result: true}, nil
	}
}
