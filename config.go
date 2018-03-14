package main

const MODE_TEXT = 0
const MODE_ATTRIBUTE = 1

type Field struct {
	Name      string
	Selector  []string
	Attribute []string
	Mode      int8
}

type ParseConfig struct {
	Fields []Field
}

var parseConfig ParseConfig = NewConfig()

func NewConfig() ParseConfig {
	return ParseConfig{
		Fields: []Field{
			Field{
				Name:      "price",
				Selector:  []string{"[itemprop=\"price\"]"},
				Attribute: []string{"content"},
				Mode:      MODE_ATTRIBUTE,
			},
			Field{
				Name:      "alt_price",
				Selector:  []string{".price", "#price", ".product__price"},
				Attribute: []string{},
				Mode:      MODE_TEXT,
			},
			Field{
				Name:      "title",
				Selector:  []string{"title"},
				Attribute: []string{},
				Mode:      MODE_TEXT,
			},
			Field{
				Name:      "priceCurrency",
				Selector:  []string{"meta[itemprop=\"priceCurrency\"]"},
				Attribute: []string{"content"},
				Mode:      MODE_ATTRIBUTE,
			},
			Field{
				Name:      "description",
				Selector:  []string{"meta[name=\"description\"]"},
				Attribute: []string{"content"},
				Mode:      MODE_ATTRIBUTE,
			},
			Field{
				Name:      "sku",
				Selector:  []string{"meta[name=\"sku\"]"},
				Attribute: []string{"content"},
				Mode:      MODE_ATTRIBUTE,
			},
			Field{
				Name:      "og:title",
				Selector:  []string{"meta[property=\"og:title\"]"},
				Attribute: []string{"content"},
				Mode:      MODE_ATTRIBUTE,
			},
			Field{
				Name:      "name",
				Selector:  []string{"meta[itemprop=\"name\"]"},
				Attribute: []string{"content"},
				Mode:      MODE_ATTRIBUTE,
			},
		},
	}
}
