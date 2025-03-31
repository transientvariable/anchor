package http

// Enumeration of supported HTTP headers.
const (
	HeaderAccept             = "Accept"
	HeaderAcceptCharset      = "Accept-Charset"
	HeaderAcceptDatetime     = "Accept-Datetime"
	HeaderAcceptEncoding     = "Accept-Encoding"
	HeaderAcceptLanguage     = "Accept-Language"
	HeaderAuthorization      = "Authorization"
	HeaderCacheControl       = "Cache-Control"
	HeaderConnection         = "Connection"
	HeaderContentDisposition = "Content-Disposition"
	HeaderContentEncoding    = "Content-Encoding"
	HeaderContentLanguage    = "Content-Language"
	HeaderContentLength      = "Content-Length"
	HeaderContentMD5         = "Content-MD5"
	HeaderContentType        = "Content-Type"
	HeaderCookie             = "Cookie"
	HeaderDate               = "Date"
	HeaderDigest             = "Digest"
	HeaderETag               = "ETag"
	HeaderExpect             = "Expect"
	HeaderExpires            = "Expires"
	HeaderForwarded          = "Forwarded"
	HeaderFrom               = "From"
	HeaderHost               = "Host"
	HeaderIdempotencyKey     = "Idempotency-Key"
	HeaderIfMatch            = "If-Match"
	HeaderIfModifiedSince    = "If-Modified-Since"
	HeaderIfNoneMatch        = "If-None-Match"
	HeaderIfUnmodifiedSince  = "If-Unmodified-Since"
	HeaderXIpfsCid           = "X-Ipfs-Cid"
	HeaderXIpfsPath          = "X-Ipfs-path"
	HeaderXIpfsRoots         = "X-Ipfs-Roots"
	HeaderLastModified       = "Last-Modified"
	HeaderLocation           = "Location"
	HeaderOrigin             = "Origin"
	HeaderProxyAuthorization = "Proxy-Authorization"
	HeaderRange              = "Range"
	HeaderReferer            = "Referer"
	HeaderTransferEncoding   = "Transfer-Encoding"
	HeaderUserAgent          = "User-Agent"
	HeaderWantDigest         = "Want-Digest"
	HeaderWantContentDigest  = "Want-Content-Digest"
)

// Headers returns a list of all supported HTTP headers.
func Headers() []string {
	return []string{
		HeaderAccept,
		HeaderAcceptCharset,
		HeaderAcceptDatetime,
		HeaderAcceptEncoding,
		HeaderAcceptLanguage,
		HeaderAuthorization,
		HeaderCacheControl,
		HeaderConnection,
		HeaderContentDisposition,
		HeaderContentEncoding,
		HeaderContentLanguage,
		HeaderContentLength,
		HeaderContentMD5,
		HeaderContentType,
		HeaderCookie,
		HeaderDate,
		HeaderDigest,
		HeaderETag,
		HeaderExpect,
		HeaderExpires,
		HeaderForwarded,
		HeaderFrom,
		HeaderHost,
		HeaderIdempotencyKey,
		HeaderIfMatch,
		HeaderIfModifiedSince,
		HeaderIfNoneMatch,
		HeaderIfUnmodifiedSince,
		HeaderXIpfsCid,
		HeaderXIpfsPath,
		HeaderXIpfsRoots,
		HeaderLastModified,
		HeaderLocation,
		HeaderOrigin,
		HeaderProxyAuthorization,
		HeaderRange,
		HeaderReferer,
		HeaderTransferEncoding,
		HeaderUserAgent,
		HeaderWantDigest,
		HeaderWantContentDigest,
	}
}
