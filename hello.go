// package main

// import (
// 	"fmt"
// 	"time"

// 	"github.com/tj/go-spin"
// )

// func main() {
// 	s := spin.New()
// 	for i := 0; i < 30; i++ {
// 		fmt.Printf("\r  \033[36mcomputing\033[m %s ", s.Next())
// 		time.Sleep(100 * time.Millisecond)
// 	}
// 	fmt.Printf("hello")
// }

package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/reechar-goog/go-prompt"
	spin "github.com/tj/go-spin"
	"google.golang.org/api/iterator"
)

func stringToSuggest(input string) prompt.Suggest {
	var result prompt.Suggest
	// result := new(prompt.Suggest)
	result.Text = "gs://" + input
	return result
}

func mymap(vs []string, f func(string) prompt.Suggest) []prompt.Suggest {
	vsm := make([]prompt.Suggest, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func spinForever(c chan string) {
	s := spin.New()
	working := true
	var proj string
	for working {
		select {
		case msg := <-c:
			proj = msg
			fmt.Print("\r                                           \r")
			if msg == "done" {
				working = false
			}
		default:
			if proj == "" {
				fmt.Printf("\r  %s \033[36mListing All Projects\033[m", s.Next())
			} else {
				fmt.Printf("\r  %s Scanning: \033[36m"+proj+"\033[m", s.Next())
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
	fmt.Println("DONE")
}

func gen(suggests []prompt.Suggest) func(prompt.Document) []prompt.Suggest {
	var result = func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterContains(suggests, d.GetWordBeforeCursor(), true)
	}

	return result
}

func dummyExecutor(in string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("Getting: %v", in)
	// q := storage.Query{Delimiter: "/"
	q := storage.Query{}
	bucket := client.Bucket(strings.Replace(in, "gs://", "", 1))
	objects := bucket.Objects(ctx, &q)
	for {
		obj, err := objects.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
			log.Fatalf("%v", err)
		}
		fmt.Println(obj.Name)
		// results = append(results, bucketAttrs.Name)
	}
}
func main() {
	// fmt.Println("Please select table.")
	// t := prompt.Input("> ", completer)
	// fmt.Println("\033[2J")
	// fmt.Println("You selected " + t)

	// projects := getProjectIDs()
	// for i, project := range projects {
	// 	fmt.Printf("%v: %v\n", i, project)
	// 	getBuckets(project)
	// }

	// buckets := getAllBuckets()

	buckets := []string{
		"p-reech-api-tf10-api-police",
		"p-reechar-terraform",
		"artifacts.reechar-kubernetes.appspot.com",
		"mrmavtest",
		"reechar-dumb-transfer",
		"databunker-reechar-pubsub",
		"dataflow-staging-us-central1-967354505112",
		"eu.artifacts.thermal-cathode-170521.appspot.com",
		"euro-bucket-reechar-london",
		"recserve_thermal-cathode-170521",
		"reechar-billing-test-export",
		"reechar-cap-one-test",
		"reechar-data-transfer-test",
		"reechar-stage-dbpub",
		"reechar-test-1",
		"reechar-test-auth",
		"reechar-utility",
		"staging.thermal-cathode-170521.appspot.com",
		"thermal-cathode-170521",
		"thermal-cathode-170521-daisy-bkt",
		"thermal-cathode-170521-data-115553",
		"thermal-cathode-170521-terraform-master",
		"thermal-cathode-170521.appspot.com",
		"us.artifacts.thermal-cathode-170521.appspot.com",
	}
	suggests := mymap(buckets, stringToSuggest)
	pt := prompt.New(dummyExecutor, gen(suggests), prompt.OptionShowCompletionAtStart(), prompt.OptionPrefix("Select a bucket >"))

	dummyExecutor(pt.Input())

	// prompt := promptui.Select{
	// 	Label: "Select Day",
	// 	Items: []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
	// 		"Saturday", "Sunday"},
	// }

	// _, result, err := prompt.Run()

	// if err != nil {
	// 	fmt.Printf("Prompt failed %v\n", err)
	// 	return
	// }

	// fmt.Printf("You choose %q\n", result)

	// items := []string{"Vim", "Emacs", "Sublime", "VSCode", "Atom"}
	// index := -1
	// var result string
	// var err error

	// for index < 0 {
	// 	prompt := promptui.SelectWithAdd{
	// 		Label:    "What's your text editor",
	// 		Items:    items,
	// 		AddLabel: "Add your own",
	// 	}

	// 	index, result, err = prompt.Run()
	// 	prompt.

	// 	if index == -1 {
	// 		items = append(items, result)
	// 	}
	// }

	// if err != nil {
	// 	fmt.Printf("Prompt failed %v\n", err)
	// 	return
	// }

	// fmt.Printf("You choose %s\n", result)

	// fmt.Println("You selected " + t)
	// for _, bucket := range buckets {
	// 	fmt.Println(bucket)
	// }

}

func getAllBuckets() []string {
	var results []string
	var c = make(chan string)
	go spinForever(c)
	projects := getProjectIDs()
	for _, project := range projects {
		// fmt.Printf("%v: %v\n", i, project)
		// bucketsFromProj := getBuckets(project)
		c <- project
		results = append(results, getBuckets(project)...)
	}
	c <- "done"
	return results
}

func getProjectIDs() []string {
	var results []string
	cmd := exec.Command("gcloud", "projects", "list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("could not get gcloud")
	}
	outputStr := string(out)
	outputArray := strings.Split(outputStr, "\n")
	for i, line := range outputArray {
		if i == 0 {
			continue
		}
		if line == "" {
			continue
		}
		// fmt.Println()
		results = append(results, strings.Fields(line)[0])
	}
	return results
}

func getBuckets(projectID string) []string {
	var results []string
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}
	buckets := client.Buckets(ctx, projectID)
	// storage.Query.
	// client.Bucket("sd").Objects
	// buckets.
	for {
		bucketAttrs, err := buckets.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}
		// fmt.Println(bucketAttrs.Name)
		results = append(results, bucketAttrs.Name)
	}
	return results
}
