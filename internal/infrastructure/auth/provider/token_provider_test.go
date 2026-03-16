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
	called   int
	gotCtx   context.Context
	gotInput *token_model.RefreshTokenInput
	output   *token_model.RefreshTokenOutput
	err      error
}

func (s *tokenClientSpy) RefreshToken(ctx context.Context, input *token_model.RefreshTokenInput) (*token_model.RefreshTokenOutput, error) {
	s.called++
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
		spy := &tokenClientSpy{output: &token_model.RefreshTokenOutput{AccessToken: "access-token", ExpiresIn: 3600}}
		provider := NewTokenProvider(input, spy)
		ctx := context.WithValue(context.Background(), requestIDKey, "abc")

		Convey("When GetValidToken is called", func() {
			token, err := provider.GetValidToken(ctx)

			Convey("Then it should return access token and forward context and input", func() {
				So(err, ShouldBeNil)
				So(token, ShouldEqual, "access-token")
				So(spy.called, ShouldEqual, 1)
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
				So(provider.(*tokenProvider).cachedToken, ShouldResemble, token_model.Token{})
			})
		})
	})
}

func TestShouldReturnCachedTokenOnSecondCallWithoutRefreshing(t *testing.T) {
	Convey("Given a token provider that already obtained a token", t, func() {
		spy := &tokenClientSpy{output: &token_model.RefreshTokenOutput{AccessToken: "access-token", ExpiresIn: 3600}}
		provider := NewTokenProvider(&token_model.RefreshTokenInput{}, spy)
		_, _ = provider.GetValidToken(context.Background())

		Convey("When GetValidToken is called again before expiry", func() {
			token, err := provider.GetValidToken(context.Background())

			Convey("Then it should return the cached token without calling the client again", func() {
				So(err, ShouldBeNil)
				So(token, ShouldEqual, "access-token")
				So(spy.called, ShouldEqual, 1)
			})
		})
	})
}

func TestShouldRefreshTokenWhenCachedTokenIsExpired(t *testing.T) {
	Convey("Given a token provider with an expired cached token", t, func() {
		spy := &tokenClientSpy{output: &token_model.RefreshTokenOutput{AccessToken: "new-token", ExpiresIn: 3600}}
		p := NewTokenProvider(&token_model.RefreshTokenInput{}, spy).(*tokenProvider)
		p.cachedToken = token_model.NewToken("old-token", -2)

		Convey("When GetValidToken is called", func() {
			token, err := p.GetValidToken(context.Background())

			Convey("Then it should refresh and return the new token", func() {
				So(err, ShouldBeNil)
				So(token, ShouldEqual, "new-token")
				So(spy.called, ShouldEqual, 1)
			})
		})
	})
}
