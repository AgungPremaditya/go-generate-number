package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

/**
 * initRedis is a function to initialize redis client
 * @return *redis.Client
 */
func initRedis() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		_ = fmt.Errorf("error loading .env file")
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPass := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Username: "default",
		Password: redisPass,
		DB:       0,
	})

	return rdb
}

/**
 * checkRedisConnection is a function to check redis connection
 * @param ctx context.Context
 * @param rdb *redis.Client
 */
func checkRedisConnection(ctx context.Context, rdb *redis.Client) {
	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Println("Error connecting to Redis:", err)
	} else {
		fmt.Println("Connected to Redis")
	}
}

func generateRandomNumber(digit int) string {
	if digit <= 0 {
		return "Digit can't be less than or equal to 0"
	}

	minVal := math.Pow(10, float64(digit-1))
	maxVal := math.Pow(10, float64(digit)) - 1

	randomNum := rand.Intn(int(maxVal-minVal)) + int(minVal)
	return strconv.Itoa(randomNum)
}

func generateUniqueNumber(ctx context.Context, digit int, rdb *redis.Client) (string, error) {
	// Generate a new random number
	randomNumber := generateRandomNumber(digit)
	key := fmt.Sprintf("random_number:%s", randomNumber)

	// Set the random number in the cache
	isSuccess, err := rdb.SetNX(ctx, key, randomNumber, 10*time.Minute).Result()
	if err != nil {
		return "", err
	}

	if isSuccess {
		return randomNumber, nil
	}

	return generateUniqueNumber(ctx, digit, rdb)
}

func handleGenerateNumber(req events.LambdaFunctionURLRequest) events.LambdaFunctionURLResponse {
	// set context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// init redis
	rdb := initRedis()

	// check redis connection
	checkRedisConnection(ctx, rdb)

	// set digits from query string
	digits := 1
	if req.RawQueryString != "" {
		// Parse the raw query string
		values, err := url.ParseQuery(req.RawQueryString)
		if err == nil {
			// Access the 'd' parameter
			if d, exists := values["d"]; exists && len(d) > 0 {
				if parsedDigits, err := strconv.Atoi(d[0]); err == nil {
					digits = parsedDigits
				}
			}
		}
	}

	// generate number
	rand, err := generateUniqueNumber(ctx, digits, rdb)
	if err != nil {
		fmt.Println("Error generating number:", err)
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       rand,
	}
}
func handler(req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	path := req.RequestContext.HTTP.Path

	var response events.LambdaFunctionURLResponse

	switch path {
	case "/api/generate":
		response = handleGenerateNumber(req)
	default:
		response = events.LambdaFunctionURLResponse{
			StatusCode: 404,
			Body:       "URL Not Found",
		}
	}

	return response, nil
}

func main() {
	lambda.Start(handler)
}
