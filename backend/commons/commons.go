package commons

import (
	"math/rand"
	"strconv"
	"time"
	"net/http"
	"log"
	"reflect"
	"strings"
	"os"
	"fmt"
	"io"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetUID() string {
	s := 6
	return randomString(s)
}

func randomString(l int) string {
	r := strconv.Itoa(rand.Intn(10000))

	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(64, 90))
	}
	return string(bytes) + r
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func IsInArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func FileBasename(s string) string {
	n := strings.LastIndexByte(s, '.')
	if n >= 0 {
		return s[:n]
	}
	return s
}

func FileOrFolderExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Thx to https://www.socketloop.com/tutorials/golang-copy-directory-including-sub-directories-files
func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		si, err := os.Stat(source)
		if err == nil {
			err = os.Chmod(dest, si.Mode())
		}
	}

	return
}

// Thx to https://www.socketloop.com/tutorials/golang-copy-directory-including-sub-directories-files
func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {
		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func WriteResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(" "+ message))
}

func Check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}
