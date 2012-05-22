package main

import (
	"dataman/api"
	"fmt"
	"time"
        "net/http"
)

var deviceIPs []string
var dm *dataman.DevConn

func beepHandler(w http.ResponseWriter, r *http.Request) {
        dm.Beep(3, 2)
        fmt.Fprintf(w, "<html><head/><body>Beep sent</body></html>")
}

func main() {
	deviceIPs, err := dataman.FindDevices("192.168.0.255:1069", 1*time.Second)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Device found...")
	fmt.Println(deviceIPs)

	dm, err = dataman.Open("127.0.0.1:2300")
	//dm, err = dataman.Open(deviceIPs[0])

	if err != nil {
		fmt.Println("error connecting")
		return
	}

	//fmt.Println(dm.GetType())
	//fmt.Println(dm.GetName())
	//fmt.Println(dm.GetSerialNumber())
	//fmt.Println(dm.GetFirmwareVersion())
	//fmt.Println(dm.GetFeatureKeys())
	fmt.Println(dm.SetTriggerType(dataman.Trig_Manual))
   for {
      b, e := dm.Read()
      if e == nil {
         fmt.Println(string(b))
      }
   }
	dm.Beep(3, 2)

	//cmd, _ := dm.NewCmd("SET SYMBOL.UPC-EAN ON", dataman.ChkSum_Yes, dataman.CmdID_Yes)
	//resp, err := dm.Write(cmd)
	//fmt.Printf("%v %v\n", string(resp), err)

	//cmd, _ = dm.NewCmd("SET SYMBOL.UPC-EAN.EXPANDED ON", dataman.ChkSum_Yes, dataman.CmdID_Yes)
	//resp, err = dm.Write(cmd)

	//fmt.Println(dm.GetTriggerType())
	//fmt.Println(dm.SetTriggerType(dataman.Trig_Manual))

        http.HandleFunc("/beep/", beepHandler)
        http.ListenAndServe("localhost:9090", nil)
}
