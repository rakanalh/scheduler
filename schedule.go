package scheduler

type Unit int

const (
	Second Unit = iota
	Minute
	Hour
	Day
	WorkDay
)

type OneTimeSchedule struct {
}

type StaticSchedule struct {
	hour   int
	minute int
	number int
	unit   Unit
}

type DynamicSchedule struct {
	number int
	unit   Unit
}

func At(hour int, minute int) *StaticSchedule {
	return &StaticSchedule{
		hour:   hour,
		minute: minute,
	}
}

func Every(number int) DynamicSchedule {
	return DynamicSchedule{
		number: number,
	}
}

func (schedule *StaticSchedule) Every(number int) *StaticSchedule {
	schedule.number = number
	return schedule
}

func (schedule *StaticSchedule) Day() {
	schedule.unit = Day
}

func (schedule *StaticSchedule) WorkDay() {
	schedule.unit = WorkDay
}

// At 10:30 every day
// At 11:00 every monday
// Every 30 seconds
// Every 5 minutes
// Every 10 hours
func (schedule *DynamicSchedule) Second() {
	schedule.unit = Second
}

func (schedule *DynamicSchedule) Minute() {
	schedule.unit = Minute
}

func (schedule *DynamicSchedule) Hour() {
	schedule.unit = Hour
}

func (schedule *DynamicSchedule) Day() {
	schedule.unit = Day
}

func (schedule *DynamicSchedule) WeekDay() {
	schedule.unit = WorkDay
}
