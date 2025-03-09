package cmd

import (
	"context"
	"fmt"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
	"github.com/saleh-ghazimoradi/Gophergram/utils"
	"github.com/spf13/cobra"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seeding DB",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("seed called")
		db, err := utils.PostConnection()
		if err != nil {
			sLogger.SLogger.Error(err.Error())
		}

		userRepository := repository.NewUserRepository(db, db)
		postRepository := repository.NewPostRepository(db, db)
		commentRepository := repository.NewCommentRepository(db, db)
		seed := repository.NewSeederService(userRepository, postRepository, commentRepository)

		if err := seed.Seed(context.Background(), db); err != nil {
			sLogger.SLogger.Error(err.Error())
		}

	},
}

func init() {
	rootCmd.AddCommand(seedCmd)

}
