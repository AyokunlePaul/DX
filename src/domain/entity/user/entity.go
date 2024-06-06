package user

import (
	"DX/src/domain/entity"
	"DX/src/pkg/response"
	"DX/src/utils"
	"errors"
	"strings"
	"time"
)

type Client string

type Type int

type Account struct {
	Name        string `json:"name" bson:"name"`
	Number      string `json:"number" bson:"number"`
	Type        string `json:"type" bson:"type"`
	BankCode    string `json:"bank_code" bson:"bank_code"`
	CountryCode string `json:"country_code" bson:"country_code"`
}

type User struct {
	Id                        entity.DatabaseId   `json:"id" bson:"_id"`
	FirstName                 string              `json:"first_name" bson:"first_name"`
	LastName                  string              `json:"last_name" bson:"last_name"`
	Client                    Client              `json:"client" bson:"client"`
	CreatedBy                 entity.By           `json:"created_by" bson:"created_by"`
	ModifiedBy                []entity.ModifiedBy `json:"-" bson:"modified_by"`
	AdminId                   string              `json:"-" bson:"admin_id"`
	Email                     string              `json:"email,omitempty" bson:"email"`
	ProfilePicture            string              `json:"profile_picture,omitempty" bson:"profile_picture,omitempty"`
	Password                  string              `json:"-" bson:"password"`
	Token                     string              `json:"token" bson:"token"`
	RefreshToken              string              `json:"-" bson:"refresh_token,omitempty"`
	CategoryInterest          []string            `json:"category_interest,omitempty" bson:"category_interest"`
	AccountNumbers            []Account           `json:"account_numbers" bson:"account_numbers"`
	UserType                  Type                `json:"-" bson:"user_type"`
	Type                      string              `json:"-" bson:"type"`
	PhoneNumber               string              `json:"phone_number" bson:"phone_number"`
	Verification              int                 `json:"verification" bson:"verification,omitempty"`
	HasVerifiedPhone          bool                `json:"-" bson:"has_verified_phone"`
	HasVerifiedBankingDetails bool                `json:"-" bson:"has_verified_banking_details"`
	HasVerifiedEmail          bool                `json:"-" bson:"has_verified_email"`
	HasVerifiedAddress        bool                `json:"-" bson:"has_verified_address"`
	HasTransactionPin         bool                `json:"-" bson:"has_transaction_pin"`
	UserId                    string              `json:"-" bson:"user_id"`
	Ratings                   []float64           `json:"-" bson:"ratings"`
	Rating                    float64             `json:"rating" bson:"rating"`
	ErrandsCompleted          int64               `json:"errands_completed" bson:"errands_completed"`
	ErrandsCancelled          int64               `json:"errands_cancelled" bson:"errands_cancelled"`
	IsSuspended               bool                `json:"is_suspended" bson:"is_suspended"`
	SuspendedBy               string              `json:"-" bson:"is_suspended_by"`
	IsDeleted                 bool                `json:"-" bson:"is_deleted"`
	DeletedBy                 string              `json:"-" bson:"is_deleted_by"`
	CreationDate              string              `json:"creation_date,omitempty" bson:"creation_date"`
	CreatedAt                 time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt                 time.Time           `json:"updated_at" bson:"updated_at"`
}

const (
	Android Client = "android"
	iOS     Client = "ios"
	Web     Client = "web"
)

const (
	Normal Type = iota
	ClientManager
	Admin
	SuperAdmin
)

func (t Type) String() string {
	if t == Normal {
		return "User"
	}
	if t == ClientManager {
		return "Client Manager"
	}
	if t == Admin {
		return "Admin"
	}
	if t == SuperAdmin {
		return "Super Admin"
	}
	return ""
}

func (t Type) Id() string {
	if t == Normal {
		return "user"
	}
	if t == ClientManager {
		return "client-manager"
	}
	if t == Admin {
		return "admin"
	}
	if t == SuperAdmin {
		return "super-admin"
	}
	return ""
}

func Create(data map[string]interface{}) (*User, error) {
	var firstName, lastName, phone, password, client string
	var catInterest []string
	var ok bool
	currentTime := time.Now()

	nUser := &User{
		Id:               entity.NewDatabaseId(),
		UserId:           entity.NewDefaultId().String(),
		UserType:         Normal,
		Type:             Normal.Id(),
		CreatedAt:        currentTime,
		UpdatedAt:        currentTime,
		CreatedBy:        entity.ByUser,
		Ratings:          []float64{},
		Rating:           0.0,
		ErrandsCompleted: 0,
		AccountNumbers:   []Account{},
		ModifiedBy:       []entity.ModifiedBy{},
	}

	if firstName, ok = data["first_name"].(string); !ok {
		return nil, errors.New("first name is required")
	}
	if lastName, ok = data["last_name"].(string); !ok {
		return nil, errors.New("last name is required")
	}
	if phone, ok = data["phone_number"].(string); !ok {
		return nil, errors.New("phone number is required")
	}
	if password, ok = data["password"].(string); !ok {
		return nil, errors.New("password is required")
	}
	if client, ok = data["client"].(string); !ok {
		return nil, errors.New("client is required")
	}
	if _, ok = data["admin"].(bool); ok {
		nUser.UserType = Admin
		nUser.Type = Admin.Id()
	}
	if tInterests, ok := data["interests"].([]interface{}); ok {
		for _, interest := range tInterests {
			catInterest = append(catInterest, interest.(string))
		}
	}

	if client != string(Web) && client != string(Android) && client != string(iOS) {
		return nil, errors.New("invalid client type")
	}

	nUser.FirstName = strings.TrimSpace(firstName)
	nUser.LastName = strings.TrimSpace(lastName)
	nUser.PhoneNumber = strings.TrimSpace(phone)
	nUser.Password = strings.TrimSpace(password)
	nUser.Client = Client(strings.TrimSpace(client))

	return nUser, nil
}

func CreateForAdmin(data map[string]interface{}) (*User, error) {
	var firstName, lastName, phone, email string
	var ok bool
	currentTime := time.Now()

	nUser := &User{
		Id:               entity.NewDatabaseId(),
		UserId:           entity.NewDefaultId().String(),
		UserType:         Normal,
		Type:             Normal.Id(),
		CreatedBy:        entity.ByAdmin,
		CreatedAt:        currentTime,
		UpdatedAt:        currentTime,
		CreationDate:     currentTime.Format("January 02, 2006"),
		AccountNumbers:   []Account{},
		CategoryInterest: []string{},
		ModifiedBy:       []entity.ModifiedBy{},
		Ratings:          []float64{},
	}

	if firstName, ok = data["first_name"].(string); !ok {
		return nil, errors.New("first name is required")
	}
	if lastName, ok = data["last_name"].(string); !ok {
		return nil, errors.New("last name is required")
	}
	if email, ok = data["email"].(string); ok {
		nUser.Email = strings.TrimSpace(email)
	}
	if phone, ok = data["phone_number"].(string); !ok {
		return nil, errors.New("phone number is required")
	} else {
		if !utils.IsValidPhoneNumber(phone) {
			return nil, errors.New("phone number is not valid")
		}
	}
	if tInterests, ok := data["interests"].([]interface{}); ok {
		for _, interest := range tInterests {
			nUser.CategoryInterest = append(nUser.CategoryInterest, interest.(string))
		}
	}
	if tAccount, ok := data["account"].(map[string]interface{}); ok {
		account, err := NewAccountDetails(tAccount)
		if err != nil {
			return nil, err
		}
		nUser.AccountNumbers = append(nUser.AccountNumbers, account)
	}

	nUser.FirstName = strings.TrimSpace(firstName)
	nUser.LastName = strings.TrimSpace(lastName)
	nUser.PhoneNumber = strings.TrimSpace(phone)
	nUser.HasVerifiedPhone = true
	nUser.UpdateVerification()
	nUser.Client = Web

	return nUser, nil
}

func NewAccountDetails(data map[string]interface{}) (Account, error) {
	account := Account{
		CountryCode: "NG",
	}
	if name, ok := data["name"].(string); !ok {
		return Account{}, errors.New("account name is required")
	} else {
		account.Name = name
	}
	if number, ok := data["number"].(string); !ok {
		return Account{}, errors.New("account number is required")
	} else {
		account.Number = number
	}
	if aType, ok := data["type"].(string); !ok {
		return Account{}, errors.New("account type is required")
	} else {
		account.Type = aType
	}
	if bankCode, ok := data["bank_code"].(string); !ok {
		return Account{}, errors.New("bank code is required")
	} else {
		account.BankCode = bankCode
	}

	return account, nil
}

func UpdateForAdmin(data map[string]interface{}) (*User, error) {
	var firstName, lastName, phone, email string
	var ok bool

	nUser := &User{}

	if firstName, ok = data["first_name"].(string); ok {
		nUser.FirstName = strings.TrimSpace(firstName)
	}
	if lastName, ok = data["last_name"].(string); ok {
		nUser.LastName = strings.TrimSpace(lastName)
	}
	if email, ok = data["email"].(string); ok {
		nUser.Email = strings.TrimSpace(email)
	}
	if phone, ok = data["phone_number"].(string); ok {
		if !utils.IsValidPhoneNumber(phone) {
			return nil, errors.New("phone number is not valid")
		} else {
			nUser.PhoneNumber = strings.TrimSpace(phone)
		}
	}

	return nUser, nil
}

func CreateForLogin(data map[string]interface{}) (*User, error) {
	var phone, password string
	var ok bool

	if phone, ok = data["phone_number"].(string); !ok {
		return nil, errors.New("phone number is required")
	}
	if password, ok = data["password"].(string); !ok {
		return nil, errors.New("password is required")
	}

	nUser := &User{
		PhoneNumber: phone,
		Password:    password,
	}
	return nUser, nil
}

func (u *User) IsValidForInitialCreation() *response.BaseResponse {
	if u.IsAdmin() && u.Client != Web {
		return response.NewBadRequestError("admin can only be created on web clients")
	}
	return nil
}

func (u *User) IsOffline() bool {
	return u.CreatedBy == entity.ByAdmin && !u.IsSuspended && !u.IsDeleted
}

func (u *User) IsSuperAdmin() bool {
	return u.UserType == SuperAdmin
}

func (u *User) IsAdmin() bool {
	return u.UserType == Admin
}

func (u *User) UpdateVerification() {
	percentage := 0
	if u.HasVerifiedPhone == true {
		percentage += 30
	}
	if u.HasVerifiedBankingDetails == true {
		percentage += 20
	}
	if u.HasVerifiedEmail == true {
		percentage += 10
	}
	if u.HasVerifiedAddress == true {
		percentage += 20
	}
	u.Verification = percentage
}

func (u *User) UpdateUserDataForAdmin(nUser *User, adminId string) {
	if nUser.FirstName != "" {
		u.FirstName = nUser.FirstName
	}
	if nUser.LastName != "" {
		u.LastName = nUser.LastName
	}
	if nUser.Email != "" {
		u.Email = nUser.Email
	}
	if nUser.PhoneNumber != "" {
		u.PhoneNumber = nUser.PhoneNumber
	}
	modBy := u.ModifiedBy
	if modBy == nil {
		modBy = []entity.ModifiedBy{}
	}
	nMod := entity.ModifiedBy{
		Id:   adminId,
		Date: time.Now(),
	}
	modBy = append(modBy, nMod)
	u.ModifiedBy = modBy
}
