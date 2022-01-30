package files_processor

type FilesProcessor interface {
	Add(file string, path string) error
	Remove(file string) error
	Update(files map[string]string) error
	Exists(file string) bool
	Files() ([]string, error)
}
