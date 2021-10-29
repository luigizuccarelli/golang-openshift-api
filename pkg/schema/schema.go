package schema

type RulesSchema struct {
	Rules []Rule
}

type Rule struct {
	Name      string `json:"name"`
	Api       string `json:"api"`
	Filter    string `json:"filter"`
	Regex     string `json:"regex"`
	Occurence string `json:"occurence"`
	Operand   string `json:"operand"`
}

type Request struct {
	Id      string `json:"id,omitemptye"`
	Message string `json:"message"`
}

// Response schema
type Response struct {
	Code    int    `json:"code,omitempty"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
