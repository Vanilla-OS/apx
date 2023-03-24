package cmdr

import (
	"embed"
	"errors"
	"log"
	"os"
	"path"
	"time"

	"github.com/fitv/go-i18n"
	"github.com/spf13/viper"
	"github.com/vanilla-os/orchid"
	"github.com/vanilla-os/orchid/roff"
)

// The App struct represents the cli application
// with supporting functionality like internationalization
// and logging.
type App struct {
	Name        string
	Version     string
	RootCommand *Command
	Logger      *log.Logger
	logFile     *os.File
	locales     embed.FS
	*i18n.I18n
}

// NewApp creates a new command line application.
// It requires an embed.FS with a top level directory
// named 'locales'.
func NewApp(name string, version string, locales embed.FS) *App {
	// for application logs
	orchid.InitLog(name+" : 	", log.LstdFlags)

	viper.SetEnvPrefix(name)
	viper.AutomaticEnv()

	i18n, err := i18n.New(locales, "locales")
	if err != nil {
		Error.Println(err)
		os.Exit(1)
	}
	i18n.SetDefaultLocale(orchid.Locale())
	a := &App{
		Name:    name,
		Logger:  log.Default(),
		I18n:    i18n,
		Version: version,
		locales: locales,
	}
	err = a.logSetup()
	if err != nil {
		log.Printf("error setting up logging: %v", err)
	}
	return a

}
func (a *App) logSetup() error {
	err := a.ensureLogDir()
	if err != nil {
		return err
	}
	logDir, err := getLogDir(a.Name)
	if err != nil {
		return err
	}
	logFile := path.Join(logDir, a.Name+".log")
	//create your file with desired read/write permissions
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	a.logFile = f

	//set output of logs to f
	log.SetOutput(a.logFile)
	return nil

}
func (a *App) CreateRootCommand(c *Command) {
	a.RootCommand = c
	c.DisableAutoGenTag = true
	manCmd := NewManCommand(a)
	a.RootCommand.AddCommand(manCmd)

	docsCmd := NewDocsCommand(a)
	a.RootCommand.AddCommand(docsCmd)
}

func (a *App) Run() error {
	if a.logFile != nil {
		defer a.logFile.Close()
	}
	if a.RootCommand != nil {
		return a.RootCommand.Execute()
	}
	return errors.New("no root command defined")
}

func (a *App) ensureLogDir() error {
	logPath, err := getLogDir(a.Name)
	if err != nil {
		return err
	}
	return os.MkdirAll(logPath, 0755)
}

func getLogDir(app string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".local", "share", app), nil
}

func (a *App) rootDocs(d *roff.Document) {
	a.docHeading(d)
	a.docName(d)
	a.docSynopsis(d)
	a.docDescription(d)
	a.docOptions(d)
	a.docCommands(d)
	d.Section("Authors")
	d.Paragraph()
	d.Text("Written by Vanilla OS contributors.")
	d.EndSection()
	d.Section("Report bugs to")
	d.Paragraph()
	d.Text("https://github.com/vanilla-os/" + a.Name + "/issues")
	d.EndSection()
	d.Section("License")
	d.Paragraph()
	d.Text("GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>.")
	d.EndSection()

}

func (a *App) docHeading(d *roff.Document) {
	d.Heading(1, a.Name, "User Manual", time.Now())

}

func (a *App) docName(d *roff.Document) {
	d.Section("Name")
	d.Indent(4)
	d.Text(a.RootCommand.Name() + " - " + a.RootCommand.Short)
	d.IndentEnd()
}

func (a *App) docSynopsis(d *roff.Document) {
	d.Section("Synopsis")
	d.Indent(4)
	d.Text(a.RootCommand.Name() + " [command] [flags] [arguments]")
	d.IndentEnd()
}

func (a *App) docDescription(d *roff.Document) {
	d.Section("Description")
	d.Indent(4)
	d.Text(a.RootCommand.Long)
	d.IndentEnd()

}

func (a *App) docOptions(d *roff.Document) {
	d.Section("Options")
	d.Text(a.RootCommand.Flags().FlagUsages())

}

func (a *App) docCommands(d *roff.Document) {
	d.Section(a.RootCommand.Name() + " Commands")
	d.Indent(4)
	for _, c := range a.RootCommand.Children() {
		if c.Hidden {
			continue
		}
		d.TextBold(c.Name())
		d.Indent(4)
		d.Text(c.Short + "\n")
		d.IndentEnd()
	}
	d.IndentEnd()
}
