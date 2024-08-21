package environment

import "os"

func IsDockerEnv() bool {
	return os.Getenv("ENVIRONMENT") == "docker"
}
