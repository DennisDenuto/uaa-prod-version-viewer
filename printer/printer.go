package printer

import (
	"github.com/DennisDenuto/uaa-prod-version-viewer/pcfnotes"
	"io"
	"fmt"
	"github.com/johnmccabe/bitbar"
)

type BoshPackage struct {
	Name    string
	Version string
}

type BoshRelease struct {
	Name     string
	Version  string
	Packages []BoshPackage
}

func (br BoshRelease) GithubURL() string {
	return fmt.Sprintf("https://github.com/cloudfoundry/%s/tree/%s", br.Name, br.Version)
}

func (bp BoshPackage) GithubURL() string {
	return fmt.Sprintf("https://github.com/cloudfoundry/%s/tree/%s", bp.Name, bp.Version)
}

type LineItem struct {
	BoshRelease BoshRelease
	PcfVersion  pcfnotes.Version
}

type Printer interface {
	Print([]LineItem) error
}

type BitBarPrinter struct {
	Writer           io.Writer
	StatusIconBase64 string
	ColorText        string
}

func (bb BitBarPrinter) Print(lineItems []LineItem) error {
	bitbarStyle.Color = bb.ColorText

	b := &bitbar.Plugin{}
	b.StatusLine("").Image(bb.StatusIconBase64)
	menu := b.NewSubMenu()
	for _, item := range lineItems {
		bb.printPCFLine(b, menu, item)
	}

	bb.printCSVExportOption(b, menu)

	_, err := bb.Writer.Write([]byte(b.Render()))
	if err != nil {
		panic(err)
	}

	return nil
}
func (b BitBarPrinter) printStatus() error {
	_, err := b.Writer.Write([]byte(fmt.Sprintf(`| image=%s`, b.StatusIconBase64)))
	return err
}

var bitbarStyle = bitbar.Style{
	Font:  "UbuntuMono-Bold",
	Color: "black",
	Size:  12,
}

func (b BitBarPrinter) printPCFLine(bitbarPlugin *bitbar.Plugin, menu *bitbar.SubMenu, item LineItem) {
	menu.HR()
	menu.Line(fmt.Sprintf("PCF %s", item.PcfVersion.String())).Style(bitbarStyle)
	menu.HR()
	menu.Line(fmt.Sprintf("%s: %s", item.BoshRelease.Name, item.BoshRelease.Version)).Style(bitbarStyle).Href(item.BoshRelease.GithubURL())

	for _, boshPackage := range item.BoshRelease.Packages {
		submenu := bitbarPlugin.SubMenu.NewSubMenu()
		submenu.Line(fmt.Sprintf("%s: version %s", boshPackage.Name, boshPackage.Version)).Href(boshPackage.GithubURL())
	}
}

func (printer BitBarPrinter) printCSVExportOption(bitbarPlugin *bitbar.Plugin, menu *bitbar.SubMenu) {
	menu.HR()
	menu.Line("Reports").Style(bitbarStyle)

	cmd := bitbar.Cmd{
		Bash:   "uaa-prod-version-viewer",
		Params: []string{"csv"},
	}

	bitbarPlugin.SubMenu.NewSubMenu().Line("CSV").Command(cmd)
}
