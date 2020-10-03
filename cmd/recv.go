package cmd

import (
	"bytes"
	"encoding/binary"
	"ft/constants"
	"io"
	"net"
	"os"

	"github.com/spf13/cobra"
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

		f, err := os.OpenFile(savePath, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}

		defer f.Close()

		for {
			// revice size
			sizeBytes := make([]byte, 8)
			_, err = io.ReadFull(conn, sizeBytes)
			if err != nil {
				return err
			}
			size := binary.BigEndian.Uint64(sizeBytes)

			// recv data
			data := make([]byte, size)
			n, err := io.ReadFull(conn, data)
			if err != nil {
				return nil
			}

			if size == constants.DefaultSLenSize && bytes.Compare(data, constants.SLine) == 0 {
				break
			}

			f.Write(data[:n])
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(recvCmd)

	recvCmd.Flags().StringVarP(&savePath, "path", "p", "./fss", "the save dir path")
	recvCmd.Flags().StringVarP(&serverAddr, "addr", "a", ":9990", "server connect addr")
}
