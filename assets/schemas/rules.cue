package docgen

// 1. Centralized list of all valid component names.
#AllComponentNames:
	"DocumentCategoryTitle" |
	"DocumentTitle" |
	"DocumentSubject" |
	"TestBlock" |
	"AuthorBlock" |
	"HeadingOneSection" |
	"HeadingTwoSection" |
	"RevisionBlock" |
	"ReferenceTable" |
	"AcronymList" |
	"ImageFigure" |
	"TableOfContents" |
	"TableList" |
	"FigureList" |
	"StandardHeader" |
	"StandardFooter" |
	"DisclaimerStatement" |
	"DisclaimerPage"

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
	if component == "HeadingOneSection" {
		props: {
			heading_text: string & !=""
			content_text: string & !=""
		}
	}
	if component == "HeadingTwoSection" {
		props: {
			heading_text: string & !=""
			content_text: string & !=""
		}
	}
	if component == "RevisionBlock" {
		props: {
			revision_letter:       string & !=""
			revision_date:         string & =~"^\\d{4}/\\d{2}/\\d{2}$"
			changes_description:   string & !=""
			authority_name:        string & !=""
		}
	}
	if component == "ReferenceTable" {
		props: {
			table_caption:           string & !=""
			row1_number:             string & !=""
			row1_title:              string & !=""
			row1_document_number:    string & !=""
			row1_origin:             string & !=""
			row2_number:             string & !=""
			row2_title:              string & !=""
			row2_document_number:    string & !=""
			row2_origin:             string & !=""
			row3_number:             string & !=""
			row3_title:              string & !=""
			row3_document_number:    string & !=""
			row3_origin:             string & !=""
			row4_number:             string & !=""
			row4_title:              string & !=""
			row4_document_number:    string & !=""
			row4_origin:             string & !=""
		}
	}
	if component == "AcronymList" {
		props: {
			list_title:  string & !=""
			acronym1:    string & !=""
			definition1: string & !=""
			acronym2:    string & !=""
			definition2: string & !=""
			acronym3:    string & !=""
			definition3: string & !=""
			acronym4:    string & !=""
			definition4: string & !=""
			acronym5:    string & !=""
			definition5: string & !=""
		}
	}
	if component == "ImageFigure" {
		props: {
			image_ref:       string & !=""
			figure_caption:  string & !=""
		}
	}
	if component == "TableOfContents" {
		props: {
			section1_title: string & !=""
			section1_page:  string & !=""
			section2_number: string & !=""
			section2_title: string & !=""
			section2_page:  string & !=""
			section3_title: string & !=""
			section3_page:  string & !=""
		}
	}
	if component == "TableList" {
		props: {
			list_title:     string & !=""
			table1_number:  string & !=""
			table1_title:   string & !=""
			table1_page:    string & !=""
			table2_number:  string & !=""
			table2_title:   string & !=""
			table2_page:    string & !=""
			table3_number:  string & !=""
			table3_title:   string & !=""
			table3_page:    string & !=""
		}
	}
	if component == "FigureList" {
		props: {
			list_title:      string & !=""
			figure1_number:  string & !=""
			figure1_title:   string & !=""
			figure1_page:    string & !=""
			figure2_number:  string & !=""
			figure2_title:   string & !=""
			figure2_page:    string & !=""
			figure3_number:  string & !=""
			figure3_title:   string & !=""
			figure3_page:    string & !=""
		}
	}
	if component == "StandardHeader" {
		props: {
			classification: string & !=""
			company_name:   string & !=""
		}
	}
	if component == "StandardFooter" {
		props: {
			restrictions_text: string & !=""
			classification:    string & !=""
		}
	}
	if component == "DisclaimerStatement" {
		props: {
			technical_rights_statement: string & !=""
			export_control_statement:   string & !=""
			handling_statement:         string & !=""
			disclaimer_statement:       string & !=""
		}
	}
	if component == "DisclaimerPage" {
		props: {
			page_title:             string & !=""
			general_instructions:   string & !=""
			warning_text:           string & !=""
			handling_instructions:  string & !=""
			second_warning_text:    string & !=""
			detailed_instructions:  string & !=""
			procedure_notes:        string & !=""
		}
	}
}