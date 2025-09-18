package docgen

// 1. Centralized list of all valid component names.
#AllComponentNames:
	"DocumentCategoryTitle" |
	"DocumentTitle" |
	"DocumentSubject" |
	"TestBlock" |
	"AuthorBlock"

// 2. Main document plan with compositional rules.
#DocumentPlan: {
	// Optional document properties
	doc_props?: {
		filename?: string
		...
	}
	body: [...#ComponentInstance]
}

// 3. Generic shape of a component instance with scalable 'if' pattern for specific prop validation.
#ComponentInstance: {
	component: #AllComponentNames
	props:     {...}

	// Specific prop validation using if statements
	if component == "DocumentCategoryTitle" {
		props: {
			category_title: string & !=""
		}
	}
	if component == "DocumentTitle" {
		props: {
			document_title: string & !=""
		}
	}
	if component == "DocumentSubject" {
		props: {
			document_subject: string & =~"^DOC-\\d{4,}, Rev [A-Z]$"
		}
	}
	if component == "TestBlock" {
		props: {
			tester_name:     string & !=""
			test_date:       string & =~"^\\d{1,2}/\\d{1,2}/\\d{4}$"
			serial_number:   string & !=""
			test_result:     "PASS" | "FAIL" | "INCOMPLETE"
			additional_info: string
		}
	}
	if component == "AuthorBlock" {
		props: {
			author_name:    string & !=""
			company_name:   string & !=""
			address_line1:  string & !=""
			address_line2?: string // Optional
			city_state_zip: string & !=""
			phone:          string & !=""
			fax?:           string // Optional
			website:        string & !=""
		}
	}
}