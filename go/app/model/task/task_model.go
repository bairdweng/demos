package task

import (
	"time"
)

type Task struct {
	ID               int64      `gorm:"primary_key" json:"id"`             //
	TaskId           int64      `json:"task_id"`                           //任务id
	Category         int        `json:"category"`                          //任务分类:1任务,2岗位
	Progress         int        `gorm:"default:'1'" json:"progress"`       //任务状态:1进行中,2过期,3完成
	PayStatus        int        `json:"pay_status"`                        //支付状态: 余额是否已返回等
	Status           int        `gorm:"default:'1'" json:"status"`         //是否正常:禁用0/启用1
	Quota            int        `json:"quota"`                             //可完成人/次数
	SingleRewardMin  float64    `json:"single_reward_min"`                 //单次任务赏金最小值
	SingleRewardMax  float64    `json:"single_reward_max"`                 //单次任务赏金最大值
	IsPublic         int        `gorm:"default:'1'" json:"is_public"`      //任务模式:1公开0私密,默认公开
	IsCanComment     int        `gorm:"default:'1'" json:"is_can_comment"` //评论模式:是否可以评论,1是0否,默认1
	IsNeedProof      int        `gorm:"default:'0'"  json:"is_need_proof"` //凭证模式:任务是否需要凭证,1是0否,默认0
	ProofType        int        `json:"proof_type"`                        //凭证类型:1视频/图片,2其他
	ProofDescription string     `json:"proof_description"`                 //凭证描述:上传凭证页面展示描述
	Extend           string     `json:"extend"`                            //冗余字段
	CreatedAt        time.Time  `json:"created_at"`                        //
	UpdatedAt        time.Time  `json:"updated_at"`                        //
	DeletedAt        *time.Time `json:"deleted_at"`                        //软删除
}

func (Task) TableName() string {
	return "iq_task"
}
