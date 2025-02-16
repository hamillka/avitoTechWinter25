package dto

// AuthRequestDto model info
// @Description Информация о пользователе для входа или регистрации
type AuthRequestDto struct {
	Username string `json:"username"` // Имя пользователя
	Password string `json:"password"` // Пароль
}

// AuthResponseDto model info
// @Description JWT-токен пользователя при входе
type AuthResponseDto struct {
	Token string `json:"token"` // JWT-токен
}

const (
	AvitoShopName = "AvitoShop"
)
