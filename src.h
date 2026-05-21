IDA Pro - A powerful disassembler, decompiler and a versatile debugger. In one tool.
Disassemble almost anything.
IDA Disassembler excels in supporting various processors and file formats. Its versatility makes it ideal for analyzing embedded systems, mobile apps, or complex software, ensuring you have the best tools for any task.
High-quality, readable, and maintainable pseudocode.
IDA decompilers focus on delivering code that is readable, maintainable, and semantically similar to the original source code thanks to high-level abstractions, semantic preservation, readability, type inference, structure recovery and more.

After downloading the target binary from the web service, load it into IDA Pro using the x64 analysis mode. Wait for the auto-analysis process to complete — the status bar at the bottom should change from “... | Down” to “... | Idle” before proceeding with the reverse engineering process.

Navigate to the Functions window on the left side of the interface. Most of the application-specific routines are prefixed with `main__`, making them easy to identify. Locate the function named `main__handleLogin`, which is responsible for processing requests for the `/login` endpoint.

Inside the function, observe the conditional branch leading to an “invalid username or password” routine. This indicates that the authentication verification logic is nearby. Slightly above this branch, identify the function call `main_credentialsMatch`. Double-click the symbol to follow the cross-reference and open the authentication validation routine.

Within this function, two important helper routines become visible:

* `main_expectedLogin`
* `main_expectedPassword`

The `main_expectedLogin` routine clearly reveals the expected username:

The password string can be identified directly before the `accept-charset` sequence in memory/string references. The extracted password is: 	

Use the recovered credentials to authenticate against the service and obtain the flag.

