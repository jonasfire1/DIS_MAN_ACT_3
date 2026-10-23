# Chit Chat

Chit Chat is a distributed chat service written in Go. Clients communicate with a central server using gRPC.

## Requirements

- Go
- Four terminals: one server and three clients

## Run the program

Open a terminal in the project root and download the dependencies:

```powershell
terminal 1
go run ./server

terminal 2
go run ./client Alice

terminal 3 
go run ./client Bob

terminal 4
go run ./client Charlie

Type this into Alices terminal:
Hello everyone

Press Ctrl+C in Charlies terminal. Alice and Bob should receive something like:
[Logical time 6] Participant Charlie left Chit Chat at logical time 6

alternativly, you can exit by simply writing leave