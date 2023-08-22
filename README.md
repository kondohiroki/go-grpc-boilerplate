# go-grpc-boilerplate :rocket:

<p align="center">
<img src="https://user-images.githubusercontent.com/49369000/262280325-e9c5caa7-844f-46a5-af19-95172add3265.png"  width="500" />
</p>

## Getting Started
### Prerequisites
-  Go 1.21
-  Docker
-  sonar-scanner - for coverage test in local
   ```sh
   brew install sonar-scanner
   ```

### Installation
1. Clone the repo
   ```sh
   git clone https://github.com/kondohiroki/go-grpc-boilerplate
    ```
2. Install Go dependencies
    ```sh
    go mod tidy
    ```
3. Copy the default configuration file
    ```sh
    cp config/config.example.yaml config/config.yaml
    ```
4. Start the database
    ```sh
    docker compose up -d
    ```
5. Migrate database
    ```sh
    go run main.go migrate
    ```
6. Run the application
    ```sh
    # Run normally
    go run main.go serve:grpc-api

    # Run with hot reload
    air serve:grpc-api
    ```
7. Testing (optional)
    ```sh
    # Run unit-test
    make unit-test

    # Run api-test
    make api-test

    # Create sonar scret
    touch .sonar.secret
    echo "your-sonar-token" > .sonar.secret

    # Add secret to .sonar.secret
    # Get from sonar web
    ```
### Protof
1. Install protobuf compiler
    ```sh
    brew install protobuf
    ```

2. Install the protocol compiler plugins for Go
    ```sh
    $ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
    $ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    ```
3. Update your PATH so that the protoc compiler can find the plugins:
    ```sh
    export PATH="$PATH:$(go env GOPATH)/bin"
    ```
4. Compile the proto file to generate the gRPC code
    ```sh
    make pb
    ```
5. Make sure all dependencies are installed
    ```sh
    go mod tidy
    ```