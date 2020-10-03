package cmd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"ft/constants"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
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
		fmt.Printf("%v, %v", basePath, addr)

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
		log.Println(sendFile(f, fi, connRef))
	}

	conn := *connRef
	defer conn.Close()

	return nil
}

func sendFile(f *os.File, fi os.FileInfo, conn *net.Conn) error {
	writeBytes := make([]byte, constants.DefaultByteSize)
	for {
		n, err := f.Read(writeBytes)
		if err != nil {
			log.Println(err)
			break
		}

		// msg := entity.Message{
		// 	Data: writeBytes[:n],
		// 	Name: fi.Name(),
		// 	Size: fi.Size(),
		// 	End:  isEOF,
		// }

		err = sendData(writeBytes[:n], conn)
		if err != nil {
			return err
		}
	}

	return sendData(constants.SLine, conn)
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
