package users

type UsersDispatcher struct {
	users map[string]UserAccount
}

func (p *UsersDispatcher) IsNewUser(id string) bool {
	_, exist := p.users[id]
	return exist
}

func (p *UsersDispatcher) CreateNewUser(u *UserAccount) {
	if p.IsNewUser(u.id) {

	}
}

func (p *UsersDispatcher) RequestSession(ip string, id string) bool {
	var temp UserAccount = p.users[id]
	temp.hasActiveSession = true
	p.users[id] = temp
	return false
}
