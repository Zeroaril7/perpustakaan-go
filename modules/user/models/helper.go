package models

import "github.com/Zeroaril7/perpustakaan-go/pkg/utils"

func (m *UserAdd) ToUser(e User) User {
	e.Username = m.Username
	e.Password = utils.HashPassword(m.Password)
	e.Role = m.Role

	return e
}
