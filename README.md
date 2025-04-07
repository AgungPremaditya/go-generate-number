# Random Number Generator Lambda

This project is an AWS Lambda function that generates unique random numbers and stores them temporarily in a Redis cache to ensure uniqueness. The function is triggered via an HTTP request and can generate random numbers with a specified number of digits.

## Features

- Generates random numbers with a specified number of digits.
- Ensures uniqueness of generated numbers using Redis.
- Provides a simple HTTP API to request random numbers.

## Prerequisites

- AWS account with permissions to deploy Lambda functions.
- Redis server for caching generated numbers.
- Go environment for building the project.

## Environment Variables

The application requires the following environment variables to be set:

- `REDIS_HOST`: The host address of the Redis server.
- `REDIS_PASSWORD`: The password for the Redis server.

These can be set in a `.env` file in the root of the project.

## Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yourusername/random-number-generator-lambda.git
   cd random-number-generator-lambda
   ```

2. **Install dependencies:**

   Ensure you have Go installed, then run:

   ```bash
   go mod tidy
   ```

3. **Set up environment variables:**

   Create a `.env` file in the root directory and add your Redis configuration:

   ```plaintext
   REDIS_HOST=your_redis_host
   REDIS_PASSWORD=your_redis_password
   ```

4. **Deploy to AWS Lambda:**

   Use the AWS CLI or AWS Management Console to deploy the function. Ensure the Lambda function has access to the Redis server.

## Usage

- **Endpoint:** `/api/generate`
- **Query Parameter:** `d` (optional) - Number of digits for the random number. Defaults to 1 if not provided.

### Example Request

```http
GET /api/generate?d=5
```

### Example Response

```json
{
  "statusCode": 200,
  "body": "12345"
}
```

## License

This project is licensed under the MIT License.
