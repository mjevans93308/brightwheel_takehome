package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mjevans93308/brightwheel_takehome/internal/models"
)

// maintain an in-memory mapping of every device we receive
// device UUID -> map of timestamp -> count

var (
	Devices          = make(map[string]models.Device)
	TimeFormatLayout = "2006-01-02T15:04:05-07:00"
)

type ReadingReq struct {
	Timestamp string `json:"timestamp"`
	Count     int    `json:"count"`
}

type DeviceRequest struct {
	ID       string       `json:"id"`
	Readings []ReadingReq `json:"readings"`
}

func DeviceHandler(w http.ResponseWriter, r *http.Request) {
	// kick req out if content-type header not set
	if r.Header.Get("Content-Type") == "" || r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	// enforce a max size for our json payload of 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	deviceReq := DeviceRequest{}
	if err := dec.Decode(&deviceReq); err != nil {
		localErr := fmt.Sprintf("Error: %s\n", err)
		fmt.Printf(localErr)
		http.Error(w, localErr, http.StatusInternalServerError)
		return
	}

	// check whether we have processed this device yet
	existingDevice, matchedDevice := Devices[deviceReq.ID]

	readings := make(map[time.Time]int)
	device := models.Device{}
	latestReading := &models.Reading{}
	var totalCount int
	var notAllReadingsProcessed bool
	for _, rReq := range deviceReq.Readings {
		parsedTime, err := time.Parse(TimeFormatLayout, rReq.Timestamp)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			// jump to next reading
			notAllReadingsProcessed = true
			continue
		}
		// if we already have an entry for that timestamp for that device
		// we want to ignore any subsequent entries for it
		if _, ok := existingDevice.ReadingsMap[parsedTime]; ok && matchedDevice {
			continue
		}
		if latestReading == nil {
			latestReading = &models.Reading{
				Timestamp: parsedTime,
				Count:     rReq.Count,
			}
		} else {
			// compare latestReading's timestamp with current reading timestamp
			if parsedTime.After(latestReading.Timestamp) {
				latestReading.Timestamp = parsedTime
				latestReading.Count = rReq.Count
			}
		}
		readings[parsedTime] = rReq.Count
		totalCount += rReq.Count
	}

	// update the latest reading if necessary and if we have a device already
	if matchedDevice {
		if latestReading.Timestamp.After(existingDevice.LatestReading.Timestamp) {
			existingDevice.LatestReading = *latestReading
		}
		existingDevice.TotalCount += totalCount
		// add new readings into existing device's readings map
		for k, v := range existingDevice.ReadingsMap {
			existingDevice.ReadingsMap[k] = v
		}
	} else {
		device.TotalCount = totalCount
		device.ReadingsMap = readings
		device.ID = deviceReq.ID
		device.LatestReading = *latestReading

		// store device into in-mem mapping
		Devices[device.ID] = device
	}

	printDeviceState()

	if notAllReadingsProcessed {
		http.Error(w, "Could not process all readings", http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func printDeviceState() {
	data, _ := json.Marshal(Devices)
	fmt.Printf("data: %s\n", data)
}

type LatestTimestampResp struct {
	LatestTimestamp string `json:"latest_timestamp"`
}

func LatestTimestampHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("deviceId")
	if deviceID == "" {
		localErr := fmt.Sprintf("Error: %s\n", errors.New("no deviceID value found in query params"))
		fmt.Printf(localErr)
		http.Error(w, localErr, http.StatusBadRequest)
		return
	}
	device, ok := Devices[deviceID]
	if !ok {
		localErr := fmt.Sprintf("Error: %s\n", errors.New("no device with that deviceID found"))
		fmt.Printf(localErr)
		http.Error(w, localErr, http.StatusNotFound)
		return
	}
	var latestTimestampResp LatestTimestampResp
	latestTimestampResp.LatestTimestamp = device.LatestReading.Timestamp.Format(TimeFormatLayout)

	jsonData, err := json.Marshal(latestTimestampResp)
	if err != nil {
		localErr := fmt.Sprintf("Error: %s\n", errors.New("could not marshal latest timestamp resp to json"))
		fmt.Printf(localErr)
		http.Error(w, localErr, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

type CumulativeResp struct {
	CumulativeCount int `json:"cumulative_count"`
}

func CumulativeHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("deviceId")
	if deviceID == "" {
		localErr := fmt.Sprintf("Error: %s\n", errors.New("no deviceID value found in query params"))
		fmt.Printf(localErr)
		http.Error(w, localErr, http.StatusBadRequest)
		return
	}
	device, ok := Devices[deviceID]
	if !ok {
		localErr := fmt.Sprintf("Error: %s\n", errors.New("no device with that deviceID found"))
		fmt.Printf(localErr)
		http.Error(w, localErr, http.StatusNotFound)
		return
	}
	cumulativeResp := CumulativeResp{
		CumulativeCount: device.TotalCount,
	}

	jsonData, err := json.Marshal(cumulativeResp)
	if err != nil {
		localErr := fmt.Sprintf("Error: %s\n", errors.New("could not marshal cumulative count resp to json"))
		fmt.Printf(localErr)
		http.Error(w, localErr, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
