package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/utils"
	"gorm.io/gorm"
)

type MemberService struct {
	db *gorm.DB
}

func NewMemberService(db *gorm.DB) *MemberService {
	return &MemberService{db: db}
}

var (
	ErrUniqueErr  = models.NewBusinessError(400, "已添加过该数据")
	ErrNoSuchData = models.NewBusinessError(400, "找不到数据")
)

// 针对MemberPlan的操作
func (s *MemberService) AddNewMemberPlan(ctx context.Context, name string, planType string, value int64, des string) (string, error) {
	var newData models.MemberPlan
	newData.Name = name
	newData.Type = planType
	newData.Value = value
	newData.Description = des
	// 直接利用Type的唯一性标签
	err := gorm.G[models.MemberPlan](s.db).Create(ctx, &newData)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return "", ErrUniqueErr
		}
		return "", ErrDatabaseErr
	}
	return newData.UUID, nil
}

func (s *MemberService) UpdatePlanData(ctx context.Context, planId string, name string, value int64, des string) (int, error) {
	var newData models.MemberPlan
	newData.Name = name
	newData.Value = value
	newData.Description = des
	var res int
	// 在事务中处理
	err := s.db.Transaction(func(tx *gorm.DB) error {
		check, e := gorm.G[models.MemberPlan](tx).Where("uuid = ?", planId).Limit(1).Find(ctx)
		if e != nil {
			log.Panicln(e.Error())
			return e
		}
		if len(check) == 0 {
			return ErrNoSuchData
		}
		res, e = gorm.G[models.MemberPlan](tx).Where("uuid = ?", planId).Updates(ctx, newData)
		if e != nil {
			log.Panicln(e.Error())
			return e
		}
		return nil
	})
	if err != nil {
		return 0, ErrDatabaseErr
	}
	return res, nil
}

func (s *MemberService) GetMemberPlanData(ctx context.Context, planId string) (models.MemberPlanRes, error) {
	var datas []models.MemberPlan
	datas, err := gorm.G[models.MemberPlan](s.db).Where("uuid = ?", planId).Find(ctx)
	if err != nil {
		return models.MemberPlanRes{}, ErrDatabaseErr
	}
	if len(datas) == 0 {
		return models.MemberPlanRes{}, ErrNoSuchData
	}
	data := datas[0]
	var resData models.MemberPlanRes
	resData.UUID = data.UUID
	resData.Name = data.Name
	resData.Type = data.Type
	resData.Value = data.Value
	resData.Description = data.Description
	return resData, nil
}

func (s *MemberService) GetAllMemberPlans(ctx context.Context) ([]models.MemberPlanRes, error) {
	// 预计不会有很多数据,不做分页查询
	var rawDatas []models.MemberPlan
	rawDatas, err := gorm.G[models.MemberPlan](s.db).Limit(5).Find(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil, ErrDatabaseErr
	}
	if len(rawDatas) == 0 {
		return nil, nil
	}
	// 创建响应数据数组
	var resData []models.MemberPlanRes
	for _, v := range rawDatas {
		var temp = models.MemberPlanRes{
			UUID:        v.UUID,
			Name:        v.Name,
			Type:        v.Type,
			Value:       v.Value,
			Description: v.Description,
		}
		resData = append(resData, temp)
	}
	return resData, nil
}

func (s *MemberService) GenMemberPlanOrder(ctx context.Context, userId string, planId string) (string, error) {
	// 新建订单
	var newData models.MemberOrders
	newData.PlanId = planId
	newData.UserId = userId
	newData.Status = 0
	// 在事务中处理,确保一致性
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 检查用户和会员计划确实存在
		var checkUser []models.User
		var checkPlan []models.MemberPlan
		checkUser, e := gorm.G[models.User](tx).Where("uuid = ?", userId).Limit(1).Find(ctx)
		if e != nil {
			log.Println(e.Error())
			return e
		}
		if len(checkUser) == 0 {
			return ErrNoSuchUser
		}
		checkPlan, e = gorm.G[models.MemberPlan](tx).Where("uuid = ?", planId).Limit(1).Find(ctx)
		if e != nil {
			log.Println(e.Error())
			return e
		}
		if len(checkPlan) == 0 {
			return ErrNoSuchData
		}
		e = gorm.G[models.MemberOrders](tx).Create(ctx, &newData)
		if e != nil {
			log.Println(e.Error())
			return ErrDatabaseErr
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return newData.UUID, nil
}

// 针对MemberList的操作
func (s *MemberService) GetUserMemberData(ctx context.Context, userId string) (models.MemberListRes, error) {
	var rawData []models.MemberList
	rawData, err := gorm.G[models.MemberList](s.db).Where("uuid = ?", userId).Limit(1).Find(ctx)
	if err != nil {
		log.Println(err.Error())
		return models.MemberListRes{}, ErrDatabaseErr
	}
	if len(rawData) == 0 {
		return models.MemberListRes{}, nil
	}
	data := rawData[0]
	var resData = models.MemberListRes{
		UUID:   data.UUID,
		UserId: userId,
		Status: data.Status,
		EndAt:  data.EndAt,
	}
	return resData, nil
}

// 针对MemberOrder的操作
func (s *MemberService) GetUserMemberOrderData(ctx context.Context, userId string, page int, pageSize int) ([]models.MemberOrderRes, error) {
	var rawData []models.MemberOrders
	tx := s.db.Model(&models.MemberOrders{}).WithContext(ctx).Where("user_id = ?", userId).Scopes(utils.Paginate(page, pageSize)).Find(&rawData)
	if tx.Error != nil {
		log.Println(tx.Error.Error())
		return nil, ErrDatabaseErr
	}
	if len(rawData) == 0 {
		return nil, nil
	}
	var resData []models.MemberOrderRes
	for _, v := range rawData {
		var temp = models.MemberOrderRes{
			UUID:   v.UUID,
			PlanId: v.PlanId,
			UserId: v.UserId,
			Value:  v.Value,
			Status: v.Status,
		}
		resData = append(resData, temp)
	}
	return resData, nil
}

func (s *MemberService) GetUserMemberOrderDataById(ctx context.Context, userId string, orderId string) (models.MemberOrderRes, error) {
	var rawData []models.MemberOrders
	rawData, err := gorm.G[models.MemberOrders](s.db).Where("uuid = ?", orderId).Limit(1).Find(ctx)
	if err != nil {
		log.Println(err.Error())
		return models.MemberOrderRes{}, ErrDatabaseErr
	}
	if len(rawData) == 0 {
		return models.MemberOrderRes{}, ErrNoSuchData
	}
	data := rawData[0]
	if data.UserId != userId {
		return models.MemberOrderRes{}, ErrWrongUser
	}
	var resData = models.MemberOrderRes{
		UUID:   data.UUID,
		PlanId: data.PlanId,
		UserId: data.UserId,
		Value:  data.Value,
		Status: data.Status,
	}
	return resData, nil
}

func (s *MemberService) CancelMemberOrder(ctx context.Context, userId string, orderId string) (int, error) {
	var res int
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var check []models.MemberOrders
		check, e := gorm.G[models.MemberOrders](tx).Where("uuid = ?", orderId).Limit(1).Find(ctx)
		if e != nil {
			log.Println(e.Error())
			return ErrDatabaseErr
		}
		if len(check) == 0 {
			return ErrNoSuchData
		}
		data := check[0]
		if data.UserId != userId || data.Status == -1 {
			return ErrWrongUser
		}
		// 实际的更新操作
		res, e = gorm.G[models.MemberOrders](tx).Where("uuid = ?", orderId).Update(ctx, "status", -1)
		return nil
	})
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (s *MemberService) FinishMemberOrder(ctx context.Context, userId string, orderId string) (string, error) {
	var newData models.Pay
	newData.BusinessType = "member"
	newData.BusinessId = orderId
	newData.Status = 0
	newData.ExpireTime = time.Now().Add(30 * time.Minute).Unix()
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var check []models.MemberOrders
		check, e := gorm.G[models.MemberOrders](tx).Where("uuid = ?", orderId).Limit(1).Find(ctx)
		if e != nil {
			return ErrDatabaseErr
		}
		if len(check) == 0 {
			return ErrNoSuchData
		}
		data := check[0]
		if data.UserId != userId || data.Status != 0 {
			return ErrWrongUser
		}
		// 通过planId获取价格
		var planData models.MemberPlan
		planData, e = gorm.G[models.MemberPlan](tx).Where("uuid = ?", data.PlanId).First(ctx)
		if e != nil {
			log.Println(e)
			return ErrNoSuchData
		}
		newData.Value = planData.Value
		// 新建支付订单
		e = gorm.G[models.Pay](tx).Create(ctx, &newData)
		return nil
	})
	if err != nil {
		return "", err
	}
	return newData.UUID, nil
}
