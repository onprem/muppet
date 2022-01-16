# `muppet`

A poor man's distributed puppet service.

This project contains two applications:

- `muppet-service`: An HTTP service that listens for requests to run commands on hosts.
- `command-agent`: An agent that runs on hosts and polls `muppet-service` for any commands that it need ti run on the host.

## Building

- The project is written in Go. You'll need to Go toolchain installed on your machine. See: [https://go.dev/doc/install](https://go.dev/doc/install).
- Run `make build` to build both binaries. Or you can build the service binary by running `go build ./cmd/service`, and similarly the agent binary by running `go build ./cmd/agent` directly.

### Building containers

- Run `make container` to build docker containers for both services using the provided Dockerfiles.

## Usage

You need to run one instance of `muppet-service` in a cental place and on avery host server, a `command-agent`. The command agent will constantly poll the `muppet-service` for commands to run for the specified hostname.

### Adding commands for a host

You can add a command to the service by using the follwing `curl` request:

```bash
curl --header "Content-Type: application/json" \
      --request POST \
      --data '{"shell_command": "apt update", "host": "host001"}' \
      localhost:8080/api/v1/commands
```

For a full reference of the API, take a look at the [OpenAPI spec here](./pkg/api/spec.yaml).

### Muppet Service

You can run it by running the binary directly

```bash
./muppet-service
```

or by running it using docker

```bash
docker run -p 8080:8080 onprem/muppet-service
```

### Command Agent

You can run it by running the binary on the host

```bash
./command-agent --hostname $HOSTNAME
```

or by running it using docker

```bash
docker run onprem/command-agent --hostname=host001 --service-url=http://muppet-service:8080
```

> NOTE: If you use docker to run the command agent, the commands will be executed inside the container.

[embedmd]: # "tmp/help-agent.txt"

```txt
Usage of ./command-agent:
      --hostname string      The hostname to fetch commands for.
      --interval uint        The interval at which to poll the muppet service for commands to execute, given in seconds. (default 60)
      --service-url string   The URL of muppet service to fetch commands from. (default "http://localhost:8080")
pflag: help requested
```

## Security Considerations

- The `muppet-service` and `command-agent` communicate over unencrypted HTTP connection. This means they are suspectible to [Man-In-The-Middle](https://en.wikipedia.org/wiki/Man-in-the-middle_attack) attacks. A malicious actor can impersonate `muppet-service` and execute any command on the hosts running `command-agent`. A potential mitigation of this RCE vulnerability is serving the `muppet-service` API over HTTPS.

- Although we can mitigate the MITM attack by using HTTPS, there is still an Information Disclosure vulnerability present. Since the `command-agents` are not authenticated, any malicious actor can see what commands are there in queue for a given host. The malicious actor can also mark any command as done, before it is actually picked by the `command-agent`. To mitigate this, we can use [mutual TLS](https://www.cloudflare.com/en-in/learning/access-management/what-is-mutual-tls/) (`mTLS`) authentication between `muppet-service` and `command-agent`.

- Similar to this, anyone can submit a command to be run at any host, causing RCE. This can also be mitigated by using `mTLS` authentication between `muppet-service` and the system administrator.

## Scalability Considerations

- As of now, the `muppet-service` is not horizontally scalable. The running instance contains some state in-memory, like the commands in queue for a host. This state won't be shared with other instances, so we can't just run multiple replicas of `muppet-service` and load balance between them. The solution is just to add another store implementation that can be shared between multiple instances. For example, a SQL database backed store.

- From reliability point of view, certain things can be improved. The code doesnt implement retries in a lot of plcaes where it can help. For example, if marking a command done fails, we just log a warning and skip it. This can be improved by rertying the request with exponential backoff.
