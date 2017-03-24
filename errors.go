package httpreserve

const errorBlankProtocol = "blank protocol"
const errorUnknownProtocol = "unknown protocol"

const errorNoIALink = "no internet archive record"

// ErrorIAExists so that we can identify links we
// do not need to process a second time, or send to IA
const ErrorIAExists = "already an internet archive record"
