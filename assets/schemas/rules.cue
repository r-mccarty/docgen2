package docgen

// Document plan schema
#DocumentPlan: {
	// Optional document properties
	doc_props?: {
		filename?: string
		...
	}

	// Body must be an array of component instances
	body: [...#ComponentInstance]
}

// Component-specific validation - no generic fallback
#ComponentInstance: #DocumentCategoryTitleComponent | #DocumentTitleComponent | #DocumentSubjectComponent | #TestBlockComponent | #AuthorBlockComponent

#DocumentCategoryTitleComponent: {
	component: "DocumentCategoryTitle"
	props: {
		category_title: string & !=""
	}
}

#DocumentTitleComponent: {
	component: "DocumentTitle"
	props: {
		document_title: string & !=""
	}
}

#DocumentSubjectComponent: {
	component: "DocumentSubject"
	props: {
		document_subject: string & =~"^DOC-\\d{4,}, Rev [A-Z]$"
	}
}

#TestBlockComponent: {
	component: "TestBlock"
	props: {
		tester_name:     string & !=""
		test_date:       string & =~"^\\d{1,2}/\\d{1,2}/\\d{4}$"
		serial_number:   string & !=""
		test_result:     "PASS" | "FAIL" | "INCOMPLETE"
		additional_info: string
	}
}

#AuthorBlockComponent: {
	component: "AuthorBlock"
	props: {
		author_name:     string & !=""
		company_name:    string & !=""
		address_line1:   string & !=""
		address_line2:   string & !=""
		city_state_zip:  string & !=""
		phone:           string & !=""
		fax:             string & !=""
		website:         string & !=""
	}
}