// SPDX-License-Identifier: AGPL-3.0-only

package api

import "net/http"

func setPrivateResponseHeaders(w http.ResponseWriter) {
	setPrivateResponseHeaders(w)
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
	w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
}

// writeAccepted emits a bodyless successful mutation response with the same
// privacy and content-sniffing protections used by JSON API responses.
func writeAccepted(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusAccepted)
}
