package portrepository

// File represents a file to be stored or retrieved
type File struct {
	Name string
	Data []byte
}

// StoragePort defines the interface for interacting with storage systems
type StoragePort interface {
	Upload(file File) error
	// Download(fileName string) (File, error)
	// Delete(fileName string) error
}
