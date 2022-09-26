package apis

func LoginEmail(email, password string) (*AuthResponse, error) {
	data := &SignInRequest{
		Device:   1,
		Email:    email,
		Password: password,
	}
	resp := &AuthResponse{}
	err := postJson("auth/signin", data, resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func LoginGuest() (*GuestAuthResponse, error) {
	data := &GuestRegisterRequest{
		Avatar:   "",
		Nickname: "",
	}
	resp := &GuestAuthResponse{}
	err := postJson("auth/guest/signin", data, resp)
	return resp, err
}

func RegisterGuest(nickname string, avatar string) (*AuthResponse, error) {
	data := &GuestRegisterRequest{
		Avatar:   avatar,
		Nickname: nickname,
	}
	resp := &AuthResponse{}
	err := postJson("auth/guest", data, resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Register(request *RegisterRequest) error {
	resp := &GuestAuthResponse{}
	err := postJson("auth/guest/signin", request, resp)
	return err
}

func GetVerifyCode(email string) (*GuestAuthResponse, error) {
	data := &VerifyCodeRequest{
		Email: email,
		Mode:  "",
	}
	resp := &GuestAuthResponse{}
	err := postJson("auth/verifyCode", data, resp)
	return resp, err
}
