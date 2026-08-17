package service

import (
	"context"
	"log"
	"time"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/utils"
	"gorm.io/gorm"
)

type MemberService struct {
	db      *gorm.DB
	payment *PaymentService
}

func NewMemberService(db *gorm.DB, payment *PaymentService) *MemberService {
	return &MemberService{db: db, payment: payment}
}

var (
	ErrUniqueErr             = models.NewBusinessError(400, "已添加过该数据")
	ErrNoSuchData            = models.NewBusinessError(400, "找不到数据")
	ErrExistMemberData       = models.NewBusinessError(400, "会员数据已经存在")
	ErrPaidOrderCannotCancel = models.NewBusinessError(400, "已支付订单不可取消")
	ErrOrderFinished         = models.NewBusinessError(400, "订单已完成")
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
		if err == gorm.ErrDuplicatedKey {
			return "", ErrUniqueErr
		}
		return "", models.NewDatabaseErr(err)
	}
	return newData.UUID, nil
}

func (s *MemberService) UpdatePlanData(ctx context.Context, planId string, name string, value int64, des string) (int, error) {
	var newData models.MemberPlan
	newData.Name = name
	if value != 0 {
		newData.Value = value
	}
	newData.Description = des
	var res int
	// 在事务中处理
	err := s.db.Transaction(func(tx *gorm.DB) error {
		check, e := gorm.G[models.MemberPlan](tx).Select("uuid").Where("uuid = ?", planId).Limit(1).Find(ctx)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		if len(check) == 0 {
			return ErrNoSuchData
		}
		res, e = gorm.G[models.MemberPlan](tx).Where("uuid = ?", planId).Updates(ctx, newData)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (s *MemberService) GetMemberPlanData(ctx context.Context, planId string) (models.MemberPlanRes, error) {
	var datas []models.MemberPlan
	datas, err := gorm.G[models.MemberPlan](s.db).Where("uuid = ?", planId).Find(ctx)
	if err != nil {
		return models.MemberPlanRes{}, models.NewDatabaseErr(err)
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
	rawDatas, err := gorm.G[models.MemberPlan](s.db).Find(ctx)
	if err != nil {
		return nil, models.NewDatabaseErr(err)
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
		checkUser, e := gorm.G[models.User](tx).Select("uuid").Where("uuid = ?", userId).Limit(1).Find(ctx)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		if len(checkUser) == 0 {
			return ErrNoSuchUser
		}
		checkPlan, e = gorm.G[models.MemberPlan](tx).Select("uuid").Where("uuid = ?", planId).Limit(1).Find(ctx)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		if len(checkPlan) == 0 {
			return ErrNoSuchData
		}
		e = gorm.G[models.MemberOrders](tx).Create(ctx, &newData)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return newData.UUID, nil
}

// 针对MemberList的操作
func (s *MemberService) GetAllMemberListData(ctx context.Context, page int, pageSize int) ([]models.MemberListRes, error) {
	var datas []models.MemberListRes
	tx := s.db.Model(&models.MemberList{}).Select("uuid, user_id, status, end_at").
		WithContext(ctx).Scopes(utils.Paginate(page, pageSize)).Find(&datas)
	if tx.Error != nil {
		return nil, models.NewDatabaseErr(tx.Error)
	}
	if len(datas) == 0 {
		return nil, nil
	}
	return datas, nil
}

func (s *MemberService) GetMemberListDataByMemberId(ctx context.Context, memberId string) (models.MemberListRes, error) {
	var rawData []models.MemberList
	rawData, err := gorm.G[models.MemberList](s.db).Where("uuid = ?", memberId).Limit(1).Find(ctx)
	if err != nil {
		log.Println(err.Error())
		return models.MemberListRes{}, models.NewDatabaseErr(err)
	}
	if len(rawData) == 0 {
		return models.MemberListRes{}, ErrNoSuchData
	}
	data := rawData[0]
	var resData = models.MemberListRes{
		UUID:   data.UUID,
		UserId: data.UserId,
		Status: data.Status,
		EndAt:  data.EndAt,
	}
	return resData, nil
}

func (s *MemberService) GetUserMemberData(ctx context.Context, userId string) (models.MemberListRes, error) {
	var rawData []models.MemberList
	rawData, err := gorm.G[models.MemberList](s.db).Where("user_id = ?", userId).Limit(1).Find(ctx)
	if err != nil {
		log.Println(err.Error())
		return models.MemberListRes{}, models.NewDatabaseErr(err)
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

func (s *MemberService) AddNewMemberListData(ctx context.Context, userId string, endAt int64) (string, error) {
	// 表中user_id有唯一性约束
	var newData = models.MemberList{
		UserId: userId,
		Status: 1,
		EndAt:  endAt,
	}
	err := gorm.G[models.MemberList](s.db).Create(ctx, &newData)
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			return "", ErrExistMemberData
		}
		return "", models.NewDatabaseErr(err)
	}
	return newData.UUID, nil
}

func (s *MemberService) UpdateMemberListData(ctx context.Context, memberId string, status int64, endAt int64) (int, error) {
	var res int
	var newData models.MemberList
	newData.Status = status
	if newData.Status == 0 {
		newData.EndAt = time.Now().Unix()
	}
	if endAt == 0 && newData.Status != 0 {
		newData.EndAt = time.Now().Unix()
	} else if endAt > 0 && newData.Status != 0 {
		newData.EndAt = endAt
	}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 检查会员信息是否存在
		var check []models.MemberList
		check, e := gorm.G[models.MemberList](tx).Select("uuid").Where("uuid = ?", memberId).Find(ctx)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		if len(check) == 0 {
			return ErrNoSuchData
		}
		// 更新数据
		res, e = gorm.G[models.MemberList](tx).Where("uuid = ?", memberId).Updates(ctx, newData)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return res, nil
}

// 针对MemberOrder的操作
func (s *MemberService) GetUserMemberOrderData(ctx context.Context, userId string, page int, pageSize int) ([]models.MemberOrderRes, error) {
	var datas []models.MemberOrderRes
	tx := s.db.Model(&models.MemberOrders{}).WithContext(ctx).
		Select("uuid, plan_id, user_id, value, status").Where("user_id = ?", userId).
		Scopes(utils.Paginate(page, pageSize)).Find(&datas)
	if tx.Error != nil {
		return nil, models.NewDatabaseErr(tx.Error)
	}
	if len(datas) == 0 {
		return nil, nil
	}

	return datas, nil
}

func (s *MemberService) GetUserMemberOrderDataById(ctx context.Context, userId string, orderId string) (models.MemberOrderRes, error) {
	var rawData []models.MemberOrders
	rawData, err := gorm.G[models.MemberOrders](s.db).Where("uuid = ?", orderId).Limit(1).Find(ctx)
	if err != nil {
		return models.MemberOrderRes{}, models.NewDatabaseErr(err)
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

func (s *MemberService) GetAllMemberOrderData(ctx context.Context, page int, pageSize int) ([]models.MemberOrderRes, error) {
	var datas []models.MemberOrderRes
	tx := s.db.Model(&models.MemberOrders{}).WithContext(ctx).
		Select("uuid, plan_id, user_id, value, status").Scopes(utils.Paginate(page, pageSize)).Find(&datas)
	if tx.Error != nil {
		log.Println(tx.Error.Error())
		return nil, models.NewDatabaseErr(tx.Error)
	}
	if len(datas) == 0 {
		return nil, ErrNoSuchData
	}
	return datas, nil
}

func (s *MemberService) GetMemberOrderDataByOrderId(ctx context.Context, orderId string) (models.MemberOrderRes, error) {
	var rawData []models.MemberOrders
	rawData, err := gorm.G[models.MemberOrders](s.db).Where("uuid = ?", orderId).Limit(1).Find(ctx)
	if err != nil {
		return models.MemberOrderRes{}, models.NewDatabaseErr(err)
	}
	if len(rawData) == 0 {
		return models.MemberOrderRes{}, ErrNoSuchData
	}
	data := rawData[0]
	var resData = models.MemberOrderRes{
		UUID:   data.UUID,
		PlanId: data.PlanId,
		UserId: data.UserId,
		Value:  data.Value,
		Status: data.Status,
	}
	return resData, nil
}

func (s *MemberService) UpdateMemberOrderData(ctx context.Context, orderId string, status int64, value int64) (int, error) {
	var res int
	var newData models.MemberOrders
	newData.Status = status
	if status == models.MemberOrderStatusCanceled {
		newData.Value = 0
	} else {
		newData.Value = value
	}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 检查订单是否存在
		var check []models.MemberOrders
		check, e := gorm.G[models.MemberOrders](tx).Select("uuid").Where("uuid = ?", orderId).Limit(1).Find(ctx)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		if len(check) == 0 {
			return ErrNoSuchData
		}
		// 更新订单
		res, e = gorm.G[models.MemberOrders](tx).Where("uuid = ?", orderId).Updates(ctx, newData)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (s *MemberService) CancelMemberOrder(ctx context.Context, userId string, orderId string) (int, error) {
	var res int
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var check []models.MemberOrders
		check, e := gorm.G[models.MemberOrders](tx).Where("uuid = ?", orderId).Limit(1).Find(ctx)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		if len(check) == 0 {
			return ErrNoSuchData
		}
		data := check[0]
		if data.UserId != userId || data.Status == models.MemberOrderStatusCanceled {
			return ErrWrongUser
		}
		if data.Status == models.MemberOrderStatusPaid {
			return ErrPaidOrderCannotCancel
		}
		// 实际的更新操作
		res, e = gorm.G[models.MemberOrders](tx).Where("uuid = ?", orderId).Update(ctx, "status", models.MemberOrderStatusCanceled)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		// 联动作废该订单的未支付支付单
		if _, e := gorm.G[models.Pay](tx).Where("business_id = ? AND status = ?", orderId, models.PayStatusPending).
			Update(ctx, "status", models.PayStatusVoid); e != nil {
			return models.NewDatabaseErr(e)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (s *MemberService) FinishMemberOrder(ctx context.Context, userId string, orderId string) (*PaymentResult, error) {
	if !s.payment.GatewayConfigured() {
		return nil, ErrGatewayMissing
	}
	var pay *models.Pay
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var check []models.MemberOrders
		check, e := gorm.G[models.MemberOrders](tx).Where("uuid = ?", orderId).Limit(1).Find(ctx)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		if len(check) == 0 {
			return ErrNoSuchData
		}
		data := check[0]
		if data.UserId != userId {
			return ErrWrongUser
		}
		if data.Status == models.MemberOrderStatusPaid {
			return ErrOrderFinished
		}
		if data.Status != models.MemberOrderStatusCreated && data.Status != models.MemberOrderStatusPendingPay {
			return ErrWrongUser
		}
		// 已有有效待支付单则幂等返回
		if data.Status == models.MemberOrderStatusPendingPay {
			if active, err := s.payment.FindActivePendingPay(ctx, tx, orderId); err != nil {
				return err
			} else if active != nil {
				pay = active
				return nil
			}
		}
		// 通过planId获取价格
		var planData models.MemberPlan
		planData, e = gorm.G[models.MemberPlan](tx).Where("uuid = ?", data.PlanId).First(ctx)
		if e != nil {
			log.Println(e)
			return ErrNoSuchData
		}
		// 修改会员订单的状态为待支付
		_, e = gorm.G[models.MemberOrders](tx).Where("uuid = ?", orderId).Update(ctx, "status", models.MemberOrderStatusPendingPay)
		if e != nil {
			return models.NewDatabaseErr(e)
		}
		// 幂等创建支付单
		pay, e = s.payment.CreatePay(ctx, tx, models.PayBusinessMember, orderId, userId, planData.Value)
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.payment.BuildResult(ctx, pay)
}
