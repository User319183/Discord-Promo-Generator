# 🚀 Discord Promo Generator 🚀

This repository contains a Go application that automates the creation of Discord promos. It's designed to be efficient and easy to use, making it perfect for developers and Discord enthusiasts alike.

## 📚 What it does 📚

The application generates Discord promos by making HTTP requests to the Discord API. It uses goroutines to handle multiple requests concurrently, maximizing efficiency and speed. The application also uses a proxy for the requests, ensuring that they are not blocked due to too many requests from a single IP.

## 🛠️ How it works 🛠️

The application is built in Go and uses the `net/http` package to make the HTTP requests. It also uses the `github.com/google/uuid` package to generate unique identifiers for each request, and the `github.com/fatih/color` package to colorize console output.

The application starts by initializing a new `App` struct, which contains the HTTP client, headers for the requests, and other necessary data. It then enters a loop where it continuously creates new promos. Each promo creation is done in a separate goroutine, allowing for concurrent execution.

## 🚀 How to run 🚀

To run the application, you need to have Go installed on your machine. Once you have Go installed, you can run the application by executing the following command in the terminal:

```bash
go run main.go
```

## 📝 Note 📝

Please use this application responsibly. Do not use it to spam or abuse the Discord API. Always respect the terms of service of any API you are using.

## 📚 License 📚

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments 🙏

Thanks to the Go community for the great resources and libraries that made this project possible.
