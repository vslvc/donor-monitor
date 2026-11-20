package model

type DonorCenter struct{
	ID	    int64
	NameCenter  string
	City	    string
	URL         string
	Statuses    []DonorStatus
}
