package cmd

import (
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/logger"
	"github.com/saleh-ghazimoradi/Gophergram/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(seedCmd)
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "launching the database seeding functionality",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("The database seeding functionality is starting", config.AppConfig.General.Listen)
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

		postDB := repository.NewPostRepository(db)
		commentDB := repository.NewCommentRepository(db)
		userDB := repository.NewUserRepository(db)
		postService := service.NewPostService(postDB, commentDB, db)
		commentService := service.NewCommentService(commentDB, db)
		userService := service.NewServiceUser(userDB, nil, db)
		seedService := service.NewSeederService(postService, commentService, userService)

		if err := seedService.SeedDatabase(cmd.Context(), db); err != nil {
			logger.Logger.Fatal("Failed to seed the database:", err)
		}
		logger.Logger.Info("Database seeding completed successfully")
	},
}
