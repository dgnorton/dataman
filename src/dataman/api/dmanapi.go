package dataman

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func FindDevices(addr string, duration time.Duration) ([]string, error) {
	wconn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}
	defer wconn.Close()

	_, err = wconn.Write([]byte("CGNM\x04@\x02\x00\x00"))

	rconn, err := net.ListenPacket("udp", "0.0.0.0:1069")
	defer rconn.Close()

	var buf [2048]byte
	t := time.Now().Add(duration)
	rconn.SetReadDeadline(t)
	var deviceIPs []string
	for time.Now().Before(t) {
		_, a, err := rconn.ReadFrom(buf[:])
		if err != nil {
			break
		}
		ip := strings.Split(a.String(), ":")[0] + ":23"
		deviceIPs = append(deviceIPs, ip)
	}

	return deviceIPs, err
}

type DevConn struct {
	Sock   net.Conn
	CmdSeq int
}

func Open(addr string) (*DevConn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &DevConn{conn, 1}, nil
}

func (c *DevConn) Close() {
	c.Close()
}

func (c *DevConn) GetType() (string, error) {
	cmd, _ := c.NewCmd("GET DEVICE.TYPE", ChkSum_Yes, CmdID_Yes)
	resp, err := c.Write(cmd)
	return string(resp), err
}

func (c *DevConn) GetName() (string, error) {
	cmd, _ := c.NewCmd("GET DEVICE.NAME", ChkSum_Yes, CmdID_Yes)
	resp, err := c.Write(cmd)
	return string(resp), err
}

func (c *DevConn) GetSerialNumber() (string, error) {
	cmd, _ := c.NewCmd("GET DEVICE.SERIAL-NUMBER", ChkSum_Yes, CmdID_Yes)
	resp, err := c.Write(cmd)
	return string(resp), err
}

func (c *DevConn) GetFirmwareVersion() (string, error) {
	cmd, _ := c.NewCmd("GET DEVICE.FIRMWARE-VER", ChkSum_Yes, CmdID_Yes)
	resp, err := c.Write(cmd)
	return string(resp), err
}

func (c *DevConn) GetFeatureKeys() (string, error) {
	cmd, _ := c.NewCmd("GET DEVICE.FEATURE-KEYS", ChkSum_Yes, CmdID_Yes)
	resp, err := c.Write(cmd)
	return string(resp), err
}

func (c *DevConn) Beep(cnt int, vol int) error {
	cmdstr := fmt.Sprintf("BEEP %v %v", cnt, vol)
	cmd, err := c.NewCmd(cmdstr, ChkSum_No, CmdID_No)
	if err != nil {
		return err
	}
	_, err = c.Write(cmd)

	return err
}

func (c *DevConn) GetResult() (string, error) {
	cmd, _ := c.NewCmd("GET RESULT", ChkSum_Yes, CmdID_Yes)
	resp, err := c.Write(cmd)
	return string(resp), err
}

func (c *DevConn) Read() ([]byte, error) {
        var buf [2048]byte
        n, e := c.Sock.Read(buf[:])
        return buf[:n], e
}

func (c *DevConn) Write(cmd *Cmd) (resp []byte, err error) {
	_, e := c.Sock.Write(cmd.Bytes())
	if e != nil {
		return nil, e
	} else if cmd.CmdID <= 0 {
		return nil, nil
	}

	c.Sock.SetReadDeadline(time.Now().Add(time.Second*3))

	resp = make([]byte, 4096)
	n, e := c.Sock.Read(resp[:])
	if e != nil {
		return nil, e
	}

	return resp[:n], nil
}

type TrigType int

const (
        Trig_Single TrigType = iota
        Trig_Presentation
        Trig_Manual
        Trig_Burst
        Trig_Self
        Trig_Continuous
        Trig_Unknown
)

func TriggerType(b []byte) TrigType {
        switch string(b) {
                case "0":
                        return Trig_Single
                case "1":
                        return Trig_Presentation
                case "2":
                        return Trig_Manual
                case "3":
                        return Trig_Burst
                case "4":
                        return Trig_Self
                case "5":
                        return Trig_Continuous
        }

        return Trig_Unknown
}

func (t *TrigType) String() string {
        switch *t {
                case Trig_Single:
                        return "0"
                case Trig_Presentation:
                        return "1"
                case Trig_Manual:
                        return "2"
                case Trig_Burst:
                        return "3"
                case Trig_Self:
                        return "4"
                case Trig_Continuous:
                        return "5"
        }

        return ""
}

func (c *DevConn) GetTriggerType() (TrigType, error) {
        cmd, _ := c.NewCmd("GET TRIGGER.TYPE", ChkSum_Yes, CmdID_Yes)
        resp, err := c.Write(cmd)
        return TriggerType(resp), err
}

func (c *DevConn) SetTriggerType(t TrigType) error {
        cmdstr := fmt.Sprintf("SET TRIGGER.TYPE %v", t.String())
        cmd, _ := c.NewCmd(cmdstr, ChkSum_Yes, CmdID_Yes)
        _, err := c.Write(cmd)

        return err
}

type ChkSumFlag int

const (
	ChkSum_No ChkSumFlag = iota
	ChkSum_Yes
)

type CmdIDFlag int

const (
	CmdID_No CmdIDFlag = iota
	CmdID_Yes
)

type Cmd struct {
	CmdID  int
	Str    string
	ChkSum byte
}

func (cmd *Cmd) Bytes() []byte {
	return []byte(cmd.Str)
}

func (c *DevConn) NewCmd(cmdstr string, useChkSum ChkSumFlag, useCmdID CmdIDFlag) (*Cmd, error) {
	var strs []string
	cmd := &Cmd{0, "", 0}

	// build command header
	strs = append(strs, "||")

	if useChkSum == ChkSum_Yes {
		strs = append(strs, "1")
	} else if useChkSum == ChkSum_No {
		strs = append(strs, "0")
	} else {
		return nil, fmt.Errorf("invalid parameter useChkSum = %v", useChkSum)
	}

	if useCmdID == CmdID_Yes {
		cmd.CmdID = c.CmdSeq
		c.CmdSeq++
		strs = append(strs, ":", strconv.Itoa(cmd.CmdID))
	} else if useCmdID == CmdID_No {
		// do nothing
	} else {
		return nil, fmt.Errorf("invalid parameter useCmdId = %v", useCmdID)
	}

	strs = append(strs, ">")

	// add command
	strs = append(strs, cmdstr)

	// add checksum, if requested
	if useChkSum == ChkSum_Yes {
		cmd.ChkSum = CheckSum([]byte(cmdstr))
		strs = append(strs, string(cmd.ChkSum))
	}

	// add footer
	strs = append(strs, "\r\n")

	// combine all the command parts 
	cmd.Str = strings.Join(strs, "")

	return cmd, nil
}

func CheckSum(buf []byte) byte {
	var chksum byte = 0
	for _, b := range buf {
		chksum ^= b
	}

	return chksum
}
