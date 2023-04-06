package archiver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"com.geophone.observer/features/ntpclient"
)

func WriteMessage(path string, options *ArchiverOptions) {
	if !options.Enable {
		return
	}

	exist, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				options.OnErrorCallback(err)
				return
			}
		} else {
			options.OnErrorCallback(err)
			return
		}
	} else {
		if !exist.IsDir() {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				options.OnErrorCallback(err)
				return
			}
		}
	}

	data, err := json.Marshal(options.Message)
	if err != nil {
		options.OnErrorCallback(err)
		return
	}

	file, err := func(path string, timestamp int64) (*os.File, error) {
		if options.Name == "timestamp" {
			return os.Create(
				fmt.Sprintf("%s/%d.json", path, timestamp),
			)
		}

		s := ntpclient.FormatTime(timestamp, "060102150405")
		return os.Create(
			fmt.Sprintf("%s/%s.json", path, s),
		)
	}(path, ntpclient.AlignTime(options.Status.Offset))
	if err != nil {
		options.OnErrorCallback(err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		options.OnErrorCallback(err)
		return
	}

	writer.Flush()
	options.OnCompleteCallback()
}
