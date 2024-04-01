# Slackboat

Slackboat is a simple Slack bot written in Go that integrates with PostgreSQL for data storage.

## Features

- **Custom Commands:** Easily add custom commands to Slackboat to extend its functionality.
- **Persistent Storage:** Utilizes PostgreSQL to store data for seamless retrieval and persistence.
- **Easy Configuration:** Configure Slackboat quickly with environment variables.

## Prerequisites

- Go 1.16 or higher
- PostgreSQL
- Slack API token

## Installation

1. Clone this repository:

    ```bash
    git clone github.com/beingrohanpandit/slacboat-go
    ```

2. Navigate into the directory:

    ```bash
    cd slackboat
    ```

3. Install dependencies:

    ```bash
    go mod tidy
    ```

4. Set up environment variables

5. Build and run the application:

    ```bash
    go build -o slackboat .
    ./slackboat
    ```

## Usage

1. Invite Slackboat to your Slack workspace.
2. Interact with Slackboat using predefined commands or create your own custom commands.
3. Slackboat will respond to commands and interact with PostgreSQL for data storage.

## Contributing

Contributions are welcome! Please fork this repository and submit a pull request with your changes.
