package main

type Voter struct {
	ID        string `json:"id"`
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
	FullName string `json:"full_name"`
	NIN      string `json:"nin"`
	DOB      string `json:"dob"`
	State    string `json:"state"`
	Lga      string `json:"lga"`
	Phone    string `json:"phone"`
}

type PaginatedMeta struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}
