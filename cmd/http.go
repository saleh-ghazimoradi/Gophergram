package cmd

import (
	"expvar"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/logger"
	"github.com/saleh-ghazimoradi/Gophergram/utils"
	"github.com/spf13/cobra"
	"runtime"
	"time"
)

func init() {
	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "launching the http rest listen server",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Logger.Infow("server has started", "addr", config.AppConfig.General.Listen, "env", config.AppConfig.Env.Env)

		// TODO: Make it more efficient
		cfg := utils.PostgresConfig{
			Host:         config.AppConfig.Database.Postgresql.Host,
			Port:         config.AppConfig.Database.Postgresql.Port,
			User:         config.AppConfig.Database.Postgresql.User,
			Password:     config.AppConfig.Database.Postgresql.Password,
			Database:     config.AppConfig.Database.Postgresql.Database,
			SSLMode:      config.AppConfig.Database.Postgresql.SSLMode,
			MaxOpenConns: config.AppConfig.Database.Postgresql.MaxOpenConns,
			MaxIdleConns: config.AppConfig.Database.Postgresql.MaxIdleConns,
			MaxIdleTime:  config.AppConfig.Database.Postgresql.MaxIdleTime,
			Timeout:      config.AppConfig.Database.Postgresql.Timeout,
		}

		db, err := utils.PostgresConnection(cfg)
		if err != nil {
			logger.Logger.Fatal(err)
		}
		defer db.Close()

		redis, err := utils.RedisConnection(config.AppConfig.Database.Redis.Addr, config.AppConfig.Database.Redis.PW, config.AppConfig.Database.Redis.DB)
		if err != nil {
			logger.Logger.Fatal(err)
		}

		logger.Logger.Info("Postgresql connection pool established")
		logger.Logger.Info("Redis connection pool established")

		expvar.NewString("version").Set(gateway.Version)
		expvar.Publish("goroutines", expvar.Func(func() any {
			return runtime.NumGoroutine()
		}))

		expvar.Publish("database", expvar.Func(func() any {
			return db.Stats()
		}))

		expvar.Publish("timestamp", expvar.Func(func() any {
			return time.Now().Unix()
		}))

		/*-------------------repo---------------------*/
		cacheDB := repository.NewCacheRepo(redis)
		postDB := repository.NewPostRepository(db)
		commentDB := repository.NewCommentRepository(db)
		userDB := repository.NewUserRepository(db)
		followDB := repository.NewFollowRepository(db)
		roleDB := repository.NewRoleRepository(db)
		rateLimitDB := repository.NewRateLimitRepo(redis)

		/*-------------------service---------------------*/
		postService := service.NewPostService(postDB, commentDB, db)
		commentService := service.NewCommentService(commentDB, db)
		userService := service.NewServiceUser(userDB, cacheDB, db)
		followService := service.NewFollowService(followDB, db)
		mailerService := service.NewSendGridMailer(config.AppConfig.General.Mail.SendGrid.ApiKey, config.AppConfig.General.Mail.SendGrid.FromEmail)
		jwtAuthentication := service.NewJWTAuthenticator(config.AppConfig.General.Auth.Token.Secret, config.AppConfig.General.Auth.Token.TokenHost, config.AppConfig.General.Auth.Token.TokenHost)
		roleService := service.NewRoleService(roleDB)
		rateLimitService := service.NewRateLimitService(rateLimitDB)
		/*-------------------handler----------------------*/
		postHandler := gateway.NewPostHandler(postService, commentService)
		userHandler := gateway.NewUserHandler(userService, followService)
		feedHandler := gateway.NewFeedHandler(postService)
		authHandler := gateway.NewAuth(userService, mailerService, jwtAuthentication)
		authMiddleware := gateway.NewMiddleware(userService, jwtAuthentication, postService, roleService, rateLimitService)

		routeHandlers := gateway.Handlers{
			CreatePostHandler:      postHandler.CreatePost,
			GetPostHandler:         postHandler.GetPost,
			DeletePostHandler:      postHandler.DeletePost,
			UpdatePostHandler:      postHandler.UpdatePost,
			GetUserHandler:         userHandler.GetUserByID,
			FollowUserHandler:      userHandler.FollowUserHandler,
			UnfollowUserHandler:    userHandler.UnfollowUserHandler,
			GetUserFeedHandler:     feedHandler.GetUserFeedHandler,
			RegisterUserHandler:    authHandler.RegisterUserHandler,
			ActivateUserHandler:    userHandler.ActivateUserHandler,
			CreateTokenHandler:     authHandler.CreateTokenHandler,
			PostsContextMiddleware: postHandler.PostsContextMiddleware,
			AuthTokenMiddleware:    authMiddleware.AuthToken,
			CheckPostOwnership:     authMiddleware.CheckPostOwnership,
			RateLimitMiddleware:    authMiddleware.RateLimitMiddleware,
		}

		if err := gateway.Server(gateway.Routes(routeHandlers, rateLimitService)); err != nil {
			logger.Logger.Fatal(err)
		}
	},
}
