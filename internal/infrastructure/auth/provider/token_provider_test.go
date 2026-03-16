package token_provider

import (
	"context"
	"errors"
	"testing"

	token_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
	. "github.com/smartystreets/goconvey/convey"
)

type contextKey string

const requestIDKey contextKey = "request-id"

type tokenClientSpy struct {
	gotCtx   context.Context
	gotInput *token_model.RefreshTokenInput
	output   *token_model.RefreshTokenOutput
	err      error
}

func (s *tokenClientSpy) RefreshToken(ctx context.Context, input *token_model.RefreshTokenInput) (*token_model.RefreshTokenOutput, error) {
	s.gotCtx = ctx
	s.gotInput = input
	if s.err != nil {
		return nil, s.err
	}
	return s.output, nil
}

func TestShouldReturnAccessTokenAndPassContextAndInputWhenRefreshSucceeds(t *testing.T) {
	Convey("Given a token provider with valid refresh input and token client response", t, func() {
		input := &token_model.RefreshTokenInput{RefreshToken: "refresh"}
		spy := &tokenClientSpy{output: &token_model.RefreshTokenOutput{AccessToken: "access-token"}}
		provider := NewTokenProvider(input, spy)
		ctx := context.WithValue(context.Background(), requestIDKey, "abc")

		Convey("When GetValidToken is called", func() {
			token, err := provider.GetValidToken(ctx)

			Convey("Then it should return access token and forward context and input", func() {
				So(err, ShouldBeNil)
				So(token, ShouldEqual, "access-token")
				So(spy.gotCtx, ShouldEqual, ctx)
				So(spy.gotInput, ShouldEqual, input)
			})
		})
	})
}

func TestShouldReturnErrorAndEmptyTokenWhenRefreshFails(t *testing.T) {
	Convey("Given a token provider with token client refresh error", t, func() {
		refreshErr := errors.New("refresh failed")
		provider := NewTokenProvider(&token_model.RefreshTokenInput{}, &tokenClientSpy{err: refreshErr})

		Convey("When GetValidToken is called", func() {
			token, err := provider.GetValidToken(context.Background())

			Convey("Then it should return the refresh error and an empty token", func() {
				So(token, ShouldBeBlank)
				So(errors.Is(err, refreshErr), ShouldBeTrue)
			})
		})
	})
}
