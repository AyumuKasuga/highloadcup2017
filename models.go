package main

type user struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	BirthDate int    `json:"birth_date"`
}

type users struct {
	Users []user `json:"users"`
}

type location struct {
	Id       int    `json:"id"`
	Place    string `json:"place"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Distance int    `json:"distance"`
}

type locations struct {
	Locations []location `json:"locations"`
}

type visit struct {
	Id        int `json:"id"`
	Location  int `json:"location"`
	User      int `json:"user"`
	VisitedAt int `json:"visited_at"`
	Mark      int `json:"mark"`
}

type visits struct {
	Visits []visit `json:"visits"`
}

type userVisit struct {
	Mark      int    `json:"mark"`
	VisitedAt int    `json:"visited_at"`
	Place     string `json:"place"`
}
