import marked from "marked"
import escapeHtml from "escape-html"

marked.setOptions({
	renderer: new marked.Renderer(),
	pedantic: true,
	sanitizer: (input) => escapeHtml(input),
})

export default marked
