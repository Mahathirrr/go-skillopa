package data

type Interest struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

var Interests = []Interest{
	{
		Type:  "subCategory",
		Title: "Web Development",
		Slug:  "web-development",
	},
	{
		Type:  "subCategory",
		Title: "User Experience Design",
		Slug:  "user-experience-design",
	},
	{
		Type:  "subCategory",
		Title: "Growth Hacking",
		Slug:  "growth-hacking",
	},
	{
		Type:  "subCategory",
		Title: "Management",
		Slug:  "management",
	},
	{
		Type:  "subCategory",
		Title: "Communication",
		Slug:  "communication",
	},
	{
		Type:  "subCategory",
		Title: "Digital Marketing",
		Slug:  "digital-marketing",
	},
}