package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/wish"
	"github.com/creack/pty"
	"github.com/charmbracelet/ssh"
	"github.com/muesli/termenv"
)

type sshOutput struct {
	ssh.Session
	tty *os.File
}

func (s *sshOutput) Write(p []byte) (int, error) {
	return s.Session.Write(p)
}

func (s *sshOutput) Read(p []byte) (int, error) {
	return s.Session.Read(p)
}

func (s *sshOutput) Name() string {
	return s.tty.Name()
}

func (s *sshOutput) Fd() uintptr {
	return s.tty.Fd()
}

type sshEnviron struct {
	environ []string
}

func (s *sshEnviron) Getenv(key string) string {
	for _, v := range s.environ {
		if strings.HasPrefix(v, key+"=") {
			return v[len(key)+1:]
		}
	}
	return ""
}

func (s *sshEnviron) Environ() []string {
	return s.environ
}

func outputFromSession(s ssh.Session) *termenv.Output {
	sshPty, _, _ := s.Pty()
	_, tty, err := pty.Open()
	if err != nil {
		panic(err)
	}
	o := &sshOutput{
		Session: s,
		tty:     tty,
	}
	environ := s.Environ()
	environ = append(environ, fmt.Sprintf("TERM=%s", sshPty.Term))
	e := &sshEnviron{
		environ: environ,
	}
	return termenv.NewOutput(o, termenv.WithUnsafe(), termenv.WithEnvironment(e))
}

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(":2345"),
		wish.WithHostKeyPath("termenv"),
		wish.WithMiddleware(
			func(sh ssh.Handler) ssh.Handler {
				return func(s ssh.Session) {
					output := outputFromSession(s)

					p := output.ColorProfile()
					fmt.Fprintf(s, "\tColor Profile: %d\n", p)

					fmt.Fprintf(s, "\n\t%s %s %s %s %s",
						output.String("bold").Bold(),
						output.String("faint").Faint(),
						output.String("italic").Italic(),
						output.String("underline").Underline(),
						output.String("crossout").CrossOut(),
					)

					fmt.Fprintf(s, "\n\t%s %s %s %s %s %s %s",
						output.String("red").Foreground(p.Color("#E88388")),
						output.String("green").Foreground(p.Color("#A8CC8C")),
						output.String("yellow").Foreground(p.Color("#DBAB79")),
						output.String("blue").Foreground(p.Color("#71BEF2")),
						output.String("magenta").Foreground(p.Color("#D290E4")),
						output.String("cyan").Foreground(p.Color("#66C2CD")),
						output.String("gray").Foreground(p.Color("#B9BFCA")),
					)

					fmt.Fprintf(s, "\n\t%s %s %s %s %s %s %s\n\n",
						output.String("red").Foreground(p.Color("0")).Background(p.Color("#E88388")),
						output.String("green").Foreground(p.Color("0")).Background(p.Color("#A8CC8C")),
						output.String("yellow").Foreground(p.Color("0")).Background(p.Color("#DBAB79")),
						output.String("blue").Foreground(p.Color("0")).Background(p.Color("#71BEF2")),
						output.String("magenta").Foreground(p.Color("0")).Background(p.Color("#D290E4")),
						output.String("cyan").Foreground(p.Color("0")).Background(p.Color("#66C2CD")),
						output.String("gray").Foreground(p.Color("0")).Background(p.Color("#B9BFCA")),
					)

					fmt.Fprintf(s, "\n\t%s %s\n", output.String("Has foreground color").Bold(), output.ForegroundColor())
					fmt.Fprintf(s, "\t%s %s\n", output.String("Has background color").Bold(), output.BackgroundColor())
					fmt.Fprintf(s, "\t%s %t\n", output.String("Has dark background?").Bold(), output.HasDarkBackground())
					fmt.Fprintln(s)

					hw := "Hello, world!"
					output.Copy(hw)
					fmt.Fprintf(s, "\t%q copied to clipboard\n", hw)
					fmt.Fprintln(s)

					fmt.Fprintf(s, "\t%s", output.Hyperlink("http://example.com", "This is a link"))
					fmt.Fprintln(s)

					sh(s)
				}
			},
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on %s", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
