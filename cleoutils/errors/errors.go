// Package errors is a wrapper for github.com/cockroachdb/errors
// to make it easier to replace later, if necessary.
package errors

import "github.com/cockroachdb/errors"

func New(msg string) error {
	return errors.New(msg)
}

func Newf(format string, args ...interface{}) error {
	return errors.Newf(format, args...)
}

func Wrap(err error, msg string) error {
	return errors.Wrap(err, msg)
}

func Wrapf(err error, msg string, args ...interface{}) error {
	return errors.Wrapf(err, msg, args...)
}

func WithHint(err error, msg string) error {
	return errors.WithHint(err, msg)
}

func WithHintf(err error, format string, args ...interface{}) error {
	return errors.WithHintf(err, format, args...)
}

func FlattenHints(err error) string {
	return errors.FlattenHints(err)
}
