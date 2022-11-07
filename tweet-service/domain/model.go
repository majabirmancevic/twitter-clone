package domain

type RegularProfile struct {
	Name          string
	Lastname      string
	Gender        string
	Age           int32
	PlaceOfLiving string
	Username      string //mora biti jedinstveno
	Password      string
	IsPrivate     bool
	Tweets        []Tweet
}

type BusinessProfile struct {
	CompanyName string
	Email       string
	WebSite     string
	Username    string //mora biti jedinstveno
	Password    string
	Tweets      []Tweet
}

type Tweet struct {
	Name        string
	Description string
	//Profile     Profile
	NumberOfLikes int32
	Likes         []RegularProfile
}
