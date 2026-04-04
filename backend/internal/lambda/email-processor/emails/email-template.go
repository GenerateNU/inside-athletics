package emailtemplate

import (
	"bytes"
	"fmt"
	"html/template"
	sqs "inside-athletics/internal/sqs"
	"strings"
)

var replyEmailHTMLTemplate = template.Must(template.New("disasterEmail").Parse(`
<!DOCTYPE html>
<html>
<head>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&display=swap" rel="stylesheet">
</head>
<body style="display:flex; flex-direction:column; background-color:#f3f4f6; margin:0; padding:0; font-family:'Inter', Arial, sans-serif; justify-content:center; align-items:center;">
  <div style="display:flex; flex-direction:column; background-color:#ffffff; max-width:36rem; margin:0 auto;">
    <div style="background:#00804D; padding:1.5rem 2.5rem;">
      <h1 style="color:#ffffff; font-weight:600; font-size:1.75rem; margin:0; font-family:'Inter', Arial, sans-serif;">Inside Athletics</h1>
    </div>
    <div style="padding:2.5rem;">
      <p style="font-size:1.125rem; margin-bottom:2rem;">
        <strong>New Reply in Inside Athletics:</strong>
      </p>
      <p style="font-size:1rem; margin-bottom:2rem; color:#444;">Lowkey this school is ahhh don't join</p>
      <hr style="border:none; border-top:1px solid #d1d5db; margin:2rem 0;" />
      <img style="width:12.5rem; height:5rem; object-fit:contain;" src="https://prisere.com/wp-content/uploads/2023/09/Prisere-logo-transparent.png" />
    </div>
  </div>
</body>
</html>
`))

type replyEmailTemplateData struct {
	Message string
}

func RenderReplyEmailHTML(message sqs.ReplyEmailMessage) string {
	data := replyEmailTemplateData{
		Message: message.Message,
	}
	var buf bytes.Buffer
	if err := replyEmailHTMLTemplate.Execute(&buf, data); err != nil {
		return fmt.Sprintf("error rendering email: %v", err)
	}
	return buf.String()
}

func RenderReplyEmailText(message sqs.ReplyEmailMessage) string {
	return strings.TrimSpace(fmt.Sprintf(`New Reply in Inside Athletics:

%s`, message.Message))
}
