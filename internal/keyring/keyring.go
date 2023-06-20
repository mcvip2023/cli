// Package keyring is a simple wrapper that adds timeouts to the zalando/go-keyring package.
package keyring

import (
	"context"
	"time"

	"github.com/zalando/go-keyring"
)

type TimeoutError struct {
	message string
}

func (e *TimeoutError) Error() string {
	return e.message
}

// Set secret in keyring for user.
func Set(service, user, secret string) error {
	duration := time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		ch <- keyring.Set(service, user, secret)
	}()
	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return &TimeoutError{"timeout while trying to set secret in keyring"}
	}
}

// Get secret from keyring given service and user name.
func Get(service, user string) (string, error) {
	duration := time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	ch := make(chan struct {
		val string
		err error
	}, 1)
	go func() {
		defer close(ch)
		val, err := keyring.Get(service, user)
		ch <- struct {
			val string
			err error
		}{val, err}
	}()
	select {
	case res := <-ch:
		return res.val, res.err
	case <-ctx.Done():
		return "", &TimeoutError{"timeout while trying to get secret from keyring"}
	}
}

// Delete secret from keyring.
func Delete(service, user string) error {
	duration := time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		ch <- keyring.Delete(service, user)
	}()
	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return &TimeoutError{"timeout while trying to delete secret from keyring"}
	}
}
