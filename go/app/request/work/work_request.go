package work

type GetWorkByUserIdAndStatusRequest struct {
	CompanyId int    `form:"company_id" json:"company_id"`
	UserId    string `form:"user_id" json:"user_id" binding:"required"`
	Status    int    `form:"status" json:"status"`
	PageNum   int    `form:"page_num" json:"page_num"`
	PageSize  int    `form:"page_size" json:"page_size"`
}
