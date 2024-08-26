# CVE Notifier

CVE Notifier is a Golang application that fetches the latest Common Vulnerabilities and Exposures (CVE) data and sends notifications to a specified Discord channel. This tool helps you stay updated on new security vulnerabilities that might affect your systems.

## Features

- Retrieves the latest CVE data from a reliable source.
- Sends real-time notifications to a Discord channel.
- Configurable notification settings.
- Easy to set up and deploy.

## Requirements

- [Go](https://golang.org/dl/) 1.16 or later
- A Discord Webhocks token

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/cve-notifier.git
    cd cve-notifier
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. Build the project:
    ```bash
    go build -o cve-notifier
    ```

## Configuration
- Add the webhock in the feed.go
