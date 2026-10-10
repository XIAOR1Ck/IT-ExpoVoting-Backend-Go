# 11TH GCES IT EXPO Project Voting System Backend
### This is a simple voting system written in Go to choose people's choice awards for the project display.
## How it works
- You Register for voting by entering your Information like Name, Email, address, affiliation
- You'll be registered for voting and a QR code will be generated with the voting frontend link which can be used to cast vote from anywhere.
- In the frontend, you can choose a single project to vote for people's choce award.

## Components Used
- [Chi Router](github.com/go-chi/chi/v5) + Go's native `net/http` for request handling
- `go-redis` + Redis for low-latency data storage.
- [go-qrcode](https://github.com/skip2/go-qrcode) for QR Code Generation.
