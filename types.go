package main

// Site represents an IQ4 site (location grouping controllers).
type Site struct {
	ID                 int    `json:"id"`
	CompanyID          int    `json:"companyId"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	TimeZone           string `json:"timeZone"`
	PostalCode         string `json:"postalCode"`
	NumberOfSatellites int    `json:"numberOfSatellites"`
}

// Controller represents an IQ4 satellite/controller.
type Controller struct {
	ID                   int     `json:"id"`
	Name                 string  `json:"name"`
	Type                 int     `json:"type"`
	SiteID               int     `json:"siteId"`
	SiteName             string  `json:"siteName"`
	StationCount         int     `json:"stationCount"`
	ProgramsUsedCount    int     `json:"programsUsedCount"`
	MacAddress           string  `json:"macAddress"`
	Version              string  `json:"version"`
	IsMQTT               bool    `json:"isMQTT"`
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	RainDelay            int     `json:"rainDelay"`
	SimultaneousStations int     `json:"simultaneousStations"`
	IsShutdown           bool    `json:"isShutdown"`
	SatelliteEnabled     bool    `json:"satelliteEnabled"`
}

// ConnectionStatus represents the MQTT connection state of a controller.
type ConnectionStatus struct {
	ID          int    `json:"id"`
	IsConnected bool   `json:"isConnected"`
	Timestamp   string `json:"timestamp"`
}

// ConnectionStatusResponse wraps the isConnected endpoint response.
type ConnectionStatusResponse struct {
	Satellites []ConnectionStatus `json:"satellites"`
}

// Program represents an irrigation program on a controller.
type Program struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	ShortName              string `json:"shortName"`
	ProgramAdjust          int    `json:"programAdjust"`
	Type                   int    `json:"type"`
	SatelliteID            int    `json:"satelliteId"`
	SatelliteName          string `json:"satelliteName,omitempty"`
	SiteID                 int    `json:"siteId,omitempty"`
	SiteName               string `json:"siteName,omitempty"`
	WeekDays               string `json:"weekDays"`
	HybridWeekDays         string `json:"hybridWeekDays,omitempty"`
	SkipDays               int    `json:"skipDays"`
	Number                 int    `json:"number"`
	IsEnabled              bool   `json:"isEnabled"`
	NumberOfProgramSteps   int    `json:"numberOfProgramSteps"`
	StationDelay           string `json:"stationDelay"`
	SimultaneousStations   int    `json:"simultaneousStations,omitempty"`
	ETAdjustType           int    `json:"etAdjustType,omitempty"`
	ApplySeasonalByMonth   bool   `json:"applySeasonalAdjustByMonth,omitempty"`
}

// ProgramDetail is the full program object from GetProgram (used for updates).
type ProgramDetail map[string]any

// StartTime represents a program start time.
type StartTime struct {
	ID             int    `json:"id"`
	DateTime       string `json:"dateTime"`
	DateTimeOffset string `json:"dateTimeOffset,omitempty"`
	ProgramID      int    `json:"programId"`
	Enabled        bool   `json:"enabled"`
	CompanyID      int    `json:"companyId,omitempty"`
	LastEditUserID int    `json:"lastEditUserId,omitempty"`
	AssetID        int    `json:"assetID,omitempty"`
	DateTimeLocal  string `json:"dateTimeLocal,omitempty"`
	DateTimeUTC    string `json:"dateTimeUTC,omitempty"`
}

// Station represents a physical valve zone on a controller.
type Station struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	SatelliteID       int      `json:"satelliteId"`
	Terminal          int      `json:"terminal"`
	YearlyAdjFactor   float64  `json:"yearlyAdjFactor"`
	ETAdjustFactor    int      `json:"etAdjustFactor"`
	Priority          int      `json:"priority"`
	LandscapeType     string   `json:"areaLevel2Name,omitempty"`
	SprinklerType     string   `json:"areaLevel3Name,omitempty"`
	Arc               *float64 `json:"arc,omitempty"`
	SoilType          *string  `json:"soilType,omitempty"`
	PrecRate          *float64 `json:"precRateFinal,omitempty"`
	CropCoefficient   float64  `json:"cropCoefficient,omitempty"`
	Slope             int      `json:"slope,omitempty"`
}

// StationRuntime maps a station to its program assignments.
type StationRuntime struct {
	StationID                int              `json:"stationId"`
	RuntimeProgramAssignedList []RuntimeAssignment `json:"runtimeProgramAssignedList"`
}

// RuntimeAssignment is a single program step assigned to a station.
type RuntimeAssignment struct {
	ProgramStepID    int    `json:"programStepId"`
	ProgramID        int    `json:"programId"`
	ProgramShortName string `json:"programShortName"`
	AdjustedRunTime  string `json:"adjustedRunTime"`
	BaseRunTime      string `json:"baseRunTime"`
	ETAdjustType     int    `json:"etAdjustType"`
	ETAdjustFactor   int    `json:"etAdjustFactor"`
}

// ProgramStep is the full step detail from GetProgramStepById (used for updates).
type ProgramStep struct {
	ID             int    `json:"id"`
	ProgramID      int    `json:"programId"`
	StationID      int    `json:"stationId"`
	SequenceNumber int    `json:"sequenceNumber"`
	ActionID       int    `json:"actionId"`
	RunTime        string `json:"runTime"`
	RunTimeLong    int64  `json:"runTimeLong"`
	CompanyID      int    `json:"companyId"`
	AssetID        int    `json:"assetID"`
}

// NewProgramStep is the payload format for CreateProgramSteps (matches IQ4 web UI).
type NewProgramStep struct {
	ActionID   string `json:"actionId"`
	ProgramID  string `json:"programId"`
	RunTimeLong *int64 `json:"runTimeLong"`
	StationID  int    `json:"stationId"`
}
