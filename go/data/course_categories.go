package data

type SubCategory struct {
	ID    int    `json:"id,omitempty"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Category struct {
	ID            int           `json:"id,omitempty"`
	Title         string        `json:"title"`
	URL           string        `json:"url"`
	SubCategories []SubCategory `json:"subCategories"`
}

var Categories = []Category{
	{
		ID:    100,
		Title: "Development",
		URL:   "/courses/development/",
		SubCategories: []SubCategory{
			{
				ID:    101,
				Title: "Web Development",
				URL:   "/courses/development/web-development/",
			},
			{
				ID:    102,
				Title: "Mobile Development",
				URL:   "/courses/development/mobile-apps/",
			},
			{
				ID:    103,
				Title: "Programming Languages",
				URL:   "/courses/development/programming-languages/",
			},
			{
				ID:    104,
				Title: "Game Development",
				URL:   "/courses/development/game-development/",
			},
			{
				ID:    105,
				Title: "Database Design and Development",
				URL:   "/courses/development/database/",
			},
			{
				ID:    106,
				Title: "Software Testing",
				URL:   "/courses/development/software-testing/",
			},
		},
	},
	{
		ID:    200,
		Title: "Business",
		URL:   "/courses/business/",
		SubCategories: []SubCategory{
			{
				ID:    201,
				Title: "Entrepreneurship",
				URL:   "/courses/business/entrepreneurship/",
			},
			{
				ID:    202,
				Title: "Communication",
				URL:   "/courses/business/communications/",
			},
			{
				ID:    203,
				Title: "Management",
				URL:   "/courses/business/management/",
			},
			{
				ID:    204,
				Title: "Sales",
				URL:   "/courses/business/sales/",
			},
			{
				ID:    205,
				Title: "Business Strategy",
				URL:   "/courses/business/strategy/",
			},
		},
	},
	// ... other categories follow the same pattern
}