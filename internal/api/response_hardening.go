// SPDX-License-Identifier: AGPL-3.0-only

package api

import "net/http"

// writeAccepted emits a bodyless successful mutation response with the same
// privacy and content-sniffing protections used by JSON API responses.
func writeAccepted(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusAccepted)
}
