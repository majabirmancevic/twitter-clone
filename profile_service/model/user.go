package model

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

//type SignUpRegularRequest struct {
//	Name            string `bson:"name" json:"name" binding:"required"`
//	Lastname        string `bson:"lastname" json:"lastname" binding:"required"`
//	Gender          string `bson:"gender" json:"gender" binding:"required"`
//	Age             uint16 `bson:"age" json:"age" binding:"required"`
//	PlaceOfLiving   string `bson:"placeOfLiving" json:"placeOfLiving" binding:"required"`
//	Username        string `bson:"username" json:"username" binding:"required"`
//	Password        string `bson:"password" json:"password" binding:"required,min=8"`
//	PasswordConfirm string `bson:"passwordConfirm,omitempty" json:"passwordConfirm" binding:"required"`
//	Role            string `json:"role" bson:"role"`
//}

type RegularProfiles []*RegularProfile
type BusinessProfiles []*BusinessProfile

// Model for create and read user from db
type RegularProfile struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	Name             string             `bson:"name" json:"name" validate:"required,min=2,max=30"`
	Lastname         string             `bson:"lastname" json:"lastname" validate:"required,min=2,max=30"`
	Gender           string             `bson:"gender" json:"gender" validate:"required"`
	Age              int32              `bson:"age" json:"age" validate:"required,min=13"`
	PlaceOfLiving    string             `bson:"placeOfLiving" json:"placeOfLiving" validate:"required"`
	Email            string             `json:"email" bson:"email" validate:"email,required"`
	Username         string             `bson:"username" json:"username" validate:"required"`
	Password         string             `bson:"password" json:"password" validate:"required,min=8"`
	VerificationCode string             `bson:"verificationCode" json:"verificationCode" `
	Verified         bool               `bson:"verified" json:"verified" `
	Role             string             `json:"role" bson:"role"`
}

// Response for client without password
type DBRegularResponse struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Name          string             `bson:"name" json:"name" validate:"required,min=2,max=30"`
	Lastname      string             `bson:"lastname" json:"lastname" validate:"required,min=2,max=30"`
	Gender        string             `bson:"gender" json:"gender" validate:"required"`
	Age           int32              `bson:"age" json:"age" validate:"required,min=13"`
	PlaceOfLiving string             `bson:"placeOfLiving" json:"placeOfLiving" validate:"required"`
	Email         string             `json:"email" bson:"email" validate:"email,required"`
	Username      string             `bson:"username" json:"username" validate:"required"`
	//VerificationCode string             `bson:"verificationCode" json:"verificationCode" `
	Verified bool   `bson:"verified" json:"verified" `
	Role     string `json:"role" bson:"role"`
}

type SignInRequest struct {
	Username string `bson:"username" json:"username" binding:"required"`
	Password string `bson:"password" json:"password" binding:"required"`
}

type SignInResponseRegular struct {
	Token          string            `bson:"token" json:"token" `
	RegularProfile DBRegularResponse `bson:"user" json:"user" `
}

// -------------------------------------------------------------------------------------------------
func (p *PasswordDto) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (p *PasswordDto) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

//--------------------------------------------------------------------------------------------------

type PasswordDto struct {
	OldPassword string `bson:"oldPassword" json:"oldPassword"`
	NewPassword string `bson:"newPassword" json:"newPassword"`
}

func NewUserResponse(user *RegularProfile) DBRegularResponse {
	return DBRegularResponse{
		ID:            user.ID,
		Name:          user.Name,
		Lastname:      user.Lastname,
		Gender:        user.Gender,
		Age:           user.Age,
		PlaceOfLiving: user.PlaceOfLiving,
		Email:         user.Email,
		Username:      user.Username,
		//VerificationCode: user.VerificationCode,
		Verified: user.Verified,
		Role:     user.Role,
	}
}

// -------------------------------------------------------------------------------------------------------------
type BusinessProfile struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	CompanyName      string             `bson:"companyName" json:"companyName" validate:"required,min=2,max=30"`
	Email            string             `bson:"email" json:"email" validate:"email,required"`
	WebSite          string             `bson:"webSite" json:"webSite" validate:"required"`
	Username         string             `bson:"username" json:"username" validate:"required"`
	Password         string             `bson:"password" json:"password" validate:"required,min=8"`
	VerificationCode string             `bson:"verificationCode" json:"verificationCode" `
	Verified         bool               `bson:"verified" json:"verified" `
	Role             string             `json:"role" bson:"role"`
}

func (p *BusinessProfile) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (p *BusinessProfile) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}
func (p *BusinessProfiles) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

//---------------------------------------------------------------

func (p *RegularProfiles) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *RegularProfile) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *RegularProfile) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
func (p *SignInRequest) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
