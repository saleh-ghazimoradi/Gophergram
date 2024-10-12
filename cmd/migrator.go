package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/logger"
	"github.com/saleh-ghazimoradi/Gophergram/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func getSourceURL(path string) string {
	if path != "" {
		return path
	}
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalln("Error: can not find project root path", err)
	}

	var sourceURL string

	pattern := regexp.MustCompile(`^.*\.(down|up)\.sql$`)

	err = filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		match := pattern.MatchString(info.Name())
		if match {
			sourceURL = filepath.Dir(path)
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	if sourceURL == "" {
		log.Println("Can not find the SourceURL of database schemes directory.")
		log.Fatalln("Please ensure you are in the root of the project when run or build the project.")
	}
	return sourceURL
}

var migratorCmd = &cobra.Command{
	Use:   "migrator",
	Short: "Manages your database migrations",
	Long: `The migrator command allows you to manage your database migrations.
You can apply all up migrations using 'migrator up' or apply all down migrations using 'migrator down'.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please use 'migrator up' to apply all up migrations or 'migrator down' to apply all down migrations.")
	},
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all up migrations",
	Long: `The 'migrator up' command applies all up migrations to your database.
	This will update your database schema to the latest version.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("Applying all up migrations...")
		path := cmd.Flag("path").Value.String()
		repository.MigrateUp("file://"+getSourceURL(path),
			utils.PostgresUrl(
				config.AppConfig.Database.Postgresql.Host,
				config.AppConfig.Database.Postgresql.Port,
				config.AppConfig.Database.Postgresql.User,
				config.AppConfig.Database.Postgresql.Password,
				config.AppConfig.Database.Postgresql.Database,
				config.AppConfig.Database.Postgresql.SSLMode,
			),
		)
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Apply all down migrations",
	Long: `The 'migrator down' command applies all down migrations to your database.
	This will revert your database schema to the previous version.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("Applying all down migrations...")
		path := cmd.Flag("path").Value.String()
		repository.MigrateDown("file://"+getSourceURL(path),
			utils.PostgresUrl(
				config.AppConfig.Database.Postgresql.Host,
				config.AppConfig.Database.Postgresql.Port,
				config.AppConfig.Database.Postgresql.User,
				config.AppConfig.Database.Postgresql.Password,
				config.AppConfig.Database.Postgresql.Database,
				config.AppConfig.Database.Postgresql.SSLMode,
			),
		)
	},
}

func init() {
	upCmd.PersistentFlags().StringP("path", "p", "", "path to database schemes directory")
	downCmd.PersistentFlags().StringP("path", "p", "", "path to database schemes directory")
	migratorCmd.AddCommand(upCmd)
	migratorCmd.AddCommand(downCmd)
}
