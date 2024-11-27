// Package values is a collection of constants and types that are used throughout the Kairos init system.
// This are mainly interfaces and constants that we can refer from everywhere and that other packages can import without
// messing with circular dependencies and such.
// Not only interfaces can be added, but if a type is generic enough and can be used in different packages, it can be
// added here. See for example the System struct. We need it and its fully complete, we dont need to have an interface for it as
// its not going to import anything else.
package values
