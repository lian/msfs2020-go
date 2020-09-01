package simconnect

import "fmt"

// MSFS-SDK/SimConnect\ SDK/include/SimConnect.h

const E_FAIL uint32 = 0x80004005

type DWORD uint32

const UNUSED DWORD = 0xffffffff // special value to indicate unused event, ID
const OBJECT_ID_USER DWORD = 0  // proxy value for User vehicle ObjectID

const (
	DATATYPE_INVALID      DWORD = iota // invalid data type
	DATATYPE_INT32                     // 32-bit integer number
	DATATYPE_INT64                     // 64-bit integer number
	DATATYPE_FLOAT32                   // 32-bit floating-point number (float)
	DATATYPE_FLOAT64                   // 64-bit floating-point number (double)
	DATATYPE_STRING8                   // 8-byte string
	DATATYPE_STRING32                  // 32-byte string
	DATATYPE_STRING64                  // 64-byte string
	DATATYPE_STRING128                 // 128-byte string
	DATATYPE_STRING256                 // 256-byte string
	DATATYPE_STRING260                 // 260-byte string
	DATATYPE_STRINGV                   // variable-length string
	DATATYPE_INITPOSITION              // see SIMCONNECT_DATA_INITPOSITION
	DATATYPE_MARKERSTATE               // see SIMCONNECT_DATA_MARKERSTATE
	DATATYPE_WAYPOINT                  // see SIMCONNECT_DATA_WAYPOINT
	DATATYPE_LATLONALT                 // see SIMCONNECT_DATA_LATLONALT
	DATATYPE_XYZ                       // see SIMCONNECT_DATA_XYZ

	DATATYPE_MAX // enum limit
)

const (
	TEXT_TYPE_SCROLL_BLACK DWORD = iota
	TEXT_TYPE_SCROLL_WHITE
	TEXT_TYPE_SCROLL_RED
	TEXT_TYPE_SCROLL_GREEN
	TEXT_TYPE_SCROLL_BLUE
	TEXT_TYPE_SCROLL_YELLOW
	TEXT_TYPE_SCROLL_MAGENTA
	TEXT_TYPE_SCROLL_CYAN
)

const (
	TEXT_TYPE_PRINT_BLACK DWORD = iota + 0x0100
	TEXT_TYPE_PRINT_WHITE
	TEXT_TYPE_PRINT_RED
	TEXT_TYPE_PRINT_GREEN
	TEXT_TYPE_PRINT_BLUE
	TEXT_TYPE_PRINT_YELLOW
	TEXT_TYPE_PRINT_MAGENTA
	TEXT_TYPE_PRINT_CYAN
)

const TEXT_TYPE_MENU DWORD = 0x0200

// Notification Group priority values
const GROUP_PRIORITY_HIGHEST DWORD = 1                 // highest priority
const GROUP_PRIORITY_HIGHEST_MASKABLE DWORD = 10000000 // highest priority that allows events to be masked
const GROUP_PRIORITY_STANDARD DWORD = 1900000000       // standard priority
const GROUP_PRIORITY_DEFAULT DWORD = 2000000000        // default priority
const GROUP_PRIORITY_LOWEST DWORD = 4000000000         // priorities lower than this will be ignored

func derefDataType(fieldType string) (DWORD, error) {
	var dataType DWORD
	switch fieldType {
	case "int32":
		dataType = DATATYPE_INT32
	case "int64":
		dataType = DATATYPE_INT64
	case "float32":
		dataType = DATATYPE_FLOAT32
	case "float64":
		dataType = DATATYPE_FLOAT64
	case "[8]byte":
		dataType = DATATYPE_STRING8
	case "[32]byte":
		dataType = DATATYPE_STRING32
	case "[64]byte":
		dataType = DATATYPE_STRING64
	case "[128]byte":
		dataType = DATATYPE_STRING128
	case "[256]byte":
		dataType = DATATYPE_STRING256
	case "[260]byte":
		dataType = DATATYPE_STRING260
	default:
		return 0, fmt.Errorf("DATATYPE not implemented: %s", fieldType)
	}

	return dataType, nil
}

const (
	RECV_ID_NULL DWORD = iota
	RECV_ID_EXCEPTION
	RECV_ID_OPEN
	RECV_ID_QUIT
	RECV_ID_EVENT
	RECV_ID_EVENT_OBJECT_ADDREMOVE
	RECV_ID_EVENT_FILENAME
	RECV_ID_EVENT_FRAME
	RECV_ID_SIMOBJECT_DATA
	RECV_ID_SIMOBJECT_DATA_BYTYPE
	RECV_ID_WEATHER_OBSERVATION
	RECV_ID_CLOUD_STATE
	RECV_ID_ASSIGNED_OBJECT_ID
	RECV_ID_RESERVED_KEY
	RECV_ID_CUSTOM_ACTION
	RECV_ID_SYSTEM_STATE
	RECV_ID_CLIENT_DATA
	RECV_ID_EVENT_WEATHER_MODE
	RECV_ID_AIRPORT_LIST
	RECV_ID_VOR_LIST
	RECV_ID_NDB_LIST
	RECV_ID_WAYPOINT_LIST
	RECV_ID_EVENT_MULTIPLAYER_SERVER_STARTED
	RECV_ID_EVENT_MULTIPLAYER_CLIENT_STARTED
	RECV_ID_EVENT_MULTIPLAYER_SESSION_ENDED
	RECV_ID_EVENT_RACE_END
	RECV_ID_EVENT_RACE_LAP

	RECV_ID_PICK
)

const (
	SIMOBJECT_TYPE_USER DWORD = iota
	SIMOBJECT_TYPE_ALL
	SIMOBJECT_TYPE_AIRCRAFT
	SIMOBJECT_TYPE_HELICOPTER
	SIMOBJECT_TYPE_BOAT
	SIMOBJECT_TYPE_GROUND
)

const (
	FACILITY_LIST_TYPE_AIRPORT DWORD = iota
	FACILITY_LIST_TYPE_WAYPOINT
	FACILITY_LIST_TYPE_NDB
	FACILITY_LIST_TYPE_VOR
	FACILITY_LIST_TYPE_COUNT // invalid
)

type Recv struct {
	Size    DWORD
	Version DWORD
	ID      DWORD
}

type RecvOpen struct {
	Recv
	ApplicationName         [256]byte
	ApplicationVersionMajor DWORD
	ApplicationVersionMinor DWORD
	ApplicationBuildMajor   DWORD
	ApplicationBuildMinor   DWORD
	SimConnectVersionMajor  DWORD
	SimConnectVersionMinor  DWORD
	SimConnectBuildMajor    DWORD
	SimConnectBuildMinor    DWORD
	Reserved1               DWORD
	Reserved2               DWORD
}

type RecvEvent struct {
	Recv
	//static const DWORD UNKNOWN_GROUP = DWORD_MAX;
	GroupID DWORD
	EventID DWORD
	Data    DWORD // uEventID-dependent context
}

type RecvSimobjectData struct {
	Recv
	RequestID   DWORD
	ObjectID    DWORD
	DefineID    DWORD
	Flags       DWORD // SIMCONNECT_DATA_REQUEST_FLAG
	entrynumber DWORD // if multiple objects returned, this is number <entrynumber> out of <outof>.
	outof       DWORD // note: starts with 1, not 0.
	DefineCount DWORD // data count (number of datums, *not* byte count)
	//SIMCONNECT_DATAV(   dwData, dwDefineID, ); // data begins here, dwDefineCount data items
}

type RecvSimobjectDataByType struct {
	RecvSimobjectData
}

type RecvException struct {
	Recv
	Exception DWORD // see SIMCONNECT_EXCEPTION
	//static const DWORD UNKNOWN_SENDID = 0;
	SendID DWORD // see SimConnect_GetLastSentPacketID
	//static const DWORD UNKNOWN_INDEX = DWORD_MAX;
	Index DWORD // index of parameter that was source of error
}

type RecvFacilityList struct {
	Recv
	RequestID   DWORD
	ArraySize   DWORD
	EntryNumber DWORD // when the array of items is too big for one send, which send this is (0..dwOutOf-1)
	OutOf       DWORD // total number of transmissions the list is chopped into
}

type RecvFacilityAirportList struct {
	RecvFacilityList
	List [1]DataFacilityAirport
}

type DataFacilityAirport struct {
	Icao      [9]byte // ICAO of the object
	Latitude  float64 // degrees
	Longitude float64 // degrees
	Altitude  float64 // meters
}

type RecvFacilityWaypointList struct {
	RecvFacilityList
	List [1]DataFacilityWaypoint
}

type DataFacilityWaypoint struct {
	DataFacilityAirport
	MagVar float64 // Magvar in degrees
}
