package cmd

import (
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/logger"
	"github.com/saleh-ghazimoradi/Gophergram/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "launching the http rest listen server",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("http rest server is starting", config.AppConfig.General.Listen)

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

		logger.Logger.Info("database connection pool established")

		/*-------------------repo---------------------*/
		postDB := repository.NewPostRepository(db)
		commentDB := repository.NewCommentRepository(db)
		userDB := repository.NewUserRepository(db)

		/*-------------------service---------------------*/
		postService := service.NewPostService(postDB, commentDB)
		commentService := service.NewCommentService(commentDB)
		userService := service.NewServiceUser(userDB)
		/*-------------------handler----------------------*/
		postHandler := gateway.NewPostHandler(postService, commentService)
		userHandler := gateway.NewUserHandler(userService)

		routeHandlers := gateway.Handlers{
			CreatePostHandler:      postHandler.CreatePost,
			GetPostHandler:         postHandler.GetPost,
			DeletePostHandler:      postHandler.DeletePost,
			UpdatePostHandler:      postHandler.UpdatePost,
			GetUserHandler:         userHandler.GetUserByID,
			PostsContextMiddleware: postHandler.PostsContextMiddleware,
		}

		if err := gateway.Server(gateway.Routes(routeHandlers)); err != nil {
			logger.Logger.Fatal(err)
		}
	},
}
