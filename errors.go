package httpreserve

const errorBlankProtocol = "blank protocol"
const errorUnknownProtocol = "unknown protocol"

//snapshot errors and variables

// generate screenshots...
var snapshot = true

// SnapshotNotEnabled helps end users verify server status
var SnapshotNotEnabled = "snapshots are not currently enabled"

// GenerateSnapshotErr tells us we've something else wrong
var GenerateSnapshotErr = "error generating snapshot"

// ResponseIncorrect tells us we haven't created a screenshot because
// the domain is no longer in existence...
var ResponseIncorrect = "snapshots not created for response codes zero or greater than 400"
