package main

import (
	"fmt"
	"time"

	"github.com/lian/msfs2020-go/simconnect"
)

// ported from: MSFS-SDK/Samples/SimConnectSamples/RequestData/RequestData.cpp
// build: GOOS=windows GOARCH=amd64 go build github.com/lian/msfs2020-go/examples/request_data

type Report struct {
	simconnect.RecvSimobjectDataByType
	Title     [256]byte `name:"TITLE"`
	Kohlsman  float64   `name:"Kohlsman setting hg" unit:"inHg"`
	Altitude  float64   `name:"Plane Altitude" unit:"feet"`
	Latitude  float64   `name:"Plane Latitude" unit:"degrees"`
	Longitude float64   `name:"Plane Longitude" unit:"degrees"`
}

func (r *Report) RequestData(s *simconnect.SimConnect) {
	defineID := s.GetDefineID(r)
	requestID := defineID
	s.RequestDataOnSimObjectType(requestID, defineID, 0, simconnect.SIMOBJECT_TYPE_USER)
}

func main() {
	s, err := simconnect.New("Request Data")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to Flight Simulator!")

	report := &Report{}
	s.RegisterDataDefinition(report)
	report.RequestData(s)

	/*
		fmt.Println("SubscribeToSystemEvent")
		eventSimStartID := simconnect.DWORD(0)
		s.SubscribeToSystemEvent(eventSimStartID, "SimStart")
	*/

	for {
		ppData, r1, err := s.GetNextDispatch()

		if r1 < 0 {
			if uint32(r1) == simconnect.E_FAIL {
				// skip error, means no new messages?
				continue
			} else {
				panic(fmt.Errorf("GetNextDispatch error: %d %s", r1, err))
			}
		}

		recvInfo := *(*simconnect.Recv)(ppData)
		//fmt.Println(ppData, pcbData, recvInfo)

		switch recvInfo.ID {
		case simconnect.RECV_ID_EXCEPTION:
			recvErr := *(*simconnect.RecvException)(ppData)
			fmt.Printf("SIMCONNECT_RECV_ID_EXCEPTION %#v\n", recvErr)

		case simconnect.RECV_ID_OPEN:
			recvOpen := *(*simconnect.RecvOpen)(ppData)
			fmt.Println("SIMCONNECT_RECV_ID_OPEN", fmt.Sprintf("%s", recvOpen.ApplicationName))
			//spew.Dump(recvOpen)
		case simconnect.RECV_ID_EVENT:
			recvEvent := *(*simconnect.RecvEvent)(ppData)
			fmt.Println("SIMCONNECT_RECV_ID_EVENT")
			//spew.Dump(recvEvent)

			switch recvEvent.EventID {
			//case eventSimStartID:
			//	fmt.Println("SimStart Event")
			default:
				fmt.Println("unknown SIMCONNECT_RECV_ID_EVENT", recvEvent.EventID)
			}

		case simconnect.RECV_ID_SIMOBJECT_DATA_BYTYPE:
			recvData := *(*simconnect.RecvSimobjectDataByType)(ppData)
			fmt.Println("SIMCONNECT_RECV_SIMOBJECT_DATA_BYTYPE")

			switch recvData.RequestID {
			case s.DefineMap["Report"]:
				report := (*Report)(ppData)
				fmt.Printf("REPORT: %s: GPS: %.6f,%.6f Altitude: %.0f\n", report.Title, report.Latitude, report.Longitude, report.Altitude)
				report.RequestData(s)
			}

		default:
			fmt.Println("recvInfo.dwID unknown", recvInfo.ID)
		}

		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("close")

	if err = s.Close(); err != nil {
		panic(err)
	}
}
