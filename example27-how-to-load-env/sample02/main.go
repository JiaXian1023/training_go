package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)
//autoload >init>godotenv.Load()
func main() {
	s3Bucket := os.Getenv("S3_BUCKET")
	secretKey := os.Getenv("SECRET_KEY")

	fmt.Println(s3Bucket)
	fmt.Println(secretKey)
}
