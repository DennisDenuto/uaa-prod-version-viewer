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
}

func (bb BitBarPrinter) Print(lineItems []LineItem) error {
	b := &bitbar.Plugin{}
	b.StatusLine("").Image(bb.StatusIconBase64)
	menu := b.NewSubMenu()
	for _, item := range lineItems {
		bb.printPCFLine(b, menu, item)
	}

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

func (b BitBarPrinter) printPCFLine(bitbarPlugin *bitbar.Plugin, menu *bitbar.SubMenu, item LineItem) {
	s := bitbar.Style{
		Font:  "UbuntuMono-Bold",
		Color: "white",
		Size:  12,
	}
	menu.HR()
	menu.Line(fmt.Sprintf("PCF %s", item.PcfVersion.String())).Style(s)
	menu.HR()
	menu.Line(fmt.Sprintf("%s: %s", item.BoshRelease.Name, item.BoshRelease.Version)).Style(s).Href(item.BoshRelease.GithubURL())

	for _, boshPackage := range item.BoshRelease.Packages {
		submenu := bitbarPlugin.SubMenu.NewSubMenu()
		submenu.Line(fmt.Sprintf("%s: version %s", boshPackage.Name, boshPackage.Version)).Href(boshPackage.GithubURL())
	}
}
