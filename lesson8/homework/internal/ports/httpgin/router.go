package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/app"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.POST("/ads", createAd(a))                    // Метод для создания объявления (ad)
	r.PUT("/ads/:ad_id/status", changeAdStatus(a)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	r.PUT("/ads/:ad_id", updateAd(a))              // Метод для обновления текста(Text) или заголовка(Title) объявления
	r.GET("/ads/:ad_id", getAd(a))                 // Метод для получения объявления по ID

	r.GET("/ads", listAds(a)) // Метод для получения списка объявлений с фильтрами (по published, userID, date, title)

	r.POST("/users", createUser(a))         // Метод для создания пользователя (user)
	r.GET("/users/:user_id", getUser(a))    // Метод для получения пользователя по ID
	r.PUT("/users/:user_id", updateUser(a)) // Метод для обновления имени(Nickname) или почты(Email) пользователя
}
