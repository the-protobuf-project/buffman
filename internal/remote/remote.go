// Package remote provides abstractions for interacting with various remote Git
// repositories. It includes a factory and a manager for handling different
// remote provider types, such as GitHub.
package remote

import "errors"

// RemoteType defines the type of remote repository provider.
type RemoteType string

const (
	// Github represents a standard Git repository, typically hosted on GitHub.
	Github RemoteType = "github"
)

// PullOptions specifies the parameters for a remote pull operation.
type PullOptions struct {
	// Out is the local directory path where the repository will be cloned.
	Out string
	// Url is the remote URL of the Git repository to clone.
	Url string
	// Commit is an optional Git commit hash. If provided, the repository
	// will be checked out to this specific commit after cloning.
	Commit *string
}

// Remote defines the interface for all remote repository operations.
type Remote interface {
	// Pull clones a repository according to the specified options.
	Pull(opts PullOptions) error
}

// NewRemote acts as a factory, returning a concrete implementation of the Remote
// interface based on the provided RemoteType. It returns an error if the
// requested type is unsupported.
func NewRemote(remote RemoteType) (Remote, error) {
	switch remote {
	case Github:
		return &github{}, nil
	default:
		return nil, errors.New("unsupported remote type")
	}
}

// Manager holds and manages a collection of registered Remote implementations.
type Manager struct {
	remote map[RemoteType]Remote
}

// NewManager initializes and returns a new, empty Manager.
func NewManager() *Manager {
	return &Manager{
		remote: make(map[RemoteType]Remote),
	}
}

// RegisterRemote adds a Remote implementation to the manager, associating it
// with a specific RemoteType.
func (m *Manager) RegisterRemote(remoteType RemoteType, r Remote) {
	m.remote[remoteType] = r
}

// GetRemote retrieves a registered Remote implementation by its type. It returns
// nil if no remote provider is registered for the given type.
func (m *Manager) GetRemote(remoteType RemoteType) Remote {
	return m.remote[remoteType]
}
