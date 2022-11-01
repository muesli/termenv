package termenv

// Notification triggers a notification using OSC777.
func Notification(title, body string) {
	output.Notification(title, body)
}

// Notification triggers a notification using OSC777.
func (o *Output) Notification(title, body string) {
	_, _ = o.WriteString(OSC + "777;notify;" + title + ";" + body + ST)
}
