package pb

import (
	"context"
	"fmt"
	"github.com/lookeme/short-url/internal/app/domain/shorten"
	"github.com/lookeme/short-url/internal/models"
	"github.com/lookeme/short-url/internal/security"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IShortenService struct {
	URLService shorten.URLService
}

func (s *IShortenService) mustEmbedUnimplementedShortenURLServiceServer() {
	//TODO implement me
	panic("implement me")
}

func NewShortenService(urlService shorten.URLService) *IShortenService {
	return &IShortenService{URLService: urlService}
}

func (s *IShortenService) CreateShortURL(ctx context.Context, in *CreateAndSaveRequest) (*CreateAndSaveResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	tokenArr := md["authorization"]
	userID := security.GetUserID(tokenArr[0])
	result, err := s.URLService.CreateAndSave(in.Url, userID)
	if err != nil {
		return nil, err
	}
	return &CreateAndSaveResponse{
		Result: result,
	}, nil
}
func (s *IShortenService) GetByID(ctx context.Context, in *FindByURLRequest) (*FindByURLResponse, error) {
	val, ok := s.URLService.FindByKey(in.Key)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "value is not found")
	}
	if val.DeletedFlag {
		return nil, status.Errorf(codes.NotFound, "value deleted")
	}
	return &FindByURLResponse{
		Data: &ShortenData{
			ID:            val.ID,
			CorrelationID: val.CorrelationID,
			ShortURL:      val.ShortURL,
			UserID:        int64(val.UserID),
			DeletedFlag:   val.DeletedFlag,
		},
	}, nil
}

func (s *IShortenService) CreateShortBatch(ctx context.Context, in *CreateAndSaveBatchRequest) (*CreateAndSaveBatchResponse, error) {
	batch := make([]models.BatchRequest, len(in.Req))
	for _, request := range in.Req {
		batch = append(batch, models.BatchRequest{
			CorrelationID: request.CorrelationID,
			OriginalURL:   request.OriginalURL,
		})
	}
	saveBatch, err := s.URLService.CreateAndSaveBatch(batch)
	if err != nil {
		return nil, err
	}
	batchResp := make([]*BatchResponse, len(in.Req))
	for _, response := range saveBatch {
		batchResp = append(batchResp, &BatchResponse{
			CorrelationID: response.CorrelationID,
			ShortURL:      response.ShortURL,
		})
	}
	return &CreateAndSaveBatchResponse{
		Response: batchResp,
	}, nil
}
func (s *IShortenService) GetAllUserURLs(ctx context.Context, _ *emptypb.Empty) (*FindAllResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	tokenArr := md["authorization"]
	userID := security.GetUserID(tokenArr[0])
	shortenData, err := s.URLService.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}
	resp := make([]*ShortenData, len(shortenData))
	for _, val := range shortenData {
		resp = append(resp, &ShortenData{
			ID:            val.ID,
			CorrelationID: val.CorrelationID,
			ShortURL:      val.ShortURL,
			UserID:        int64(val.UserID),
			DeletedFlag:   val.DeletedFlag,
		})
	}
	return &FindAllResponse{
		Data: resp,
	}, nil

}
func (s *IShortenService) Ping(ctx context.Context, _ *emptypb.Empty) (*ResponseStatus, error) {
	err := s.URLService.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &ResponseStatus{
		Status: "200",
	}, nil
}
func (s *IShortenService) DeleteByShortURLs(ctx context.Context, in *DeleteByShortURLSRequest) (*ResponseStatus, error) {

	err := s.URLService.DeleteByShortURLs(in.Urls)
	if err != nil {
		return nil, err
	}
	return &ResponseStatus{
		Status: "202",
	}, nil
}
