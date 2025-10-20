============================================================
💬 CHIT CHAT — DISTRIBUTED CHAT SYSTEM
============================================================

In this assignment, you will design and implement Chit Chat,
a distributed chat service where participants can join, exchange
messages, and leave the conversation at any time.

Chit Chat is a lively playground for exploring the essence of
distributed systems — communication, coordination, and the
ordering of events in a world without a single shared clock.


============================================================
🧩 SYSTEM SPECIFICATION
============================================================

S1. Communication
-----------------
- Chit Chat is a distributed service that enables clients to exchange
  chat messages using gRPC.
- You must design the gRPC API, including:
  * All service methods
  * All message types


S2. Topology
-------------
- The system follows a distributed topology consisting of:
  * One service process (the server)
  * Multiple client processes
- Each client runs independently and communicates with the service
  via gRPC.
- Minimal configuration:
  * 1 service instance
  * At least 3 concurrently active clients


S3. Publishing Messages
------------------------
- Each participant can publish a valid chat message at any time.
- A valid message is:
  * UTF-8 encoded
  * Maximum 128 characters
- Publishing is done via a gRPC call to the Chit Chat service.


S4. Broadcasting
-----------------
- The service must broadcast each published message to all currently
  active participants.
- Each broadcast must include:
  * Message content
  * Logical timestamp (e.g., Lamport timestamp)


S5. Joining
------------
- When a new participant X joins, the service must broadcast:
  "Participant X joined Chit Chat at logical time L"
- This message must be delivered to all participants, including X.


S6. Leaving
------------
- When a participant X leaves, the service must broadcast:
  "Participant X left Chit Chat at logical time L"
- This message must be delivered to all remaining participants.


S7. Receiving Messages
-----------------------
When a participant receives any broadcast message, it must:
1. Display the message content and logical timestamp.
2. Log the message content and logical timestamp.


============================================================
⚙️ TECHNICAL REQUIREMENTS
============================================================

- Implementation language: Go
- Communication framework: gRPC
- Message definitions: Protocol Buffers (.proto)
- Logging: Go’s "log" standard library
- Concurrency: Use goroutines and channels
- Each process (client/server) runs independently
- The server must:
  * Handle each client connection in a dedicated goroutine
  * Support multiple concurrent client connections
  * Ensure non-blocking message delivery


============================================================
🪵 SYSTEM LOGGING
============================================================

You must log the following events:

Event Type                      | Description
--------------------------------|---------------------------------
Server startup/shutdown          | When the server starts or stops
Client connection/disconnection  | When clients join or leave
Broadcasts                       | Join/leave message broadcasts
Message delivery                 | When a message is delivered to a client

Each log entry must include:
- Timestamp
- Component name (Server/Client)
- Event type
- Relevant identifiers (e.g., Client ID)


============================================================
🧠 MINIMUM TEST CONFIGURATION
============================================================

- At least 3 nodes:
  * 1 server
  * 2 clients
- Must demonstrate:
  * A client joining
  * A client leaving


============================================================
📄 HAND-IN REQUIREMENTS
============================================================

Submit a single PDF report via LearnIT and include a link to your
Git repository.

Your report must include:

1. Streaming Discussion
   - Explain whether you use:
     * Server-side streaming
     * Client-side streaming
     * Bidirectional streaming

2. System Architecture
   - Describe your architecture (client–server, peer-to-peer, etc.)

3. RPC Methods
   - List all implemented RPC methods and their types
   - Describe message types used for communication

4. Timestamp Implementation
   - Explain how logical (Lamport) timestamps are calculated

5. Sequence Diagram
   - Show a sequence of RPC calls with Lamport timestamps, e.g.:
     Client X joins → Client X publishes → ... → Client X leaves

6. System Logs
   - Include logs that demonstrate the requirements are met
   - Logs must appear:
     * In the appendix of your report
     * In your repository


============================================================
📁 REPOSITORY STRUCTURE
============================================================

project-root/
├── client/        # contains the client code
├── grpc/          # contains the .proto file
├── server/        # contains the server code
└── readme.md      # explains how to run the program


============================================================
📘 README REQUIREMENTS
============================================================

Your readme.md must clearly explain:
- How to build and run the system (server + clients)
- Any dependencies or setup instructions
- Example usage or test scenario
