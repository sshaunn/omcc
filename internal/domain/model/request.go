package model

type UpdateCustomerStatusRequest struct {
	CustomerId   string  `json:"customer_id" binding:"required"`
	SocialId     string  `json:"social_id" binding:"required"`
	Status       *string `json:"status" binding:"omitempty"`
	MemberStatus *string `json:"member_status" binding:"omitempty"`
}

type DeleteCustomerRequest struct {
	IdList []string `json:"id_list" binding:"required"`
}
