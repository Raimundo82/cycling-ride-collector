package model

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsExpiredReturnFalseIfTokenIsNotExpired(t *testing.T) {
	Convey("Given a token that expires in more than one minute", t, func() {
		token := &Token{
			AccessToken:  "test-access-token",
			ExpiresAt:    time.Now().Add(2 * time.Minute),
			RefreshToken: "test-refresh-token",
		}

		Convey("Then it returns the expected token", func() {
			So(token.IsExpired(), ShouldBeFalse)
		})
	})
}

func TestIsExpiredReturnFalseIfTokenIsExpiringInOneMinute(t *testing.T) {
	Convey("Given a token that expires in exactly one minute one second", t, func() {
		token := &Token{
			AccessToken:  "test-access-token",
			ExpiresAt:    time.Now().Add(1*time.Minute + 1*time.Second),
			RefreshToken: "test-refresh-token",
		}

		Convey("Then it returns the expected token", func() {
			So(token.IsExpired(), ShouldBeFalse)
		})
	})
}

func TestIsExpiredReturnTrueIfTokenIsExpiringInLessThanOneMinute(t *testing.T) {
	Convey("Given a token that expires in less than one minute", t, func() {
		token := &Token{
			AccessToken:  "test-access-token",
			ExpiresAt:    time.Now().Add(30 * time.Second),
			RefreshToken: "test-refresh-token",
		}

		Convey("Then it returns the expected token", func() {
			So(token.IsExpired(), ShouldBeTrue)
		})
	})
}

func TestIsExpiredReturnTrueIfTokenIsAboutToExpire(t *testing.T) {
	Convey("Given a token that is about to expire", t, func() {
		token := &Token{
			AccessToken:  "test-access-token",
			ExpiresAt:    time.Now(),
			RefreshToken: "test-refresh-token",
		}

		Convey("Then it returns the expected token", func() {
			So(token.IsExpired(), ShouldBeTrue)
		})
	})
}

func TestIsExpiredReturnTrueIfTokenIsZero(t *testing.T) {
	Convey("Given a token with zero expiration time", t, func() {
		token := &Token{
			AccessToken:  "test-access-token",
			ExpiresAt:    time.Time{},
			RefreshToken: "test-refresh-token",
		}

		Convey("Then it returns the expected token", func() {
			So(token.IsExpired(), ShouldBeTrue)
		})
	})
}
