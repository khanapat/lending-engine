# Lending-Engine (Backend API)

[![Github Actions Status](https://github.com/gbrlsnchs/jwt/workflows/Linux,%20macOS%20and%20Windows/badge.svg)](https://github.com/gbrlsnchs/jwt/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/gbrlsnchs/jwt)](https://goreportcard.com/report/github.com/gbrlsnchs/jwt)
[![GoDoc](https://godoc.org/github.com/gbrlsnchs/jwt?status.svg)](https://pkg.go.dev/github.com/gbrlsnchs/jwt/v3)
[![Version compatibility with Go 1.11 onward using modules](https://img.shields.io/badge/compatible%20with-go1.11+-5272b4.svg)](https://github.com/gbrlsnchs/jwt#installing)

<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

## About The Project

This backend api connect with lending service web for ICFIN project. Using postgres, redis, websocket and mail api.

[![Go-logo-aqua.png](https://i.postimg.cc/FHJggSg7/Go-logo-aqua.png)](https://postimg.cc/75w2qfRk)

### Built With

This project use these components to run service.

* [Go](https://golang.org) 
* [Redis](https://redis.io/)
* [Postgres](https://www.postgresql.org/)
* [Mail-API](https://github.com/khanapat/mail-api)
* [Bitkub-Websocket](https://github.com/khanapat/bitkub-websocket)

<!-- GETTING STARTED -->

## Getting Started

These are steps to set up this project locally. To get a local running, please follow there simple steps.

### Prerequisites

* Go
```bash
brew install go
```

### Installation

1. Clone the repo
```bash
git clone https://github.com/khanapat/lending-engine.git
```

### Database

1. Running postgres docker
```bash
make create-postgres
```

2. Create database & table
SQL script in [SQL Script](https://github.com/khanapat/lending-engine/blob/master/lending.sql)

3. Stopping postgres docker
``` bash
make delete-postgres
```

### Redis

1. Running redis docker
```bash
make create-redis
```

2. Stopping redis
```bash
make delete-redis
```

<!-- USAGE EXAMPLES -->
## Usage

There are swagger api for testing. It can use after running app. The URL is `http://localhost:9090/swagger/index.html`. Default username is `admin` and Default password is `password`.

* Running App
```bash
make run
```

## Contact

My Contact Email - k.apiwattanawong@gmail.com
