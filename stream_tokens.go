/*
Functions to traverse directories and get file names and stream input file through a tokenizer to an output file
*/
package nlpt_tkz

import (
	"bufio"
	"bytes"
	"gopkg.in/pipe.v2"
	"io"
	"os"
	"path/filepath"
	"time"
)

func StreamTokenizedDirectory(directoryPath, writeFile, tkzType string, timeoutLimit time.Duration) {
	//overwrite the output file
	f, err := os.Create(writeFile)
	f.Close()
	if err != nil {
		panic(err)
		return
	}

	handler := NewFileHandler(directoryPath, "space")
	go func(handler *FileHandler, timeoutLimit time.Duration, tkzType, fileWrite string) {
		for _, file := range handler.FullFilePaths {
			p := pipe.Line(
				PipeFileTokens(file, tkzType),
				//pipe.Filter(func(line []byte) bool { return stopList.IsStopWord[string(line)] }),
				pipe.AppendFile(fileWrite, 0644),
				//PipeFileTokens(fileWrite, "unicode"),
				//pipe.AppendFile(fileWrite, 0644),
			)
			_, err := pipe.CombinedOutputTimeout(p, timeoutLimit)
			//output, err := pipe.CombinedOutputTimeout(p, timeoutLimit)
			if err != nil {
				panic(err)
			}

			/// *************** DEBUGGING ****************
			//Log.Debug("FILE: %v\n filter %v\n", file, string(output))
			//Log.Debug("FILE: %v\n tokens %v\n", file, string(output))
			/// *************** DEBUGGING ****************
		}
	}(handler, timeoutLimit, tkzType, writeFile)
	//fmt.Printf("read %d files for directory %s", len(handler.FullFilePaths), handler.DirName)
}

// ReadFile reads data from the file at path and writes it to the
// pipe's stdout. I've hijacked the pipe projects ReadFile function
// and stuck a text tokenzer inside of it.
// The tokenizer used here MUST be 'lex' OR 'unicode'. The latter is the fastest but less flexible and comprehensive, while the former is not much slower it will return alot of symbols and punctuation. If all you need is "words" then use the 'unicode' tokenizer.
func PipeFileTokens(readFile, tokenizer string) pipe.Pipe {
	//so we don't fail becuase of bad tokenizer input
	var tkzType string
	switch tokenizer {
	case "unicode":
		tkzType = tokenizer
	default:
		tkzType = "lex"
	}
	//Log.Debug("Using tokenizer type: %s", tkzType)

	return pipe.TaskFunc(func(s *pipe.State) error {
		file, err := os.Open(s.Path(readFile))
		//defer file.Close()
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(file)
		bufferCache := new(bytes.Buffer)
		byteLining := []byte{'\n'} //newline padding bytes for writing to file
		for scanner.Scan() {
			bufferCache.Write(
				TokenizeBytes(scanner.Bytes(), tkzType).Bytes,
			)
			//follow each buffer write with a new line
			bufferCache.Write(byteLining)
		}

		//close file as soon as we can but no sooner.
		file.Close()
		//Log.Debug("streamBytes from tokenzier: %d", bufferCache.Len())
		_, err = io.Copy(s.Stdout, bufferCache)
		if err != nil {
			panic(err)
		}
		//file.Close()
		return err
	})
}

// FileHandler contains the directory path, list of file paths, and function to create full file paths.
type FileHandler struct {
	DirName       string
	DirPath       string
	DocumentLabel string
	Tokenizer     string
	FullFilePaths []string
	FileInfo      []os.FileInfo
	FullPathFn    func(string, string, string) string
}

var separator = string(filepath.Separator)

func NewDirHandler(dirPath, dirLabel, tokenizer string) *FileHandler {
	handler := &FileHandler{
		DirName:       dirPath + separator,
		DirPath:       dirPath,
		Tokenizer:     tokenizer,
		DocumentLabel: dirLabel,
		FullPathFn:    func(dirpath, sep, filename string) string { return dirpath + sep + filename },
	}
	handler.setFileNames()
	return handler
}

func NewFileHandler(dirPath, tokenizer string) *FileHandler {
	handler := &FileHandler{
		DirName:    dirPath + separator,
		DirPath:    dirPath,
		Tokenizer:  tokenizer,
		FullPathFn: func(dirpath, sep, filename string) string { return dirpath + sep + filename },
	}
	handler.setFileNames()
	return handler
}

func (handle *FileHandler) setFileNames() {
	handle.getFileInfo()
	//Log.Debug("number of files %d:", len(handle.FileInfo))
	for _, file := range handle.FileInfo {
		if file.Mode().IsRegular() {
			handle.FullFilePaths = append(
				handle.FullFilePaths,
				handle.FullPathFn(handle.DirPath, separator, file.Name()),
			)
		}
	}
}

func (handle *FileHandler) getFileInfo() {
	//Log.Debug("GetFileInfo for new FileHandler %s:", handle.DirPath)
	d, err := os.Open(handle.DirPath)
	if err != nil {
		panic(err)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		panic(err)
	}
	handle.FileInfo = files
}

func (handle *FileHandler) FileByteSize() map[string]int64 {
	fbs := make(map[string]int64)
	for _, file := range handle.FileInfo {
		if file.Mode().IsRegular() {
			fbs[file.Name()] = file.Size()
		}
	}
	return fbs
}
