package main

import (
	"flag"
	"fmt"
	"github.com/lanfeng6/Myblog/apicontext"
	"github.com/lanfeng6/Myblog/articleapp"
	"github.com/lanfeng6/Myblog/common"
	"github.com/lanfeng6/Myblog/configapp"
	"github.com/lanfeng6/Myblog/tagapp"
	"github.com/lanfeng6/Myblog/user"
	"github.com/leyle/ginbase/middleware"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var err error
	var port string
	var cfile string

	flag.StringVar(&cfile, "c", "", "-c /config/file/path")
	flag.StringVar(&port, "p", "", "-p 8200")

	flag.Parse()

	if cfile == "" {
		fmt.Println("缺少运行的配置文件, 使用 -c /config/file/path")
		os.Exit(1)
	}

	var conf = new(configapp.Config)
	conf.Load(cfile)

	if port != "" {
		conf.Server.Port = port
	}

	closeDsCh := make(chan struct{})

	// 初始化 context 内容
	// closeDataDs := make(chan struct{})
	ctx, err := apicontext.NewContext(closeDsCh, conf)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	go httpServer(ctx)
	go MakeMigrate(ctx.Ds)
	// 捕获信号，关闭数据库连接
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	close(closeDsCh)
	time.Sleep(100 * time.Millisecond)
}

func httpServer(ctx *apicontext.Context) {
	var err error
	engine := middleware.SetupGin()
	common.AuthOption = ctx
	apiRouter := engine.Group("") // yy means yunying
	//USER app
	user.UserRouter(ctx, apiRouter.Group(""))
	tagapp.TagRouter(ctx, apiRouter.Group(""))
	articleapp.ArticleRouter(ctx, apiRouter.Group(""))
	// start server
	addr := ctx.Conf.GetServer().GetServerAddr()
	err = engine.Run(addr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
