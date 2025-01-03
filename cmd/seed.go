/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/utils"
	"log"

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
			log.Fatal(err)
		}
		postRepository := repository.NewPostRepository(db, db)
		userRepository := repository.NewUserRepository(db, db)
		commentRepository := repository.NewCommentRepository(db, db)
		postService := service.NewPostService(postRepository, db)
		userService := service.NewUserService(userRepository, db)
		commentService := service.NewCommentService(commentRepository)
		seed := service.NewSeederService(userService, postService, commentService)

		if err := seed.Seed(context.Background(), db); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)

}
