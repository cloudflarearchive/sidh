// +build !noasm

// This file is here because Go cannot handle a package where all source
// files are excluded due to build constraints, as is currently the case
// for the arm64 build. Now there is always a valid package.
package utils
