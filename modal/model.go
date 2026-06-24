package model

type Voter struct {
	ID        string `json:"id"`
	VoterID   string `json:"voter_id"`
	FullName  string `json:"full_name"`
	NIN       string `json:"nin"`
	DOB       string `json:"dob"`
	State     string `json:"state"`
	Lga       string `json:"lga"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type RegisterVoterRequest struct {
	FullName string `json:"full_name" binding:"required,min=3"`
	NIN      string `json:"nin" binding:"required, numeric"`
	DOB      string `json:"dob" binding:"required"`
	State    string `json:"state" binding:"required"`
	Lga      string `json:"lga" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

type PaginatedMeta struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

type UpdateVoterStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
