package models

type UserRegistrationEvent struct {
	Name          string `json:"name"`
	Age           int    `json:"age"`
	City          string `json:"city"`
	Gender        string `json:"gender"`
	SearchGender  string `json:"search_gender"`
	SearchAgeFrom int    `json:"search_age_from"`
	SearchAgeTo   int    `json:"search_age_to"`
	Location      int    `json:"location"`
}
