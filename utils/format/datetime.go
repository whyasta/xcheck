package format

import "time"

type Datetime struct {
	time.Time
}

type Date struct {
	time.Time
}

func (t *Datetime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02T15:04:05.000-0700"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}

func (t *Date) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}
