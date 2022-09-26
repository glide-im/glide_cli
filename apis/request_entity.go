package apis

type SignInRequest struct {
	Device   int64  `json:"device"`
	Email    string `json:"email"  validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Account  string `json:"account"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=16,min=6"`
	Captcha  string `json:"captcha" validate:"required"`
	Nickname string `json:"nickname" validate:"required"`
}

type GuestRegisterRequest struct {
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
}

type GuestRegisterV2Request struct {
	FingerprintId string `json:"fingerprint_id" validate:"required"`
	Origin        string `json:"origin"`
}

type VerifyCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
	Mode  string `json:"mode"`
}
