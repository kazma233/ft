package cmd

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"ft/constants"
	"ft/entity"
	"io"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

var (
	basePath string
	addr     string
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "start send file server",
	Long:  "when client connect the server, start file send.",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("%v, %v", basePath, addr)

		listener, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("accept error:  %v", err)
				continue
			}

			go handler(&conn)
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringVarP(&basePath, "path", "p", "./fs", "the file(dir) path")
	sendCmd.Flags().StringVarP(&addr, "addr", "a", ":9990", "server listen addr")
}

func handler(connRef *net.Conn) error {
	f, err := os.Open(basePath)
	if err != nil {
		return err
	}

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	if fi.IsDir() {
		log.Println("暂时不支持发送文件夹")
	} else {
		fName := f.Name()
		log.Printf("start send file: %s", fName)
		err = sendFile(f, fi, connRef)
		if err != nil {
			log.Printf("send file %s error: %v", fName, err)
		} else {
			log.Printf("end send file: %v", fName)
		}
	}

	conn := *connRef
	defer conn.Close()

	return nil
}

func sendFile(f *os.File, fi os.FileInfo, conn *net.Conn) error {

	fileNameMessage := &entity.Message{
		MsgType: entity.MessageType_TEXT,
		BaseMessage: &entity.BaseMessage{
			TextType: entity.TextType_FILENAME,
			FileMessage: &entity.FileMessage{
				Path: f.Name(),
				Name: fi.Name(),
			},
		},
	}

	// 发送文件名
	err := sendMessage(fileNameMessage, conn)
	if err != nil {
		return err
	}

	// 发送文件
	writeBytes := make([]byte, constants.DefaultByteSize)
	_sha1 := sha1.New()
	for {
		n, err := f.Read(writeBytes)
		if err != nil {
			// 发送结束
			if err == io.EOF {
				fileSha1Message := &entity.Message{
					MsgType: entity.MessageType_TEXT,
					FileContent: &entity.FileContent{
						Sha1: fmt.Sprintf("%X", _sha1.Sum(nil)),
					},
				}

				err = sendMessage(fileSha1Message, conn)
				if err != nil {
					return err
				}

				break
			}

			return err
		}

		wb := writeBytes[:n]
		_sha1.Write(wb)

		fileMessage := &entity.Message{
			MsgType: entity.MessageType_FILE,
			FileContent: &entity.FileContent{
				Data: wb,
				Sha1: "",
			},
		}
		err = sendMessage(fileMessage, conn)
		if err != nil {
			return err
		}
	}

	return sendData(constants.SLine, conn)
}

func sendMessage(msg *entity.Message, connRef *net.Conn) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil
	}

	return sendData(data, connRef)
}

func sendData(data []byte, connRef *net.Conn) error {
	conn := *connRef
	buf := bytes.NewBuffer(data)

	// write size
	sizeBuf := make([]byte, 8) // long; uint64
	binary.BigEndian.PutUint64(sizeBuf, uint64(len(data)))
	_, _ = conn.Write(sizeBuf)

	// write data
	_, err := buf.WriteTo(conn)
	return err
}
