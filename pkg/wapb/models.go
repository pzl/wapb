package wapb

type StoredCommon struct {
	ID string
}

type File struct {
	StoredCommon
}

type Redirect struct {
	StoredCommon
}

type Text struct {
	StoredCommon
}
