package carwise

type Interactor struct {
	services Services
}

func NewInteractor(svcs Services) *Interactor {
	return &Interactor{
		services: svcs,
	}
}

func (i *Interactor) CreateUser(request UserCreateRequest) []string {

	return nil
}

func (i *Interactor) LoginUser(request UserLoginRequest) (*User, []string) {

	return nil, nil
}

func (i *Interactor) LogoutUser(token string) []string {

	return nil
}

func (i *Interactor) IsTokenBlackListed(token string) (bool, []string) {

	return false, nil
}

func (i *Interactor) AddTokenBlackList(token string) []string {

	return nil
}
