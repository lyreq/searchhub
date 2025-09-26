package middleware

import (
	"lytemp/pkg"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		// Token boş mu kontrol et
		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"message": "Yetkilendirme başarısız: Token bulunamadı",
			})
		}

		// Bearer prefix'ini kaldır
		tokenString := strings.TrimPrefix(token, "Bearer ")

		// Token'ı doğrula
		claim, err := pkg.ValidateToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"message": "Yetkilendirme başarısız: " + err.Error(),
			})
		}

		// Kullanıcı bilgilerini context'e ekle
		c.Set("user", claim.User)
		c.Set("userid", claim.User.ID)

		return next(c)
	}
}
