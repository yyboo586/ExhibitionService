package main

import (
	"ExhibitionService/internal/drivenadapter"
	"ExhibitionService/internal/logics"
	"ExhibitionService/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func main() {
	g.Log().SetFlags(glog.F_ASYNC | glog.F_TIME_DATE | glog.F_TIME_TIME | glog.F_FILE_SHORT)
	s := g.Server()
	s.SetOpenApiPath("/api.json")
	s.SetSwaggerPath("/swagger")

	drivenFileEngine := drivenadapter.NewFileEngine()

	logicsFile := logics.NewFile()
	logicsCompany := logics.NewCompany()

	fileService := service.NewFileService(logicsFile, drivenFileEngine)
	companyService := service.NewCompanyService(logicsCompany, logicsFile)

	s.Group("/api/v1/exhibition-service", func(group *ghttp.RouterGroup) {
		group.Middleware(CORS)
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		// 展会服务相关接口
		group.Bind(
			fileService,
			companyService,
		)
	})

	s.Run()
}

func CORS(r *ghttp.Request) {
	corsOptions := r.Response.DefaultCORSOptions()
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}
