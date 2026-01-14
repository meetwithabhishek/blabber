# Blabber 🐴

Blabber is a simple, terminal-based chat application written in Go. It uses a client-server architecture with WebSockets for real-time communication, allowing multiple users to chat with each other directly from their terminals.

![Blabber Demo](https://raw.githubusercontent.com/meetwithabhishek/blabber/main/demo.gif)

## Features

*   **Real-time Chat**: Send and receive messages instantly using a WebSocket-based backend.
*   **Terminal UI**: A clean and colorful terminal user interface (TUI) built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).
*   **Persistent History**: The server saves all messages in an SQLite database, and clients receive the message history upon connection.
*   **Colored Usernames**: Each user is assigned a unique color to make conversations easy to follow.
*   **Easy Setup**: Run the server from the source and install the client with a single command. The client will prompt you for configuration on its first run.

## How It Works

Blabber consists of two main components:

1.  **The Server**: A Go application that listens for incoming WebSocket connections. It receives messages from clients, stores them in an SQLite database, and broadcasts them to all connected clients. It also serves the entire message history over a REST endpoint for new clients.
2.  **The Client**: A terminal application that connects to the server. On first launch, it prompts the user to enter a username and the server's address. It then displays the chat history and allows the user to send new messages.

## Getting Started

### Prerequisites

*   Go (version 1.23 or newer) is required.

### Running the Server

To run the server, you'll need to clone the repository.

1.  **Clone the repository**:
    ```sh
    git clone https://github.com/meetwithabhishek/blabber.git
    cd blabber
    ```

2.  **Start the server**:
    ```sh
    go run ./server
    ```
    The server will start on `localhost:8080` by default.

### Running the Client

#### Option 1: Install with `go install` (Recommended)

You can install and run the client with a single command, without cloning the repository. This is the easiest way to get started.

1.  **Install the client**:
    ```sh
    go install github.com/meetwithabhishek/blabber@latest
    ```

2.  **Run the client**:
    ```sh
    blabber
    ```
    *   On the first run, you will be prompted to enter a **username** and the **server address**.
    *   For a local server, you can enter `localhost` as the server address.

#### Option 2: Run from Source

If you have already cloned the repository to run the server, you can also run the client from the source code in a new terminal window.

```sh
go run .
```

You can open multiple terminal windows and run the client in each to simulate a multi-user chat.

## Controls

*   **Type** to compose your message.
*   **Enter** to send the message.
*   **Ctrl+C** or **q** to quit the application.

