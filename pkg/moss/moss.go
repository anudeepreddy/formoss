package moss

import (
	"errors"
	"strconv"
)

type File struct {
	path string
	name string
}

type moss struct {
	id        int
	mode      int
	language  string
	baseFiles []File
	codeFiles []File
	m         int
	comment   string
	n         int
}

func New(id int, mode int, language string) moss {
	e := moss{id, mode, language, nil, nil, 10, "", 250}
	return e
}

func (e *moss) AddBaseFile(path string, name string) {
	if e.baseFiles == nil {
		file := File{path, name}
		e.baseFiles = []File{file}
		return
	}
	e.baseFiles = append(e.baseFiles, File{path, name})
}

func (e *moss) AddCodeFile(path string, name string) {
	if e.codeFiles == nil {
		file := File{path, name}
		e.codeFiles = []File{file}
		return
	}
	e.codeFiles = append(e.codeFiles, File{path, name})
}

func (e *moss) SetParameters(m int, comment string, n int) {
	e.m = m
	e.comment = comment
	e.n = n
}

func (e moss) Start() (string,error) {
	s, err := getConnection()
	if err != nil {
		return "", errors.New("problem establishing a connection")
	}

	s.send("moss " + strconv.Itoa(e.id))
	s.send("directory " + strconv.Itoa(e.mode))
	s.send("X 0")
	s.send("maxmatches " + strconv.Itoa(e.m))
	s.send("show " + strconv.Itoa(e.n))
	s.send("language " + e.language)

	message, _ := s.recv(1024)

	if message == "no" {
		s.close()
		return "", errors.New("unsupported Language")
	}

	for _, baseFile := range e.baseFiles {
		s.uploadFile(baseFile, e.language, 0)
	}

	i := 1
	for _, codeFile := range e.codeFiles {
		s.uploadFile(codeFile, e.language, i)
		i++
	}

	s.send("query 0 "+e.comment)

	response,_ := s.recv(1024)
	s.close()

	return response, nil
}
