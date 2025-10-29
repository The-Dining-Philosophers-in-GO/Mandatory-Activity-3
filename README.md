# 💬 Chit Chat — Distributed Chat System

Chit Chat is a distributed chat service that allows participants to **join**, **exchange messages**, and **leave** the conversation at any time.  
The system is implemented in **Go** using **gRPC** and **Protocol Buffers**, and demonstrates key concepts of **distributed systems** such as communication, coordination, concurrency, and logical time (Lamport timestamps).

---

## 🏗️ System Overview

Chit Chat consists of:
- **One server process** that manages message broadcasting and participant coordination.
- **Multiple client processes** that connect to the server to send and receive messages in real time.

Each client communicates with the server through gRPC.  
Every message (including join/leave notifications) is timestamped using **Lamport logical clocks** to preserve event ordering in the absence of a global clock.

---

## ⚙️ Technical Design Summary

- **Language:** Go  
- **Communication:** gRPC with Protocol Buffers (`.proto` definitions under `/grpc`)  
- **Concurrency:** Managed using goroutines and channels  
- **Logging:** Implemented using Go’s built-in `log` library  
- **Architecture:** Client–Server  
- **Message Broadcasts:** Include both content and logical timestamp  
- **Join/Leave Messages:**  
  - “Participant X joined Chit Chat at logical time L.”  
  - “Participant X left Chit Chat at logical time L.”

---

## 🏃‍♂️ How to run program

For the program to work server have to be set up first.

Open a terminal and make sure that you're in the server folder :
  cd server

Then you can start it with :
  go run server.go

For a client to join the server, open a new terminal make sure you're in the folder :
  cd client

Then you can start it with :
  go run client.go -id YOURID

You can type a message, by just typing in the terminal.

If you want to leave the server type
  /leave

## 📦 Repository Structure

project-root/  
├── client/ # contains the client code  
├── grpc/ # contains .proto file  
├── server/ # contains the server code  
└── readme.md # this file
