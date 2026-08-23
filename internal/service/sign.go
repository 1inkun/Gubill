package service

import (
	"context"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/repository"
)

type SignService struct {
	signQuery *repository.SignQuery
}

func NewSignService(signQuery *repository.SignQuery) *SignService {
	return &SignService{
		signQuery: signQuery,
	}
}

var (
	ErrExistSignData   = models.NewBusinessError(400, "存在未结算的签到订单")
	ErrUserNoExist     = models.NewBusinessError(400, "用户不存在")
	ErrSignDataNoExist = models.NewBusinessError(400, "签到订单不存在")
	ErrWrongUser       = models.NewBusinessError(400, "结算的用户有误")
	ErrAlreadyFinished = models.NewBusinessError(400, "订单已经完成")
	ErrExistPayData    = models.NewBusinessError(400, "存在尚未支付的订单")
)

// 面向用户的接口

func (s *SignService) GenerateNewSignData(ctx context.Context, userId string) (string, error) {
	orderId, err := s.signQuery.GenerateNewSignData(ctx, userId)
	if err != nil {
		return "", err
	}
	return orderId, nil
}

func (s *SignService) GetUserSignData(ctx context.Context, userId string, status int64, page int, pageSize int) ([]models.SignRes, error) {
	results, err := s.signQuery.GetUserSignData(ctx, userId, status, page, pageSize)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *SignService) GetUserSignDataById(ctx context.Context, userId string, signId string) (models.SignRes, error) {
	res, err := s.signQuery.GetSignDataById(ctx, signId)
	if err != nil {
		return models.SignRes{}, err
	}
	if userId != res.UserId {
		return models.SignRes{}, ErrWrongUser
	}
	return res, nil
}

func (s *SignService) FinishSignData(ctx context.Context, userId string, signId string) (string, error) {
	res, err := s.signQuery.FinishSignData(ctx, userId, signId)
	if err != nil {
		return "", err
	}
	return res, nil
}

// 面向管理的接口

func (s *SignService) GetAllSignData(ctx context.Context, page int, pageSize int) ([]models.SignRes, error) {
	res, err := s.signQuery.GetAllSignData(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *SignService) GetSignDataById(ctx context.Context, signId string) (models.SignRes, error) {
	res, err := s.signQuery.GetSignDataById(ctx, signId)
	if err != nil {
		return models.SignRes{}, err
	}
	return res, nil
}

func (s *SignService) UpdateSignData(ctx context.Context, signId string, status int64, value int64) (int, error) {
	res, err := s.signQuery.UpdateSignData(ctx, signId, status, value)
	if err != nil {
		return 0, err
	}
	return res, nil
}
