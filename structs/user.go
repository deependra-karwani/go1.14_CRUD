package structs

/** Body Request */
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FCM      string `json:"fcm"`
}

type Forgot struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserId struct {
	UserId int `json:"userid"`
}

/** ~Body Request */

/** Query Results */
type ForSession struct {
	UserId         int
	Email          string
	HashedPassword string
}

type UserShort struct {
	ProfPic  string `json:"profpic"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type UserDetails struct {
	ProfPic  string `json:"profpic"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Mobile   int    `json:"mobile"`
}

/** ~Query Results */

/** Response */
type SessSuccess struct {
	Message string `json:"message"`
	UserId  int    `json:"userid"`
}

type UserShortResponse struct {
	Users []UserShort
}

type UserDetailResponse struct {
	User UserDetails
}

/** ~Response */
