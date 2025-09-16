package railway

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var projectID = os.Getenv("RAILWAY_PROJECT_ID")

var ErrNotRailway = errors.New("not running on Railway")

// Env is the Railway-provided environment variables.  See https://docs.railway.com/reference/variables#railway-provided-variables for more details.
type Env struct {
	// The public service or customer domain, of the form example.up.railway.app
	PublicDomain string

	// The private DNS name of the service.
	PrivateDomain string

	// (see TCP Proxy for details) The public TCP proxy domain for the service, if applicable. Example: roundhouse.proxy.rlwy.net
	TCPProxyDomain string

	// (see TCP Proxy for details) The external port for the TCP Proxy, if applicable. Example: 11105
	TCPProxyPort int

	// (see TCP Proxy for details) The internal port for the TCP Proxy, if applicable. Example: 25565
	TCPApplicationPort int

	// The project name the service belongs to.
	ProjectName string

	// The project id the service belongs to.
	ProjectID string

	// The environment name of the service instance.
	EnvironmentName string

	// The environment id of the service instance.
	EnvironmentID string

	// The service name.
	ServiceName string

	// The service id.
	ServiceID string

	// The replica ID for the deployment.
	ReplicaID string

	// The region where the replica is deployed. Example: us-west2
	ReplicaRegion string

	// The ID for the deployment.
	DeploymentID string

	// The snapshot ID for the deployment.
	SnapshotID string

	// The name of the attached volume, if any. Example: foobar
	VolumeName string

	// The mount path of the attached volume, if any. Example: /data
	VolumeMountPath string

	// The git SHA of the commit that triggered the deployment. Example: d0beb8f5c55b36df7d674d55965a23b8d54ad69b
	GitCommitSHA string

	// The user of the commit that triggered the deployment. Example: gschier
	GitAuthor string

	// The branch that triggered the deployment. Example: main
	GitBranch string

	// The name of the repository that triggered the deployment. Example: myproject
	GitRepoName string

	// The name of the repository owner that triggered the deployment. Example: mycompany
	GitRepoOwner string

	// The message of the commit that triggered the deployment. Example: Fixed a few bugs
	GitCommitMessage string

	// How long the old deploy will overlap with the newest one being deployed, its default value is 0. Example: 20
	DeploymentOverlapSeconds int

	// The path to the Dockerfile to be used by the service, its default value is Dockerfile. Example: Railway.dockerfile
	DockerfilePath string

	// The path to the Nixpacks configuration file relative to the root of the app, its default value is nixpacks.toml. Example: frontend.nixpacks.toml
	NixpacksConfigFile string

	// The version of Nixpacks to use, if unspecfied a default version will be used. Example: 1.29.1
	NixpacksVersion string

	// The timeout length (in seconds) of healthchecks. Example: 300
	HealthcheckTimeoutSec int

	// The SIGTERM to SIGKILL buffer time (in seconds), its default value is 0. Example: 30
	DeploymentDrainingSeconds int

	// The UID of the user which should run the main process inside the container. Set to 0 to explicitly run as root.
	RunUID int

	// This variable accepts a value in binary bytes, with a default value of 67108864 bytes (64 MB)
	SHMSizeBytes int64
}

func IsRailway() bool {
	return projectID != ""
}

func Must(env Env, err error) Env {
	if err != nil {
		panic(err)
	}
	return env
}

func MustLoad() Env {
	return Must(Load())
}

func Load() (Env, error) {
	if !IsRailway() {
		return Env{}, ErrNotRailway
	}

	tcpProxyPort, err := getEnvInt("RAILWAY_TCP_PROXY_PORT", 0)
	if err != nil {
		return Env{}, err
	}

	tcpApplicationPort, err := getEnvInt("RAILWAY_TCP_APPLICATION_PORT", 0)
	if err != nil {
		return Env{}, err
	}

	deploymentOverlapSeconds, err := getEnvInt("RAILWAY_DEPLOYMENT_OVERLAP_SECONDS", 0)
	if err != nil {
		return Env{}, err
	}

	healthcheckTimeoutSec, err := getEnvInt("RAILWAY_HEALTHCHECK_TIMEOUT_SEC", 0)
	if err != nil {
		return Env{}, err
	}

	deploymentDrainingSeconds, err := getEnvInt("RAILWAY_DEPLOYMENT_DRAINING_SECONDS", 0)
	if err != nil {
		return Env{}, err
	}

	runUID, err := getEnvInt("RAILWAY_RUN_UID", 0)
	if err != nil {
		return Env{}, err
	}

	shmSizeBytes, err := getEnvInt64("RAILWAY_SHM_SIZE_BYTES", 0)
	if err != nil {
		return Env{}, err
	}

	return Env{
		PublicDomain:              os.Getenv("RAILWAY_PUBLIC_DOMAIN"),
		PrivateDomain:             os.Getenv("RAILWAY_PRIVATE_DOMAIN"),
		TCPProxyDomain:            os.Getenv("RAILWAY_TCP_PROXY_DOMAIN"),
		TCPProxyPort:              tcpProxyPort,
		TCPApplicationPort:        tcpApplicationPort,
		ProjectName:               os.Getenv("RAILWAY_PROJECT_NAME"),
		ProjectID:                 projectID,
		EnvironmentName:           os.Getenv("RAILWAY_ENVIRONMENT_NAME"),
		EnvironmentID:             os.Getenv("RAILWAY_ENVIRONMENT_ID"),
		ServiceName:               os.Getenv("RAILWAY_SERVICE_NAME"),
		ServiceID:                 os.Getenv("RAILWAY_SERVICE_ID"),
		ReplicaID:                 os.Getenv("RAILWAY_REPLICA_ID"),
		ReplicaRegion:             os.Getenv("RAILWAY_REPLICA_REGION"),
		DeploymentID:              os.Getenv("RAILWAY_DEPLOYMENT_ID"),
		SnapshotID:                os.Getenv("RAILWAY_SNAPSHOT_ID"),
		VolumeName:                os.Getenv("RAILWAY_VOLUME_NAME"),
		VolumeMountPath:           os.Getenv("RAILWAY_VOLUME_MOUNT_PATH"),
		GitCommitSHA:              os.Getenv("RAILWAY_GIT_COMMIT_SHA"),
		GitAuthor:                 os.Getenv("RAILWAY_GIT_AUTHOR"),
		GitBranch:                 os.Getenv("RAILWAY_GIT_BRANCH"),
		GitRepoName:               os.Getenv("RAILWAY_GIT_REPO_NAME"),
		GitRepoOwner:              os.Getenv("RAILWAY_GIT_REPO_OWNER"),
		GitCommitMessage:          os.Getenv("RAILWAY_GIT_COMMIT_MESSAGE"),
		DeploymentOverlapSeconds:  deploymentOverlapSeconds,
		DockerfilePath:            os.Getenv("RAILWAY_DOCKERFILE_PATH"),
		NixpacksConfigFile:        os.Getenv("NIXPACKS_CONFIG_FILE"),
		NixpacksVersion:           os.Getenv("NIXPACKS_VERSION"),
		HealthcheckTimeoutSec:     healthcheckTimeoutSec,
		DeploymentDrainingSeconds: deploymentDrainingSeconds,
		RunUID:                    runUID,
		SHMSizeBytes:              shmSizeBytes,
	}, nil
}

func getEnvInt(env string, defaultValue int) (int, error) {
	if value := os.Getenv(env); value != "" {
		res, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("invalid %s: %w", env, err)
		}
		return res, nil
	}
	return defaultValue, nil
}

func getEnvInt64(env string, defaultValue int64) (int64, error) {
	if value := os.Getenv(env); value != "" {
		res, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid %s: %w", env, err)
		}
		return res, nil
	}
	return defaultValue, nil
}
