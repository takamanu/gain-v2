package routes

import (
	"gain-v2/configs"
	// "gain-v2/features/fashions"

	"gain-v2/features/logging"
	"gain-v2/features/users"

	// "gain-v2/features/vouchers"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func RouteUser(e *echo.Group, uh users.UserHandlerInterface, cfg configs.ProgrammingConfig) {
	e.POST("/admin/add-user", uh.AddUser())
	e.POST("/admin/login", uh.LoginAdmin())

	e.POST("/register", uh.Register())
	// e.POST("/admin/register", uh.Register())
	// e.POST("/login", uh.LoginCustomer())
	e.POST("/login", uh.Login())
	e.POST("/forgot-password", uh.ForgotPassword())
	// e.POST("/forget-password/verify", uh.ForgetPasswordVerify())
	e.POST("/reset-password", uh.ResetPassword())
	// e.POST("/refresh-token", uh.RefreshToken(), echojwt.JWT([]byte(cfg.Secret)))
	e.PUT("/admin/update", uh.UpdateProfile(), echojwt.JWT([]byte(cfg.Secret)))
	e.GET("/user/profile", uh.GetProfile(), echojwt.JWT([]byte(cfg.Secret)))
}

func RouteLogging(e *echo.Group, lh logging.LoggingHandlerInterface, cfg configs.ProgrammingConfig) {

	e.POST("/logging", lh.AddLog())
	e.GET("/logging", lh.ViewLog())
	e.GET("/logging/:log_id", lh.ViewOneLog())
	e.DELETE("/logging/:log_id", lh.DeleteLog())

}
