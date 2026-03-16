package token_model

import (
	"encoding/json"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldSetAccessTokenAndExpirationWhenNewTokenIsCalled(t *testing.T) {
	Convey("Given an access token and expires in seconds", t, func() {
		before := time.Now()
		token := NewToken("access-token", 120)
		after := time.Now()

		Convey("When NewToken is called", func() {
			minExpected := before.Add(119 * time.Second)
			maxExpected := after.Add(121 * time.Second)

			Convey("Then it should expose access token and expected expiration range", func() {
				So(token.AccessToken(), ShouldEqual, "access-token")
				So(token.ExpiresAt().Before(minExpected), ShouldBeFalse)
				So(token.ExpiresAt().After(maxExpected), ShouldBeFalse)
			})
		})
	})
}

func TestShouldReturnTrueWhenIsValidCalledWithTokenNotNearExpiry(t *testing.T) {
	Convey("Given a token with access token and enough remaining lifetime", t, func() {
		token := NewToken("access-token", 120)

		Convey("When IsValid is called", func() {
			Convey("Then it should return true", func() {
				So(token.IsValid(), ShouldBeTrue)
			})
		})
	})
}

func TestShouldReturnFalseWhenIsValidCalledWithEmptyAccessToken(t *testing.T) {
	Convey("Given a token with empty access token", t, func() {
		token := NewToken("", 120)

		Convey("When IsValid is called", func() {
			Convey("Then it should return false", func() {
				So(token.IsValid(), ShouldBeFalse)
			})
		})
	})
}

func TestShouldReturnFalseWhenIsValidCalledWithTokenNearExpiry(t *testing.T) {
	Convey("Given a token expiring within the 30 second buffer", t, func() {
		token := NewToken("access-token", 20)

		Convey("When IsValid is called", func() {
			Convey("Then it should return false", func() {
				So(token.IsValid(), ShouldBeFalse)
			})
		})
	})
}

func TestShouldPopulateRefreshTokenOutputFieldsWhenJsonUnmarshalSucceeds(t *testing.T) {
	Convey("Given valid refresh token output json", t, func() {
		var output RefreshTokenOutput
		data := []byte(`{"access_token":"token-value","expires_in":3600}`)

		Convey("When json.Unmarshal is called", func() {
			err := json.Unmarshal(data, &output)

			Convey("Then it should decode access token and expires_in", func() {
				So(err, ShouldBeNil)
				So(output.AccessToken, ShouldEqual, "token-value")
				So(output.ExpiresIn, ShouldEqual, 3600)
			})
		})
	})
}
