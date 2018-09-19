package webserver

type MultiAttachmentMode string

const (
	MultiAttachmentModeBoundary MultiAttachmentMode = "boundary"
	MultiAttachmentModeZip      MultiAttachmentMode = "zip"
)
