package services

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/1inkun/billApp/internal/models"
	"github.com/1inkun/billApp/internal/repositories"
	"golang.org/x/sync/errgroup"
)

type GoodService struct {
	goodQuery *repositories.GoodQuery
	mu        *sync.Mutex
}

func NewGoodService(goodQuery *repositories.GoodQuery) *GoodService {
	var mu sync.Mutex
	return &GoodService{
		goodQuery: goodQuery,
		mu:        &mu,
	}
}

func (s *GoodService) GetGoods(ctx context.Context) ([]models.Good, error) {
	var data []models.Good
	// 在缓存中查找
	s.mu.Lock()
	data, err := s.goodQuery.GetGoodsFromCache(ctx)
	if err != nil {
		// 未找到数据或缓存出错
		data = nil
	}
	s.mu.Unlock()
	// 得到锁以后进行doublecheck
	if data == nil {
		data, err = s.goodQuery.GetGoodsFromCache(ctx)
		if err != nil {
			// 未找到数据或缓存出错
			data = nil
		}
	} else {
		return data, nil
	}
	// 未在缓存中找到
	data, err = s.goodQuery.GetGoods(ctx)
	if err != nil {
		return nil, err
	}
	// 将输入存入缓存
	var eg errgroup.Group
	eg.Go(func() error {
		for _, v := range data {
			dataByte, err := json.Marshal(v)
			if err != nil {
				log.Printf("在GetGoods中序列化json时出错:%s", err)
				return err
			}
			key := "goodsList"
			s.goodQuery.SetCache(ctx, key, string(dataByte))
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *GoodService) BuyGood(ctx context.Context, goodId int64, userId int64) (models.BuyGoodResponse, error) {
	// 验证是否存在商品
	check, err := s.goodQuery.GetGoodsById(ctx, goodId)
	if err != nil {
		return models.BuyGoodResponse{}, err
	}
	if check.Id != goodId {
		// 不存在该商品
		return models.BuyGoodResponse{}, models.ErrNoSuchGood
	}
	// 检查是否存在未完成的支付订单
	ok, err := s.goodQuery.CheckPayBeforeNew(ctx, userId)
	if err != nil {
		return models.BuyGoodResponse{}, err
	}
	if !ok {
		return models.BuyGoodResponse{}, models.ErrExistUnFinishedPayOrder
	}
	// 新建支付订单
	data, err := s.goodQuery.BuyGood(ctx, check, userId)
	if err != nil {
		return models.BuyGoodResponse{}, nil
	}
	return data, nil
}
