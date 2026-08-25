/*
Copyright (C) 2026 Yukthi Systems Private Limited

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License version 3
as published by the Free Software Foundation.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
version 3 along with this program. If not, see
<https://www.gnu.org/licenses/>.
*/

// Package ipfire downloads and parses the IPFire Domain Blocklist (DBL)
// plaintext category lists (https://www.ipfire.org/dbl/how-to-use).
package ipfire

// Category identifies one IPFire DBL list.
type Category struct {
	Name string
	URL  string
}

// Categories is the fixed set of IPFire DBL lists this service consumes.
// Name is also used to derive the category's score environment variable
// (IPFIRE_<NAME>_SCORE) in internal/config.
var Categories = []Category{
	{Name: "gambling", URL: "https://dbl.ipfire.org/lists/gambling/domains.txt"},
	{Name: "malware", URL: "https://dbl.ipfire.org/lists/malware/domains.txt"},
	{Name: "phishing", URL: "https://dbl.ipfire.org/lists/phishing/domains.txt"},
	{Name: "pornography", URL: "https://dbl.ipfire.org/lists/porn/domains.txt"},
	{Name: "violence", URL: "https://dbl.ipfire.org/lists/violence/domains.txt"},
}
