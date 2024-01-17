package templates

import "embed"

//go:embed homepage.html.tmpl
var HomepageTemplateFS embed.FS

//go:embed base.html.tmpl list_post.html.tmpl
var ListPostTemplateFS embed.FS

//go:embed detail_post.html.tmpl base.html.tmpl
var DetailPostTemplateFS embed.FS

//go:embed new_post.html.tmpl base.html.tmpl
var NewPostTemplateFS embed.FS

//go:embed post_edit.html.tmpl base.html.tmpl
var EditPostTemplateFS embed.FS

//go:embed login_page.html.tmpl base.html.tmpl
var LoginPageFS embed.FS

//go:embed four_oh_four.html.tmpl base.html.tmpl
var MissingFS embed.FS

//go:embed four_oh_five.html.tmpl base.html.tmpl
var MethodNotAllowedFS embed.FS
