package handler

type AuthHandler struct {
	users map[string]string
}

func NewAuthHandler(users map[string]string) *AuthHandler {
	return &AuthHandler{users: users}
}

// func (a *AuthHandler) registerRoutes(g *gin.RouterGroup) {
// 	g.GET("/secret", a.handleSectret, gin.BasicAuth(a.users))
// }

// func (a *AuthHandler) handleSectret(c *gin.Context) {
// 	// /admin/secrets endpoint
// 	// hit "localhost:8080/admin/secrets

// 	// get user, it was set by the BasicAuth middleware
// 	user := c.MustGet(gin.AuthUserKey).(string)
// 	if secret, ok := a.users[user]; ok {
// 		c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
// 	}

// }
