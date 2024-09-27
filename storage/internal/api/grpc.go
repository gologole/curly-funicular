package api

import (
	"cmd/main.go/configs"
	"cmd/main.go/internal/api/rpc"
	"cmd/main.go/internal/service"
	"cmd/main.go/models"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	rpc.UnimplementedStorageServer
	cfgs      *configs.Config
	service   service.Service
	elklogger *zap.Logger
}

func NewGrpcServer(cfgs *configs.Config, service service.Service, elk *zap.Logger) rpc.StorageServer {
	return &GrpcServer{cfgs: cfgs, service: service, elklogger: elk}
}

func (c *GrpcServer) GetFile(r *rpc.FileRequest, s grpc.ServerStreamingServer[rpc.File]) error {
	file, err := c.service.Getfile(r.GetName(), string(r.GetUserid()))
	if err != nil {
		return status.Errorf(codes.Unknown, err.Error())
	}

	s.Send(&rpc.File{
		Name:     file.Name,
		Data:     file.Data,
		Hashfile: file.Hash,
	})
	return nil
}
func (c *GrpcServer) UploadFile(cs grpc.ClientStreamingServer[rpc.PutFileRequest, rpc.Response]) error {
	for {
		req, err := cs.Recv()
		if err == grpc.Errorf(codes.Canceled, "") {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Unknown, err.Error())
		}
		file := models.File{
			Name:   req.GetName(),
			Data:   req.GetData(),
			UserID: int(req.GetUserid()),
			Hash:   req.GetHashfile(),
		}
		err = c.service.UploadFile(file)
		if err != nil {
			return status.Errorf(codes.Unknown, err.Error())
		}
		cs.SendAndClose(&rpc.Response{Err: "", Success: true})
	}
	return nil
}

func (c *GrpcServer) DeleteFile(ctx context.Context, req *rpc.FileRequest) (*rpc.Response, error) {
	err := c.service.DeleteFile(req.Name, string(req.Userid))
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return &rpc.Response{Err: "", Success: true}, nil
}

// передает только имена файлов! нужен реген чтобы не было лишних полей
func (c *GrpcServer) GetFileList(ctx context.Context, req *rpc.FileListRequest) (*rpc.FileListResponse, error) {
	files, err := c.service.GetFilesByUserID(string(req.Userid))
	c.elklogger.Info("run")
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	fileList := make([]*rpc.File, 0, len(files))
	for _, file := range files {
		fileList = append(fileList, &rpc.File{
			Name: file.Name,
		})
	}
	return &rpc.FileListResponse{Files: fileList}, nil
}
