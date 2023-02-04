package cmdr

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vanilla-os/orchid/roff"
)

// Command represents a cli command which
// may have flags, arguments, and children commands.
type Command struct {
	*cobra.Command
	children []*Command
}

// AddCommand adds a command to the slice and to the underlying
// cobra command.
func (c *Command) AddCommand(commands ...*Command) {
	c.children = append(c.children, commands...)
	for _, cmd := range commands {
		c.Command.AddCommand(cmd.Command)
	}
}

// Children returns the children commands.
func (c *Command) Children() []*Command {
	return c.children
}

// WithBoolFlag adds a boolean flag to the command and
// registers the flag with environment variable injection
func (c *Command) WithBoolFlag(f BoolFlag) *Command {
	c.Command.Flags().BoolP(f.Name, f.Shorthand, f.Value, f.Usage)
	viper.BindPFlag(f.Name, c.Command.Flags().Lookup(f.Name))
	return c
}

// WithPersistentBoolFlag adds a persistent boolean flag to the command and
// registers the flag with environment variable injection
func (c *Command) WithPersistentBoolFlag(f BoolFlag) *Command {
	c.Command.PersistentFlags().BoolP(f.Name, f.Shorthand, f.Value, f.Usage)
	viper.BindPFlag(f.Name, c.Command.PersistentFlags().Lookup(f.Name))
	return c
}

// WithStringFlag adds a string flag to the command and registers
// the command with the environment variable injection
func (c *Command) WithStringFlag(f StringFlag) *Command {
	c.Command.Flags().StringP(f.Name, f.Shorthand, f.Value, f.Usage)
	viper.BindPFlag(f.Name, c.Command.Flags().Lookup(f.Name))
	return c
}

// WithPersistentStringFlag adds a persistent string flag to the command and registers
// the command with the environment variable injection
func (c *Command) WithPersistentStringFlag(f BoolFlag) *Command {
	c.Command.PersistentFlags().BoolP(f.Name, f.Shorthand, f.Value, f.Usage)
	viper.BindPFlag(f.Name, c.Command.PersistentFlags().Lookup(f.Name))
	return c
}

// NewCommand returns a new Command with the provided inputs. Alias for
// NewCommandRunE.
func NewCommand(use, long, short string, runE func(cmd *cobra.Command, args []string) error) *Command {
	return NewCommandRunE(use, long, short, runE)
}

// NewCommandRunE returns a new Command with the provided inputs. The runE function
// is used for commands that return an error.
func NewCommandRunE(use, long, short string, runE func(cmd *cobra.Command, args []string) error) *Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		RunE:  runE,
	}
	return &Command{
		Command:  cmd,
		children: make([]*Command, 0),
	}
}

// NewCommandRun returns a new Command with the provided inputs. The run function
// is used for commands that do not return an error.
func NewCommandRun(use, long, short string, run func(cmd *cobra.Command, args []string)) *Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Run:   run,
	}
	return &Command{
		Command:  cmd,
		children: make([]*Command, 0),
	}
}

// NewCustomCommand returns a Command created from
// the provided cobra.Command
func NewCommandCustom(cmd *cobra.Command) *Command {
	return &Command{
		Command:  cmd,
		children: make([]*Command, 0),
	}
}

func (c *Command) doc(d *roff.Document) {
	c.docName(d)
	c.docSynopsis(d)
	c.docDescription(d)
	c.docOptions(d)
	c.docCommands(d)
	c.docExamples(d)
}

func (c *Command) docName(d *roff.Document) {
	d.Section("subcommand " + c.Name())
	d.Indent(4)
	d.Text(c.Short)
	d.IndentEnd()
	d.EndSection()
}

func (c *Command) docSynopsis(d *roff.Document) {
	d.SubSection("Synopsis")
	d.Indent(4)
	d.TextBold(c.Name())
	d.Text(" [command] [flags] [arguments]")
	d.IndentEnd()
	d.EndSection()
}

func (c *Command) docDescription(d *roff.Document) {
	d.SubSection("Description")
	d.Indent(4)
	d.TaggedParagraph(4)
	d.Text(c.Long)
	d.IndentEnd()
	d.EndSection()

}

func (c *Command) docOptions(d *roff.Document) {
	d.SubSection("Options")
	d.Text(c.Flags().FlagUsages())
	d.SubSection("Global Options")
	d.Text(c.Parent().PersistentFlags().FlagUsages())
	d.EndSection()
}
func (c *Command) docExamples(d *roff.Document) {
	if c.Example == "" {
		return
	}
	d.SubSection("Examples")
	d.Indent(4)
	d.Text(c.Example)
	d.IndentEnd()
	d.EndSection()

}

func (c *Command) docCommands(d *roff.Document) {
	if len(c.children) == 0 {
		return
	}
	for _, child := range c.Children() {
		if child.Hidden {
			continue
		}

		d.Section(child.Name())

		d.Indent(4)

		d.Text(child.Short + "\n")
		d.IndentEnd()
	}
}
