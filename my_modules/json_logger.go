package my_modules

import (
	"fmt"
	"os"
)

var LogPath string = ""

type MyError struct{}

func (m *MyError) Error() string {
	return "Invalid file path"
}

func LogToJSON(json_obj interface{}) error {
	if LogPath == "" {
		return &MyError{}
	}

	var _json_bytes []byte
	var err error
	if _json_bytes, err = JSONMarshal(json_obj); err != nil {
		return err
	}

	// fmt.Println(json_obj)
	// fmt.Println(string(_json_bytes))

	_file, err := os.OpenFile(LogPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer _file.Close()

	fileinfo, err := _file.Stat()
	if err != nil {
		return err
	}

	BUFFERSIZE := fileinfo.Size()
	// fmt.Printf("BUFFERSIZE %v\n", BUFFERSIZE)

	// if writing file first time
	if BUFFERSIZE == 0 {
		if _, err := _file.Write([]byte(fmt.Sprintf("[%v]", string(_json_bytes)))); err != nil {
			return err
		}
		return nil
	}

	// skip EOF reperesentation
	_file.Seek(-1, 2)

	// if file already has content

	var ret int64 = BUFFERSIZE
	single_char := make([]byte, 1)
	for ret >= 0 {

		// reverse travel
		// seek to end[2], seek back [-1]
		if ret, err = _file.Seek(-1, 2); err != nil {
			return err
		}
		// read a single character
		if _, err := _file.Read(single_char); err != nil {
			return err
		}

		// fmt.Printf("c=%s,ret=%v\n", string(single_char),ret)

		// check the character is matches ]
		if string(single_char) == "]" {
			break
		}
	}
	if _, err = _file.Seek(-1, 2); err != nil {
		return err
	}

	// append json entry to json array
	if _, err := _file.Write([]byte(fmt.Sprintf(",%v]", string(_json_bytes)))); err != nil {
		return err
	}

	return nil
}
