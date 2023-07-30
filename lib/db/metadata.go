package db

// RotaSchedule is the type that defines how the members of a rota are scheduled
type RotaSchedule string

// RotaFrequency is the type that defines how long a rota lasts
type RotaFrequency string

const (
	RFDaily   = RotaFrequency("Daily")
	RFWeekly  = RotaFrequency("Weekly")
	RFMonthly = RotaFrequency("Monthly")

	RSCreated = RotaSchedule("Created At")
	RSRandom  = RotaSchedule("Randomly")
)

type RotaMetadata struct {
	Frequency      RotaFrequency `json:"frequency"`
	SchedulingType RotaSchedule  `json:"scheduling_type"`
}
