package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "task is a CLI for managing your TODOs.",
	Long:  `task is a CLI task manager that helps you organize your work.`,
}

var addCmd = &cobra.Command{
	Use:   "add [taskName]",
	Short: "Add a new task",
	Long:  "Add a new task to your TODO list",
	Args:  cobra.MinimumNArgs(1),
	Run:   addTask,
}

var doCmd = &cobra.Command{
	Use:   "do [taskNumber]",
	Short: "Mark a task",
	Long:  "Mark a task on your TODO list as complete",
	Args:  cobra.ExactArgs(1),
	Run:   doTask,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your task",
	Long:  "List all of your incomplete tasks",
	Run:   listTasks,
}

var rmCmd = &cobra.Command{
	Use:   "rm [taskNumber]",
	Short: "Delete a task",
	Long:  "Delete a task from your TODO list",
	Args:  cobra.ExactArgs(1),
	Run:   deleteTask,
}

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "List all of your completed task",
	Long:  "List all of your completed tasks",
	Run:   completeTasks,
}

func init() {
	rootCmd.AddCommand(addCmd, doCmd, listCmd, rmCmd, completeCmd)
}

func openDatabase() *bolt.DB {
	db, err := bolt.Open("tasks.db", 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func addTask(cmd *cobra.Command, args []string) {
	taskDescription := strings.Join(args, " ")

	db := openDatabase()

	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Put([]byte(taskDescription), []byte("0"))

		if err != nil {
			return err
		}

		return nil
	})
}

func doTask(cmd *cobra.Command, args []string) {
	id, _ := strconv.Atoi(args[0])

	db := openDatabase()

	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))

		if b == nil {
			return errors.New("invalid operation!!! no task exists")
		}

		c := b.Cursor()

		cnt := 0

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(v) == "0" {
				cnt++
			}

			if cnt == id {
				b.Put([]byte(k), []byte("1"))
				fmt.Printf("You have completed the '%s' task.", k)
				break
			}
		}

		return nil
	})
}

func listTasks(cmd *cobra.Command, args []string) {
	db := openDatabase()

	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))

		if b == nil {
			return errors.New("take a break!!! no task is pending right now")
		}

		c := b.Cursor()

		cnt := 1

		fmt.Println("You have the following tasks:")
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(v) == "0" {
				fmt.Printf("%d. %s\n", cnt, k)
				cnt++
			}
		}

		return nil
	})
}

func deleteTask(cmd *cobra.Command, args []string) {
	id, _ := strconv.Atoi(args[0])

	db := openDatabase()

	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))

		if b == nil {
			return errors.New("invalid operation!!! no tasks found in current database")
		}

		c := b.Cursor()

		cnt := 0

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(v) == "0" {
				cnt++
			}

			if cnt == id {
				b.Delete([]byte(k))
				fmt.Printf("You have deleted the '%s' task.", k)
				break
			}
		}

		return nil
	})
}

func completeTasks(cmd *cobra.Command, args []string) {
	db := openDatabase()

	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))

		if b == nil {
			return errors.New("invalid operation!!! no tasks found in current database")
		}

		c := b.Cursor()

		fmt.Println("You have finished the following tasks today:")
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(v) == "1" {
				fmt.Printf("- %s\n", k)
			}
		}

		return nil
	})
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
