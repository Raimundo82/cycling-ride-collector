package usecase

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseDateToUnix(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
		want    int64
	}{
		{"ValidDate", "2026-10-01", false, time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC).Unix()},
		{"EmptyString", "", true, 0},
		{"BadFormat", "01-10-2026", true, 0},
		{"LeapDay", "2024-02-29", false, time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC).Unix()},
	}

	Convey("ParseDateToUnix table-driven cases", t, func() {
		for _, tc := range cases {
			Convey(tc.name, func() {
				ts, err := ParseDateToUnix(tc.input)
				if tc.wantErr {
					So(err, ShouldNotBeNil)
				} else {
					So(err, ShouldBeNil)
					So(ts, ShouldEqual, tc.want)
				}
			})
		}
	})
}
