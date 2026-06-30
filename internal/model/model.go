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
	NIN      string `json:"nin" binding:"required,len=11,numeric"`
	DOB      string `json:"dob" binding:"required"`
	State    string `json:"state" binding:"required"`
	Lga      string `json:"lga" binding:"required"`
	Phone    string `json:"phone" binding:"required,min=11"`
}

type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// omitempty
type RegisterUserResponse struct {
	Email string `json:"email" binding:"required,email"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// omitempty
type LoginUserResponse struct {
	Email string `json:"email" binding:"required,email"`
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

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
