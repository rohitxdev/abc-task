# abc-task

Backend task for ABC fitness solutions.

## Getting Started

### Prerequisites

- Go 1.23+
- Docker
- Bash 5.2+

### Installation

1. Clone the repository:

```bash
git clone https://github.com/rohitxdev/abc-task.git
```

2. Change directory to the project:

```bash
cd abc-task
```

3. Build the project:

```bash
./run build
```

4. Run the project:

```bash
./run start
```

## Usage

### Commands

| Command | Description |
| --- | --- |
| `./run watch` | Run live development server |
| `./run build` | Build go app for production release and generate OpenAPI docs |
| `./run start` | Run go app binary |
| `./run docker.watch` | Run live development server in docker |
| `./run docker.build` | Build the Docker image for production |
| `./run docker.push` | Push the production docker image to registry |
| `./run test` | Run tests |
| `./run test.cover` | Run tests and show coverage report |
| `./run benchmark` | Run benchmarks |
| `./run clean` | Clean go mod & cache & remove build artifacts |
| `./run checkpoint` | Create a git checkpoint and push changes to origin |
| `./run pprof` | Start pprof profile |
| `./run upgrade` | Upgrade dependencies |

## Notes

- 'run' is a task runner script that provides various commands for running the project. It is located in the root directory of the project. Use only 'run' to run the project or you may encounter unexpected behavior.
- Some commands in run may not work in this project because I copied parts of the code from my other projects.
- All the environment variables must be set in the '.env' file in the root directory of the project.
