package utils

import (
	"strings"
)

// SingleLeading ensures that a string s starts with ch, but doesnt end with it
func SingleLeading(s string, ch string) string {
	/*
	Examples
	/ 		-> /
	// 		-> /
	/a/ 	-> /a
	/a/b/d 	-> /a/b/d
	a/b 	-> /a/b
	/a/b/ 	-> /a/b
	a/b/ 	-> /a/b
	*/
	s = strings.Trim(s, ch)
	return ch + s
}

// SingleTrailing ensures that a string s ends with ch, but doesnt start with it
func SingleTrailing(s string, ch string) string {
	s = strings.Trim(s, ch)
	return s + ch
}

// SingleLeadingTrailing ensures that a string s starts and ends with single occurence of ch
func SingleLeadingTrailing(s string, ch string) string {
	/*
	Examples
	/ 		-> /
	// 		-> /
	/a/ 	-> /a/
	/a/b/d 	-> /a/b/d/
	a/b 	-> /a/b/
	/a/b/ 	-> /a/b/
	a/b/ 	-> /a/b/
	*/
	s = strings.Trim(s, ch)
	s = strings.TrimRight(ch + s, ch)
	return s + ch
}

// JoinLeading joins s1 and s2, using ch, and ensures ch occurs at the start
func JoinLeading(s1 string, s2 string, ch string) string {
	s1 = strings.Trim(s1, ch)
	s2 = strings.Trim(s2, ch)
	s := s1 + ch + s2
	return SingleLeading(s, ch)
}

// JoinTrailing joins s1 and s2, using ch, and ensures ch occurs at the end
func JoinTrailing(s1 string, s2 string, ch string) string {
	s1 = strings.Trim(s1, ch)
	s2 = strings.Trim(s2, ch)
	s := s1 + ch + s2
	return SingleTrailing(s, ch)
}

// JoinLeadingTrailing joins s1 and s2, using ch, and ensures ch occurs at the start and end
func JoinLeadingTrailing(s1 string, s2 string, ch string) string {
	s1 = strings.Trim(s1, ch)
	s2 = strings.Trim(s2, ch)
	s := s1 + ch + s2
	return SingleLeadingTrailing(s, ch)
}
