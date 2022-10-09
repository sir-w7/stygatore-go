package styx

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Should we have the base file path here? It would be slightly
// redundant, but it'll be a lot more readable.
type FileInfo struct {
	WorkingDir string
	BaseName   string
	Filename   string
}

func GetFileList(args []string, ext string) []string {
	var file_list []string
	var pushFile = func(file_list []string, file string, ext string) []string {
		if filepath.Ext(file)[1:] == ext {
			abs_path, err := filepath.Abs(file)
			if err != nil {
				log.Fatal(err)
			}

			return append(file_list, abs_path)
		}
		return file_list
	}

	for _, arg := range args {
		file_info, err := os.Stat(arg)
		if err != nil {
			log.Fatal(err)
		}

		if file_info.IsDir() {
			dir_entry, err := os.ReadDir(arg)
			if err != nil {
				log.Fatal(err)
			}
			// Check extension even before concatenation for speedup.
			for _, file := range dir_entry {
				file_list = pushFile(file_list, filepath.Join(arg, file.Name()), ext)
			}
		} else {
			file_list = pushFile(file_list, arg, ext)
		}
	}
	return file_list
}

func BaseNameNoExt(file string) string {
	var base = filepath.Base(file)
	return base[:len(base)-len(filepath.Ext(base))]
}

func QueryFileInfo(file string) FileInfo {
	return FileInfo{
		WorkingDir: filepath.Dir(file),
		BaseName:   BaseNameNoExt(file),
		Filename:   filepath.Base(file),
	}
}

func Profile(time_start time.Time, name string) {
	var elapsed = (time.Since(time_start)).Microseconds()
	fmt.Printf("%s <- %.02f ms\n", name, float64(elapsed)/1000.0)
}
