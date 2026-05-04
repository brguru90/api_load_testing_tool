package my_modules

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

var LogPath string = ""
var DisableLogging bool = false

// BufferedLogging switches LogToJSON from the default per-call open/seek/append
// pattern (durable but pays open+stat+seek+write+close every call) to an
// in-memory buffered writer that opens the file once, accumulates entries in a
// 1 MiB bufio.Writer, and flushes only on FlushBufferedLogs. Trade-off: if the
// process crashes mid-run the file is left without a closing "]". For
// benchmark workloads this is almost always the right choice.
//
// Toggle either by setting BUFFERED_LOGGING=true in the environment or by
// assigning BufferedLogging = true in main() before any LogToJSON call.
var BufferedLogging bool = false

func init() {
	if os.Getenv("BUFFERED_LOGGING") == "true" {
		BufferedLogging = true
	}
}

type MyError struct{}

func (m *MyError) Error() string {
	return "Invalid file path"
}

// Per-path buffered writer. Created lazily on first write to a given path.
type bufferedLog struct {
	f       *os.File
	bw      *bufio.Writer
	started bool
	mu      sync.Mutex
}

const bufferedLogChunk = 1 << 20 // 1 MiB

var (
	bufferedLogsMu sync.Mutex
	bufferedLogs   = map[string]*bufferedLog{}
)

func getBufferedLog(path string) *bufferedLog {
	bufferedLogsMu.Lock()
	defer bufferedLogsMu.Unlock()
	if bl, ok := bufferedLogs[path]; ok {
		return bl
	}
	bl := &bufferedLog{}
	bufferedLogs[path] = bl
	return bl
}

func bufferedAppend(path string, jsonBytes []byte) error {
	bl := getBufferedLog(path)
	bl.mu.Lock()
	defer bl.mu.Unlock()
	if !bl.started {
		// Truncate on first write of the run — same effective semantics as
		// the instant-mode empty-file branch (start a fresh JSON array).
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		bl.f = f
		bl.bw = bufio.NewWriterSize(f, bufferedLogChunk)
		if _, err := bl.bw.WriteString("["); err != nil {
			return err
		}
		bl.started = true
	} else {
		if _, err := bl.bw.WriteString(","); err != nil {
			return err
		}
	}
	_, err := bl.bw.Write(jsonBytes)
	return err
}

// FlushBufferedLogs writes the closing "]", flushes the bufio.Writer, and
// closes every buffered log file. Idempotent. Called automatically from
// OnBenchmarkEnd; callers running outside the benchmark harness MUST invoke
// this themselves before process exit, otherwise the JSON arrays will be
// left without a closing "]".
func FlushBufferedLogs() error {
	bufferedLogsMu.Lock()
	all := make([]*bufferedLog, 0, len(bufferedLogs))
	for _, bl := range bufferedLogs {
		all = append(all, bl)
	}
	bufferedLogsMu.Unlock()

	var firstErr error
	for _, bl := range all {
		bl.mu.Lock()
		if bl.started {
			if _, err := bl.bw.WriteString("]"); err != nil && firstErr == nil {
				firstErr = err
			}
			if err := bl.bw.Flush(); err != nil && firstErr == nil {
				firstErr = err
			}
			if err := bl.f.Close(); err != nil && firstErr == nil {
				firstErr = err
			}
			bl.started = false
			bl.f = nil
			bl.bw = nil
		}
		bl.mu.Unlock()
	}
	return firstErr
}

func LogToJSON(json_obj interface{}, log_path interface{}) error {
	if DisableLogging {
		return nil
	}

	if LogPath == "" {
		return &MyError{}
	}
	final_log_path := LogPath
	if log_path != nil {
		final_log_path = log_path.(string)
	}

	var _json_bytes []byte
	var err error
	if _json_bytes, err = JSONMarshal(json_obj, true); err != nil {
		return err
	}

	if BufferedLogging {
		return bufferedAppend(final_log_path, _json_bytes)
	}

	// Instant mode: open + seek-to-bracket + append "," + entry + "]" + close.
	// Durable per-call (file is valid JSON after every write) but pays full
	// open/stat/write/close syscall overhead every call.
	_file, err := os.OpenFile(final_log_path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer _file.Close()

	fileinfo, err := _file.Stat()
	if err != nil {
		return err
	}

	BUFFERSIZE := fileinfo.Size()

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
