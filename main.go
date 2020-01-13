package main

import (
	"fmt"
	_ "net/url"
	_ "os"

	"github.com/jreisinger/waf-tester/target"
	"github.com/jreisinger/waf-tester/util"
)

var paths = []string{
	"etc/passwd",
	"?page=/etc/passwd",
	"?exec=/bin/bash",
	"?id=1' or '1' = '1'",
	"?<script>",
}

func main() {
	//yamls := util.ParseYamlFiles("yaml")

	//for _, yaml := range yamls {
	//	fmt.Printf("%s\n", yaml.Tests[0].Stages[0].Stage.Input.Method)
	//	fmt.Printf("%s\n", yaml.Tests[0].Stages[0].Stage.Input.URI)
	//}

	//os.Exit(0)

	ch := make(chan target.Target)
	hosts := util.Flag()

	for _, host := range hosts {
		for _, path := range paths {
			//path = url.PathEscape(path)
			go target.Test(ch, *util.Scheme, host, path)
		}
	}

	format := "%s  (%03.0f): %s\n"

	for range hosts {
		for range paths {
			t := <-ch
			if t.Err != nil {
				fmt.Printf(format, "ERR", float64(t.StatusCode), t.Err)
			} else if t.StatusCode != 403 {
				fmt.Printf(format, "FAIL", float64(t.StatusCode), t.URL)
			} else {
				fmt.Printf(format, "OK", float64(t.StatusCode), t.URL)
			}
		}
	}
}
