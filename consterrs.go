// (C)Tylor Arndt 2014, Mozilla Public License (MPL) Version 2.0
// See LICENSE file for details.
// For other licensing options please contact the author.

package errs

type ConstErr string

func (this ConstErr) Error() string { return string(this) }

const ErrNotImplemented = ConstErr("This feature is not implemented")

var _ error = ErrNotImplemented
