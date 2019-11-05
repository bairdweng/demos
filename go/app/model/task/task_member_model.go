package task

import "time"

type Member struct {
	ID                 int64      `gorm:"primary_key" json:" - "` //
	TaskId             int64      `json:"task_id"`                //任务id
	OrderId            string     `json:"order_id"`               //流水id
	PublisherId        int64      `json:"publisher_id"`           //发布者id
	ParticipantId      int64      `json:"participant_id"`         //参加者id
	Source             int        `json:"source"`                 //任务来源:1申请,2邀请,默认申请
	Progress           int        `json:"progress"`               //任务状态:1:邀请中，2已申请，3已同意参加，4申请完成，5已同意完成，6已评分，7发布者已评分，8互评,9拒绝申请,10拒绝完成,20踢出任务
	Reward             float64    `json:"reward"`                 //赏金字段
	KickOutAt          *time.Time `json:"kick_out_at"`            //踢出任务截止时间
	AutoCompleteAt     *time.Time `json:"auto_complete_at"`       //任务自动完成截止时间
	ParticipateAt      *time.Time `json:"participate_at"`         //任务参加时间
	FinishAt           *time.Time `json:"finish_at"`              //任务完成时间
	PublisherScore     int        `json:"publisher_score"`        //主态评分
	PublisherContent   string     `json:"publisher_content"`      //主态评论
	ParticipantScore   int        `json:"participant_score"`      //客态评论
	ParticipantContent string     `json:"participant_content"`    //客态评论
	ProofFileUrls      string   `json:"proof_file_urls"`        //凭证id:文件服务url
	RejectReason       string     `json:"reject_reason"`          //拒绝完成理由
	Extend             string     `json:"extend"`                 //冗余字段
	CreatedAt          *time.Time `json:"created_at"`             //
	UpdatedAt          *time.Time `json:"updated_at"`             //
	DeletedAt          *time.Time `json:"deleted_at"`             //软删除字段
}
func (Member) TableName() string {
	return "iq_task_member"
}