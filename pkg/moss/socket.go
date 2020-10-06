package moss

import (
	"errors"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

const connString string = "moss.stanford.edu:7690"

type socket struct {
	connection net.Conn
}

func getConnection() (socket, error) {
	conn, err := net.Dial("tcp", connString)
	if err != nil {
		return socket{connection: nil}, err
	}
	return socket{connection: conn}, nil
}

func (s *socket) send(message string) error {
	_, err := s.connection.Write([]byte(message + "\n"))

	if err != nil {
		return errors.New("error writing data to connection -> "+err.Error())
	}
	return nil
}

func (s *socket) recv(size int) (string, error) {
	message := make([]byte, size)
	_, err := s.connection.Read(message)
	if err != nil {
		return "", errors.New("error reading data from connection -> "+err.Error())
	}
	return string(message),nil
}

func (s *socket) uploadFile(file File, lang string, fileId int) error{
	info, err := os.Stat(file.path)
	if err != nil {
		return errors.New("error accessing file:"+file.path)
	}

	message := "file " + strconv.Itoa(fileId) + " " + lang + " " + strconv.FormatInt(info.Size(), 10) + " " + file.name

	s.send(message)
	dat, err := ioutil.ReadFile(file.path)
	if err != nil {
		return errors.New("error reading file:"+file.path)
	}
	_, e := s.connection.Write(dat)
	if e != nil {
		return errors.New("error writing file:"+file.path+":to connection")
	}
	return nil
}

func (s *socket) close() {
	s.connection.Write([]byte("end\n"))
	s.connection.Close()
}
