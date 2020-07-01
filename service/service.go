package service

import (
	"github.com/labstack/echo"
)

var S *Service

//服务发现注册
type Service struct {
	Log echo.Logger
}
