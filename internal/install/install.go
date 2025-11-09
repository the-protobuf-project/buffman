package install

import (
	"fmt"

	"github.com/machanirobotics/buffman/internal/install/flatc"
)

// BufferInstaller defines the type for different buffer package installers.
type BufferInstaller string

// Defined constants for different buffer installers.
var (
	// ProtobufInstaller is the installer for the protobuf package.
	ProtobufInstaller BufferInstaller = "protobuf"
	// FlatbuffersInstaller is the installer for the flatbuffers package.
	FlatbuffersInstaller BufferInstaller = "flatbuffers"
	// CapnProtoInstaller is the installer for the capnproto package.
	CapnProtoInstaller BufferInstaller = "capnproto"
	// NanopbInstaller is the installer for the nanopb package.
	NanopbInstaller BufferInstaller = "nanopb"
)

// Installer is the interface for buffer installers.
type Installer interface {
	Install() error
	Validate() error // To check if the installed tool works correctly (e.g., flatc --version)
	Uninstall() error
	Exists() bool                // To check if the tool is already present
	GetVersion() (string, error) // To retrieve the installed tool's version
}

// NewInstaller is a factory function that returns an appropriate Installer implementation.
func NewInstaller(installer BufferInstaller) Installer {
	switch installer {
	case ProtobufInstaller:
		// return NewProtobufInstaller()
	case FlatbuffersInstaller:
		return flatc.NewFlatbuffersInstaller()
	case CapnProtoInstaller:
		// return NewCapnProtoInstaller()
	case NanopbInstaller:
		// return NewNanopbInstaller()
	default:
		// You might want to return an error here instead of nil, or a no-op installer.
		// For simplicity, returning nil as per your original NewInstaller.
		return nil
	}
	return nil // Return nil if no valid installer is found
}

// Install installs the buffer package using the specified installer.
// This is a convenience function that uses the NewInstaller factory.
func Install(installer BufferInstaller) error {
	ins := NewInstaller(installer)
	if ins == nil {
		return fmt.Errorf("unsupported installer type: %s", installer)
	}
	return ins.Install()
}

// Manager struct for orchestrating multiple installers or high-level operations.
// This is an optional but good pattern for more complex scenarios.
type Manager struct {
	installers map[BufferInstaller]Installer
}

// NewManager creates a new Manager instance.
func NewManager() *Manager {
	return &Manager{
		installers: make(map[BufferInstaller]Installer),
	}
}

// RegisterInstaller allows you to register an installer with the manager.
func (m *Manager) RegisterInstaller(installerType BufferInstaller, ins Installer) {
	m.installers[installerType] = ins
}

// GetInstaller retrieves a registered installer.
func (m *Manager) GetInstaller(installerType BufferInstaller) Installer {
	return m.installers[installerType]
}

// InstallAll can be a manager-level function to install all registered buffer tools.
func (m *Manager) InstallAll() error {
	fmt.Println("Attempting to install all known buffer tools...")
	for _, ins := range m.installers {
		if !ins.Exists() {
			fmt.Printf("Installing %T...\n", ins) // %T prints the type of the struct
			if err := ins.Install(); err != nil {
				return fmt.Errorf("failed to install %T: %w", ins, err)
			}
			fmt.Printf("%T installed successfully.\n", ins)
		} else {
			fmt.Printf("%T already exists.\n", ins)
		}
	}
	fmt.Println("All buffer tools checked/installed.")
	return nil
}
