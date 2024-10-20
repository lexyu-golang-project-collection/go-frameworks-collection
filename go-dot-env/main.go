package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	TEST    string
	API_KEY string
	DB_PASS string
)

func main() {
	mux := http.NewServeMux()

	fmt.Println("HOME = ", os.Getenv("HOME"))

	shell, ok := os.LookupEnv("SHELL")
	if !ok {
		fmt.Println("The env var SHELL is not set")
	} else {
		fmt.Println("SHELL = ", shell)
	}

	err := os.Setenv("TEST_NAME", "TESTNAME")
	if err != nil {
		fmt.Printf("Could not set the env var TEST_NAME")
	}
	fmt.Printf("TEST_NAME = %s\n", os.Getenv("TEST_NAME"))

	// GodotEnv

	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		fmt.Printf("Could not load .env file")
		os.Exit(1)
	}

	API_KEY = os.Getenv("API_KEY")
	DB_PASS = os.Getenv("DB_PASS")
	TEST = os.Getenv("TEST")

	fmt.Printf("API_KEY = %s\n", os.Getenv("API_KEY"))
	fmt.Printf("DB_PASS = %s\n", os.Getenv("DB_PASS"))
	fmt.Println("TEST = ", TEST)

	envMap, mapErr := godotenv.Read(".env")
	if mapErr != nil {
		fmt.Printf("Error loading .env into map[string]string\n")
		os.Exit(1)
	}

	for k, v := range envMap {
		fmt.Printf("Key = %s, Value = %s\n", k, v)
	}
	fmt.Printf("API_KEY = %s\n", envMap["API_KEY"])

	mux.HandleFunc("GET /overload", func(w http.ResponseWriter, r *http.Request) {
		errEnv := godotenv.Overload(".env")
		if errEnv != nil {
			fmt.Printf("Could not load .env file")
			os.Exit(1)
		}

		TEST = os.Getenv("TEST")
		API_KEY = os.Getenv("API_KEY")
		DB_PASS = os.Getenv("DB_PASS")

		str := "Success"

		w.Write([]byte(str))
	})

	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		str := fmt.Sprintf("Welcome to the home page!, %s, %s , %s", TEST, API_KEY, DB_PASS)

		w.Write([]byte(str))
	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
