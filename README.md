# rcmd
A simple client server application to allow for commands to be executed on remote hosts. It fatures a server that is provided with the command, and a client that should be used (e.g. via cron) to poll the server and pull the command to execute. Once the command is executed it send the result back to the server and the server exits.

Written using golang, so the client application has minimal dependencies and can easily be deployed onto hosts.