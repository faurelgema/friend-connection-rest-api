package controller

import (
	"net/http"

	friendshipController "friend_connection_rest_api/controller/friendship"
	userController "friend_connection_rest_api/controller/user"
	migration "friend_connection_rest_api/migrations"
	friendshipService "friend_connection_rest_api/services/friendship"
	userService "friend_connection_rest_api/services/user"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	//_ "github.com/swaggo/gin-swagger/example/basic/docs"
	"gorm.io/gorm"
)

//Setup Manager, Migration and Routes
func Setup(db *gorm.DB) http.Handler {
	friendshipService := friendshipService.NewFriendshipManager(db)
	userService := userService.NewUserManager(db)
	migration.InitMigration(db)
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	//url := ginSwagger.URL("http://localhost:3000/docs/swagger.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/list-users", func(c *gin.Context) {
		userController.GetListUsersController(c, userService)
	})

	r.POST("/create-user", func(c *gin.Context) {
		userController.CreateNewUserController(c, userService)
	})

	r.POST("/add-friends", func(c *gin.Context) {
		friendshipController.MakeFriendController(c, friendshipService)
	})

	r.POST("/get-list-friends", func(c *gin.Context) {
		friendshipController.GetFriendsListController(c, friendshipService)
	})

	r.POST("/get-mutual-list-friends", func(c *gin.Context) {
		friendshipController.GetMutualFriendsController(c, friendshipService)
	})

	r.POST("/subscribe", func(c *gin.Context) {
		friendshipController.SubscribeController(c, friendshipService)
	})

	r.POST("/block", func(c *gin.Context) {
		friendshipController.BlockController(c, friendshipService)
	})

	r.POST("/get-list-users-receive-update", func(c *gin.Context) {
		friendshipController.GetUsersReceiveUpdateController(c, friendshipService)
	})
	return r
}
