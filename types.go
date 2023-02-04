package main

// )

// type Account struct {
// 	ID                int       `json:"id"`
// 	FirstName         string    `json:"name"`
// 	Email             string    `json:"email"`
// 	EncryptedPassword []byte    `json:"password"`
// 	Token             string    `json:"token"`
// 	CreatedAt         time.Time `json:"created_at"`
// }

// // func (a *Account) ValidPassword(pw string) bool {
// // 	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw)) == nil
// // }

// func newAccount(firstName, email, password string) (*Account, error) {
// 	// hash password
// 	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Account{
// 		FirstName:         firstName,
// 		Email:             email,
// 		EncryptedPassword: encpw,
// 		CreatedAt:         time.Now(),
// 	}, nil
// }
