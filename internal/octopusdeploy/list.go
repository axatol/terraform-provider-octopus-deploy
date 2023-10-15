package octopusdeploy

type List[T any] struct {
	ItemType       string `json:"ItemType"`
	Items          []T    `json:"Items"`
	ItemsPerPage   int    `json:"ItemsPerPage"`
	LastPageNumber int    `json:"LastPageNumber"`
	NumberOfPages  int    `json:"NumberOfPages"`
	TotalResults   int    `json:"TotalResults"`
	Links          struct {
		Page_All     string `json:"Page.All"`
		Page_Current string `json:"Page.Current"`
		Page_Last    string `json:"Page.Last"`
		Self         string `json:"Self"`
		Template     string `json:"Template"`
	} `json:"Links"`
}
