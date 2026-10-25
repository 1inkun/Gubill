package services

import (
	"context"
	"sync"

	"github.com/1inkun/billApp/internal/models"
	"github.com/1inkun/billApp/internal/repositories"
	"golang.org/x/sync/errgroup"
)

type StatusService struct {
	statusQuery *repositories.StatusQuery
	mu          *sync.Mutex
}

func NewStatusService(statusQuery *repositories.StatusQuery) *StatusService {
	var mu sync.Mutex
	return &StatusService{
		statusQuery: statusQuery,
		mu:          &mu,
	}
}

func (s *StatusService) GetSginStatus(ctx context.Context) (any, error) {
	var data []models.CurrentUser
	// 先从Redis中尝试获取数据
	s.mu.Lock()
	data, _ = s.statusQuery.GetUserDataFromCache(ctx)
	s.mu.Unlock()
	if len(data) != 0 {
		// data长度不为0表示找到缓存的数据
		return data, nil
	} else {
		// 先向order表查询,获取正在签到用户的用户id
		userIds, err := s.statusQuery.GetCurrentUser(ctx)
		if err != nil {
			return nil, err
		}
		if userIds == nil {
			// id数组为空则表示没有正在签到的用户
			return nil, nil
		}
		var eg errgroup.Group
		var mu sync.Mutex
		eg.Go(func() error {
			for _, v := range userIds {
				userData, err := s.statusQuery.GetUserData(ctx, v)
				if err != nil {
					return err
				}
				mu.Lock()
				data = append(data, userData)
				mu.Unlock()
			}
			return nil
		})
		if err := eg.Wait(); err != nil {
			return nil, err
		}
	}
	return data, nil
}
