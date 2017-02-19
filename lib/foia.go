package lib

// Result ...
type Result struct {
	Subject       string
	DocumentClass string
	PdfLink       string
	OriginalLink  string
	DocDate       int
	PostedDate    int
	From          string
	To            string
	MessageNumber string
	CaseNumber    string
}

// FoiaResults ...
type FoiaResults struct {
	Success    bool
	TotalHits  int
	Results    []Result
	QueryText  string
	FieldMatch string
	Response   string
}
