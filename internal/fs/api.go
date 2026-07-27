package fs

var fsys FileSystem

func FS() FileSystem {
	if fsys == nil {
		fsys = &osFileSystem{}
	}
	return fsys
}

func Set(f FileSystem) {
	fsys = f
}

func Reset() {
	fsys = &osFileSystem{}
}
