package services

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/1inkun/billApp/internal/models"
	"github.com/1inkun/billApp/internal/repositories"
)

type MemberService struct {
	memberQuery *repositories.MemberQuery
}

func NewMemberService(memberQuery *repositories.MemberQuery) *MemberService {
	return &MemberService{
		memberQuery: memberQuery,
	}
}

func (s *MemberService) GetMemberData(ctx context.Context, userId int64) (models.MemberData, error) {
	var data models.MemberData
	key := "Member:" + strconv.Itoa(int(userId))
	// 在缓存中查
	data, err := s.memberQuery.GetMemberDataFromCache(ctx, key)
	if err == nil {
		// 以k,v形式存在redis中的数据若不存在会报错
		return data, nil
	}
	// 缓存没命中，在数据库中查询
	data, err = s.memberQuery.GetMemberData(ctx, userId)
	if err != nil {
		return models.MemberData{}, err
	}
	dataStr, err := json.Marshal(&data)
	if err != nil {
		log.Printf("在GetMemberData中序列化数据时出错:%s", err)
		return models.MemberData{}, nil
	}
	s.memberQuery.SetCache(ctx, key, string(dataStr))
	return models.MemberData{}, nil
}
