package models

import "time"

// Official represents a city official
type Official struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Ward      *int      `json:"ward,omitempty"`
	Party     string    `json:"party"`
	Role      string    `json:"role"`
	Contact   string    `json:"contact"`
	Email     string    `json:"email"`
	ImageURL  string    `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// VotingRecord represents a voting record for an official
type VotingRecord struct {
	ID          int       `json:"id"`
	OfficialID  int       `json:"official_id"`
	BillTitle   string    `json:"bill_title"`
	Vote        string    `json:"vote"` // "yes", "no", "abstain"
	VoteDate    time.Time `json:"vote_date"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// Committee represents a city committee
type Committee struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// OfficialCommittee represents the relationship between officials and committees
type OfficialCommittee struct {
	ID          int       `json:"id"`
	OfficialID  int       `json:"official_id"`
	CommitteeID int       `json:"committee_id"`
	Role        string    `json:"role"` // "member", "chair", "vice-chair"
	CreatedAt   time.Time `json:"created_at"`
}

// WardStatistic represents statistics for a ward
type WardStatistic struct {
	ID                    int       `json:"id"`
	Ward                  int       `json:"ward"`
	InfrastructureSpending float64  `json:"infrastructure_spending"`
	Response311Time       float64   `json:"response_311_time"` // in days
	AffordableHousing     int       `json:"affordable_housing"` // units
	PermitProcessing      float64   `json:"permit_processing"` // in days
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
