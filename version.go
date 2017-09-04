package httpreserve

// User-agent to identify code being run
// e.g. "exponentialDK-httpreserve/0.0.0
const httpUSERAGENT = "exponentialDK-httpreserve/"

// Version will return a simple version number for the app.
func Version() string {
	return "0.0.8"
}

// VersionNumber is a synonym for Version()
func VersionNumber() string {
	return Version()
}

// VersionText will return the full text version information
// e.g. for the useragent to query our websites.
func VersionText() string {
	return httpUSERAGENT + Version()
}

// Dedication will return a dedication string for a colleague's
// dear brother who was lost on the day I first figured out this
// code.
func Dedication() string {
	return "Dedicated to Matthew Croad."
}
