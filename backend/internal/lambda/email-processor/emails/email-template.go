package emailtemplate

import (
	"bytes"
	"fmt"
	"html/template"
	sqs "inside-athletics/internal/sqs"
	"strings"
)

var disasterEmailHTMLTemplate = template.Must(template.New("disasterEmail").Parse(`
<!DOCTYPE html>
<html>
<body style="display:flex; flex-direction:column; background-color:#f3f4f6; margin:0; padding:0; font-family:'PT Sans', Arial, sans-serif; justify-content:center; align-items:center;">
  <div style="display:flex; flex-direction:column; background-color:#ffffff; max-width:36rem; margin:0 auto;">

    <div style="background-color:#ffffff; position:relative; padding:3rem 2.5rem 0 2.5rem; min-height:110px; max-width:36rem;">
      <div style="position:absolute; top:0; left:0; width:100%; height:100%; overflow:hidden;">
        <svg width="100%" height="100%" viewBox="0 0 100 100" preserveAspectRatio="none" style="display:block;">
          <polygon points="0,0 75,0 65,100 0,100" fill="#8a1e41" />
        </svg>
      </div>
      <h1 style="color:#ffffff; font-weight:600; font-size:clamp(1.25rem,3vw,3rem); margin:0; position:relative; z-index:10;">
        Inside Athletics
      </h1>
    </div>

    <div style="padding:2.5rem;">
      <p style="font-size:1.125rem; margin-bottom:2rem;">
        <strong>New Reply in Inside Athletics:</strong>
      </p>

      <p style="font-size:1rem; margin-bottom:2rem;">{{.Message}}</p>

      <hr style="border-color:#d1d5db; margin:2rem 0;" />

      <img style="width:12.5rem; height:5rem;" src="https://prisere.com/wp-content/uploads/2023/09/Prisere-logo-transparent.png" />
    </div>

  </div>
</body>
</html>
`))

type disasterEmailTemplateData struct {
	Message string
}

func RenderDisasterEmailHTML(message sqs.DisasterEmailMessage) string {
	data := disasterEmailTemplateData{
		Message: message.Message,
	}
	var buf bytes.Buffer
	if err := disasterEmailHTMLTemplate.Execute(&buf, data); err != nil {
		return fmt.Sprintf("error rendering email: %v", err)
	}
	return buf.String()
}

func RenderDisasterEmailText(message sqs.DisasterEmailMessage) string {
	return strings.TrimSpace(fmt.Sprintf(`New Reply in Inside Athletics:

%s`, message.Message))
}
