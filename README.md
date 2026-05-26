# Ultimate Tic-Tac-Toe

[![Go Version](https://img.shields.io/github/go-mod/go-version/AshvinBambhaniya/tic-tac-toe?filename=backend%2Fgo.mod)](https://golang.org)
[![Nuxt Version](https://img.shields.io/badge/Nuxt-4.x-00DC82?logo=nuxt.js&logoColor=white)](https://nuxt.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A high-performance, real-time implementation of **Ultimate Tic-Tac-Toe**. Built with a Go backend and a Nuxt 4 frontend, featuring competitive multiplayer and a sophisticated AI engine.

## What is Ultimate Tic-Tac-Toe?

Ultimate Tic-Tac-Toe is a strategic board game composed of nine Tic-Tac-Toe boards arranged in a 3x3 grid. 

- **The Rule**: Each move made in a sub-grid determines which sub-grid the next player must play in.
- **The Goal**: Win three sub-grids in a row to win the entire game.
- **Complexity**: It is significantly more strategic than standard Tic-Tac-Toe, requiring long-term planning and "sending" opponents to unfavorable zones.

## Features

- **Real-Time Multiplayer**: Seamless PvP experience using WebSockets for instant move broadcasting.
- **AI/Bot Mode**: Practice against a computer opponent powered by a **Depth-Limited Minimax algorithm** with Alpha-Beta pruning.
    - **Easy**: Random play.
    - **Medium**: Depth-3 search (reactive blocking).
    - **Hard**: Depth-5+ search (strategic meta-board evaluation).
- **Live Matchmaking**: Connect to the lobby and find opponents automatically.
- **Player Profiles**: Track your wins, losses, draws, and win rates over time.
- **Match History & Review**: Revisit your past matches and analyze the final board state.
- **Secure Auth**: JWT-based authentication for persistent player identities.

## Tech Stack

### Backend
- **Language**: Go (v1.26+)
- **Framework**: [Fiber v2](https://gofiber.io/)
- **Real-time**: Custom WebSocket Hub implementation.
- **Database**: MySQL / PostgreSQL with [goqu](https://github.com/doug-martin/goqu) SQL builder.
- **Migrations**: [sql-migrate](https://github.com/rubenv/sql-migrate).

### Frontend
- **Framework**: [Nuxt 4](https://nuxt.com/) (Vue 3)
- **Styling**: [Tailwind CSS](https://tailwindcss.com/) with Glassmorphism UI.
- **State Management**: Vue Composition API & Custom Composables.

## Getting Started

### Prerequisites
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)
- [Node.js](https://nodejs.org/) (for local frontend development)
- [Go](https://golang.org/) (for local backend development)

### Quick Start (Docker)

1. **Clone the repository**:
   ```bash
   git clone https://github.com/AshvinBambhaniya/tic-tac-toe.git
   cd tic-tac-toe
   ```

2. **Configure Environment**:
   ```bash
   cp backend/.env.example backend/.env
   # Edit backend/.env if necessary
   ```

3. **Spin up the stack**:
   ```bash
   docker-compose up -d --build
   ```

The application will be available at:
- **Frontend**: `http://localhost:3000`
- **Backend API**: `http://localhost:8000`

## Project Structure

```text
├── backend/                # Go source code
│   ├── controllers/        # API endpoints
│   ├── services/           # Business logic & AI Engine
│   ├── models/             # Database schemas & queries
│   ├── database/           # Migrations
│   └── pkg/websocket/      # Real-time communication hub
├── frontend/               # Nuxt application
│   ├── app/pages/          # Routing & Views
│   ├── app/components/     # UI Components
│   └── app/composables/    # Shared state & logic
└── docker-compose.yaml     # Orchestration
```

## AI Implementation

The "Hard" bot utilizes a recursive Minimax search. The heuristic function evaluates:
1. **Sub-grid Wins**: 100 points for winning a sub-grid.
2. **Sub-grid Setup**: 10 points for 2-in-a-row with an empty third cell.
3. **Meta-board Position**: Strategic weighting for the center and corners of the 9x9 grid.
4. **Active Zone Manipulation**: Deducting points for moves that send the human player to a sub-grid they are about to win.

## Contributing

Contributions are welcome! Please follow these steps:
1. Fork the Project.
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`).
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the Branch (`git push origin feature/AmazingFeature`).
5. Open a Pull Request.

## License

Distributed under the MIT License. See `LICENSE` for more information.

---

**Developed by Ashvin Bambhaniya**
