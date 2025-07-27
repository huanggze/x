package healthx

const (
	// AliveCheckPath is the path where information about the life state of the instance is provided.
	AliveCheckPath = "/health/alive"
	// ReadyCheckPath is the path where information about the ready state of the instance is provided.
	ReadyCheckPath = "/health/ready"
	// VersionPath is the path where information about the software version of the instance is provided.
	VersionPath = "/version"
)
