package printer

import (
	"io"
	"encoding/csv"
)

type CSVPrinter struct {
	Writer io.Writer
}

func (c CSVPrinter) Print(lineItems []LineItem) error {
	csvWriter := csv.NewWriter(c.Writer)

	header := []string{"PCF Version"}
	boshReleaseNames := uniqueBoshReleaseNames(lineItems)
	boshPackageNames := uniqueBoshPackageNames(lineItems)

	header = append(header, boshReleaseNames...)
	header = append(header, boshPackageNames...)

	err := csvWriter.Write(header)
	if err != nil {
		return err
	}

	for _, lineItem := range lineItems {
		csvLineItem := make([]string, len(header))
		csvLineItem[0] = lineItem.PcfVersion.String()
		if found, idx := FindIndex(lineItem.BoshRelease.Name, header); found {
			csvLineItem[idx] = lineItem.BoshRelease.Version
		}

		for _, boshPackage := range lineItem.BoshRelease.Packages {
			if found, idx := FindIndex(boshPackage.Name, header); found {
				csvLineItem[idx] = boshPackage.Version
			}
		}

		err := csvWriter.Write(csvLineItem)
		if err != nil {
			return err
		}
	}
	csvWriter.Flush()
	return nil
}

func FindIndex(s string, arrayStr []string) (bool, int) {
	for idx, v := range arrayStr {
		if s == v {
			return true, idx
		}
	}
	return false, 0
}

func uniqueBoshPackageNames(items []LineItem) []string {
	boshPackageNames := map[string]interface{}{}
	uniqueBoshPackageNames := []string{}
	for _, item := range items {
		for _, boshPkg := range item.BoshRelease.Packages {
			boshPackageNames[boshPkg.Name] = nil
		}
	}
	for uniqueBoshPackage := range boshPackageNames {
		uniqueBoshPackageNames = append(uniqueBoshPackageNames, uniqueBoshPackage)
	}
	return uniqueBoshPackageNames
}

func uniqueBoshReleaseNames(items []LineItem) []string {
	boshReleaseNames := map[string]interface{}{}
	uniqueBoshReleaseNames := []string{}
	for _, item := range items {
		boshReleaseNames[item.BoshRelease.Name] = nil
	}
	for uniqueBoshRelease := range boshReleaseNames {
		uniqueBoshReleaseNames = append(uniqueBoshReleaseNames, uniqueBoshRelease)
	}
	return uniqueBoshReleaseNames
}
