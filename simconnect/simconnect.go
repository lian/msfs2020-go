package simconnect

//go:generate go-bindata -pkg simconnect -o bindata.go -modtime 1 -prefix "../_vendor" "../_vendor/MSFS-SDK/SimConnect SDK/lib/SimConnect.dll"

// MSFS-SDK/SimConnect\ SDK/include/SimConnect.h
// MSFS-SDK/SimConnect\ SDK/lib/SimConnect.dll

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"unsafe"
)

var proc_SimConnect_Open *syscall.LazyProc
var proc_SimConnect_Close *syscall.LazyProc
var proc_SimConnect_AddToDataDefinition *syscall.LazyProc
var proc_SimConnect_SubscribeToSystemEvent *syscall.LazyProc
var proc_SimConnect_GetNextDispatch *syscall.LazyProc
var proc_SimConnect_RequestDataOnSimObjectType *syscall.LazyProc
var proc_SimConnect_SetDataOnSimObject *syscall.LazyProc

type SimConnect struct {
	handle    unsafe.Pointer
	DefineMap map[string]DWORD
}

func New(name string) (*SimConnect, error) {
	s := &SimConnect{
		DefineMap: map[string]DWORD{"_last": 0},
	}

	if proc_SimConnect_Open == nil {
		exePath, err := os.Executable()
		if err != nil {
			return nil, err
		}

		dllPath := filepath.Join(filepath.Dir(exePath), "SimConnect.dll")
		if _, err = os.Stat(dllPath); os.IsNotExist(err) {
			buf := MustAsset("MSFS-SDK/SimConnect SDK/lib/SimConnect.dll")

			if err := ioutil.WriteFile(dllPath, buf, 0644); err != nil {
				return nil, err
			}
		}

		mod := syscall.NewLazyDLL(dllPath)
		if err = mod.Load(); err != nil {
			return nil, err
		}

		proc_SimConnect_Open = mod.NewProc("SimConnect_Open")
		proc_SimConnect_Close = mod.NewProc("SimConnect_Close")
		proc_SimConnect_AddToDataDefinition = mod.NewProc("SimConnect_AddToDataDefinition")
		proc_SimConnect_SubscribeToSystemEvent = mod.NewProc("SimConnect_SubscribeToSystemEvent")
		proc_SimConnect_GetNextDispatch = mod.NewProc("SimConnect_GetNextDispatch")
		proc_SimConnect_RequestDataOnSimObjectType = mod.NewProc("SimConnect_RequestDataOnSimObjectType")
		proc_SimConnect_SetDataOnSimObject = mod.NewProc("SimConnect_SetDataOnSimObject")
	}

	// SimConnect_Open(
	//   HANDLE * phSimConnect,
	//   LPCSTR szName,
	//   HWND hWnd,
	//   DWORD UserEventWin32,
	//   HANDLE hEventHandle,
	//   DWORD ConfigIndex
	// );
	args := []uintptr{
		uintptr(unsafe.Pointer(&s.handle)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name))),
		0,
		0,
		0,
		0,
	}

	r1, _, err := proc_SimConnect_Open.Call(args...)
	if int32(r1) < 0 {
		return nil, fmt.Errorf("SimConnect_Open error: %d %s", int32(r1), err)
	}

	return s, nil
}

func (s *SimConnect) GetDefineID(a interface{}) DWORD {
	structName := reflect.TypeOf(a).Elem().Name()

	id, ok := s.DefineMap[structName]
	if !ok {
		id = s.DefineMap["_last"]
		s.DefineMap[structName] = id
		s.DefineMap["_last"] = id + 1
	}

	return id
}

func (s *SimConnect) RegisterDataDefinition(a interface{}) error {
	defineID := s.GetDefineID(a)

	v := reflect.ValueOf(a).Elem()
	for j := 1; j < v.NumField(); j++ {
		fieldName := v.Type().Field(j).Name
		nameTag, _ := v.Type().Field(j).Tag.Lookup("name")
		unitTag, _ := v.Type().Field(j).Tag.Lookup("unit")

		fieldType := v.Field(j).Kind().String()
		if fieldType == "array" {
			fieldType = fmt.Sprintf("[%d]byte", v.Field(j).Type().Len())
		}

		if nameTag == "" {
			return fmt.Errorf("%s name tag not found", fieldName)
		}

		dataType, err := derefDataType(fieldType)
		if err != nil {
			return err
		}

		s.AddToDataDefinition(defineID, nameTag, unitTag, dataType)
		//fmt.Printf("fieldName: %s  fieldType: %s  nameTag: %s unitTag: %s\n", fieldName, fieldType, nameTag, unitTag)
	}

	return nil
}

func (s *SimConnect) Close() error {
	// SimConnect_Open(
	//   HANDLE * phSimConnect,
	// );
	r1, _, err := proc_SimConnect_Close.Call(uintptr(s.handle))
	if int32(r1) < 0 {
		return fmt.Errorf("SimConnect_Close error: %d %s", int32(r1), err)
		panic("failed")
	}
	return nil
}

func (s *SimConnect) AddToDataDefinition(defineID DWORD, name, unit string, dataType DWORD) error {
	// SimConnect_AddToDataDefinition(
	//   HANDLE hSimConnect,
	//   SIMCONNECT_DATA_DEFINITION_ID DefineID,
	//   const char * DatumName,
	//   const char * UnitsName,
	//   SIMCONNECT_DATATYPE DatumType = SIMCONNECT_DATATYPE_FLOAT64,
	//   float fEpsilon = 0,
	//   DWORD DatumID = SIMCONNECT_UNUSED
	// );

	_name := []byte(name + "\x00")
	_unit := []byte(unit + "\x00")

	args := []uintptr{
		uintptr(s.handle),
		uintptr(defineID),
		uintptr(unsafe.Pointer(&_name[0])),
		uintptr(0),
		uintptr(dataType),
		uintptr(float32(0)),
		uintptr(UNUSED),
	}
	if unit != "" {
		args[3] = uintptr(unsafe.Pointer(&_unit[0]))
	}

	r1, _, err := proc_SimConnect_AddToDataDefinition.Call(args...)
	if int32(r1) < 0 {
		return fmt.Errorf("SimConnect_AddToDataDefinition for %s error: %d %s", name, r1, err)
	}

	return nil
}

func (s *SimConnect) SubscribeToSystemEvent(eventID DWORD, eventName string) error {
	// SimConnect_SubscribeToSystemEvent(
	//   HANDLE hSimConnect,
	//   SIMCONNECT_CLIENT_EVENT_ID EventID,
	//   const char * SystemEventName
	// );

	//_eventName := make([]byte, len(eventName)+1)
	//copy(_eventName, eventName)
	_eventName := []byte(eventName + "\x00")

	args := []uintptr{
		uintptr(s.handle),
		uintptr(eventID),
		uintptr(unsafe.Pointer(&_eventName[0])),
	}

	r1, _, err := proc_SimConnect_SubscribeToSystemEvent.Call(args...)
	if int32(r1) < 0 {
		return fmt.Errorf("SimConnect_SubscribeToSystemEvent for %s error: %d %s", eventName, r1, err)
	}

	return nil
}

func (s *SimConnect) RequestDataOnSimObjectType(requestID, defineID, radius, simobjectType DWORD) error {
	// SimConnect_RequestDataOnSimObjectType(
	//   HANDLE hSimConnect,
	//   SIMCONNECT_DATA_REQUEST_ID RequestID,
	//   SIMCONNECT_DATA_DEFINITION_ID DefineID,
	//   DWORD dwRadiusMeters,
	//   SIMCONNECT_SIMOBJECT_TYPE type
	// );
	args := []uintptr{
		uintptr(s.handle),
		uintptr(requestID),
		uintptr(defineID),
		uintptr(radius),
		uintptr(simobjectType),
	}

	r1, _, err := proc_SimConnect_RequestDataOnSimObjectType.Call(args...)
	if int32(r1) < 0 {
		return fmt.Errorf(
			"SimConnect_RequestDataOnSimObjectType for requestID %d defineID %d error: %d %s",
			requestID, defineID, r1, err,
		)
	}

	return nil
}

func (s *SimConnect) SetDataOnSimObject(defineID, simobjectType, flags, arrayCount, size DWORD, buf unsafe.Pointer) error {
	//s.SetDataOnSimObject(defineID, simconnect.OBJECT_ID_USER, 0, 0, size, buf)

	// SimConnect_SetDataOnSimObject(
	//   HANDLE hSimConnect,
	//   SIMCONNECT_DATA_DEFINITION_ID DefineID,
	//   SIMCONNECT_OBJECT_ID ObjectID,
	//   SIMCONNECT_DATA_SET_FLAG Flags,
	//   DWORD ArrayCount,
	//   DWORD cbUnitSize,
	//   void * pDataSet
	// );
	args := []uintptr{
		uintptr(s.handle),
		uintptr(defineID),
		uintptr(simobjectType),
		uintptr(flags),
		uintptr(arrayCount),
		uintptr(size),
		uintptr(buf),
	}

	r1, _, err := proc_SimConnect_SetDataOnSimObject.Call(args...)
	if int32(r1) < 0 {
		return fmt.Errorf(
			"SimConnect_SetDataOnSimObject for defineID %d error: %d %s",
			defineID, r1, err,
		)
	}

	return nil
}

func (s *SimConnect) GetNextDispatch() (unsafe.Pointer, int32, error) {
	var ppData unsafe.Pointer
	var ppDataLength DWORD

	r1, _, err := proc_SimConnect_GetNextDispatch.Call(
		uintptr(s.handle),
		uintptr(unsafe.Pointer(&ppData)),
		uintptr(unsafe.Pointer(&ppDataLength)),
	)

	return ppData, int32(r1), err
}
