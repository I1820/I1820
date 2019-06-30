module github.com/I1820/tm/client

go 1.12

require (
	github.com/I1820/pm v0.0.0-20190619001057-1665488291ac
	github.com/I1820/tm v0.0.0-20190619220247-f1b26a3fc484
	github.com/go-resty/resty/v2 v2.0.0
	github.com/gobuffalo/tags v2.0.15+incompatible // indirect
	github.com/microcosm-cc/bluemonday v1.0.2 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7 // indirect
)

replace github.com/go-resty/resty/v2 => github.com/1995parham/resty/v2 v2.0.0
