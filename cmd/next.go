package cmd

import (
	"errors"
	"fmt"
	"time"

	translator "crontalk/translator"

	"github.com/gorhill/cronexpr"
	"github.com/spf13/cobra"
)

var (
	nextCmd = &cobra.Command{
		Use:     "next",
		Short:   "Shows the next occurrence of a cron expression",
		Example: `crontalk next "* * * * *"`,
		RunE:    occur,
	}
	occurenceNumber = 1
	layout          = "2006-01-02 03:04PM"
)

func init() {
	nextCmd.Flags().IntVarP(&occurenceNumber, "occurence", "o", 1, "The number of occurence time")

}

func occur(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return errors.New("no cron expression detected")
	}

	translator.CronExprsn = args[0]

	if vErr := translator.Validate(); len(vErr) != 0 {
		for en, ev := range vErr {
			fmt.Printf("%v:\n", en)
			for i, e := range ev {
				fmt.Printf("%d.%v\n", i+1, e)
			}
		}
		return nil
	}
	exprns := cronexpr.MustParse(translator.CronExprsn).NextN(time.Now(), uint(occurenceNumber))
	for _, expr := range exprns {
		fmt.Println(expr.Format(layout))
	}

	return nil
}
