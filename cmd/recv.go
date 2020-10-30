package cmd

import (
	"encoding/binary"
	"errors"
	"ft/entity"
	"ft/utils"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

var (
	savePath   string
	serverAddr string
)

var recvCmd = &cobra.Command{
	Use:   "recv",
	Short: "revice file from server",
	Long:  "connect the server, recv file.",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			return err
		}

		msg, err := readMessage(&conn)
		if err != nil {
			return err
		}

		fName := msg.GetBaseMessage().GetFileMessage().GetPath()
		log.Printf("start recive file: %s", fName)
		if fName == "" {
			return errors.New("data not accept")
		}

		fp := filepath.Join(savePath, fName)
		err = os.MkdirAll(filepath.Dir(fp), 0755)
		if err != nil {
			return err
		}

		f, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}

		defer f.Close()

		for {
			dataMsg, err := readMessage(&conn)
			if err != nil {
				return err
			}

			fc := dataMsg.GetFileContent()
			sha1Val := fc.GetSha1()
			if sha1Val != "" {
				if utils.Sha1File(fName) != sha1Val {
					return errors.New("check sha1 bad")
				}
				break
			}

			f.Write(fc.GetData())
		}

		log.Printf("end revice file: %s", fName)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(recvCmd)

	recvCmd.Flags().StringVarP(&savePath, "path", "p", "./fss", "the save dir path")
	recvCmd.Flags().StringVarP(&serverAddr, "addr", "a", ":9990", "server connect addr")
}

func readMessage(connRef *net.Conn) (*entity.Message, error) {
	data, err := readData(connRef)
	if err != nil {
		return nil, err
	}

	message := &entity.Message{}
	err = proto.Unmarshal(data, message)

	return message, err
}

func readData(connRef *net.Conn) ([]byte, error) {
	conn := *connRef
	// revice size
	sizeBytes := make([]byte, 8)
	_, err := io.ReadFull(conn, sizeBytes)
	if err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint64(sizeBytes)

	// recv data
	data := make([]byte, size)
	n, err := io.ReadFull(conn, data)

	return data[:n], err
}
