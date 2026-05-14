package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

const (
	xmlRetrieve = `<?xml version="1.0" encoding="UTF-8"?>
<Package xmlns="http://soap.sforce.com/2006/04/metadata">
    <types><members>*</members><name>ExternalClientApplication</name></types>
    <types><members>*</members><name>ExtlClntAppOauthSettings</name></types>
    <types><members>*</members><name>ExtlClntAppGlobalOauthSettings</name></types>
    <types><members>*</members><name>ExtlClntAppOauthConfigurablePolicies</name></types>
    <version>58.0</version>
</Package>`

	xmlDeploy = `<?xml version="1.0" encoding="UTF-8"?>
<Package xmlns="http://soap.sforce.com/2006/04/metadata">
    <types><members>*</members><name>ExternalClientApplication</name></types>
    <types><members>*</members><name>ExtlClntAppOauthSettings</name></types>
    <types><members>*</members><name>ExtlClntAppOauthConfigurablePolicies</name></types>
    <version>65.0</version>
</Package>`
)

func main() {
	srcOrg := flag.String("s", "", "Source org alias")
	targetOrg := flag.String("t", "", "Target org alias")
	orgAlias := flag.String("o", "", "Org alias for password")

	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		printUsage()
		return
	}

	switch args[0] {
	case "create":
		if args[1] == "external" {
			handleCreateExternal(*srcOrg, *targetOrg)
		} else {
			printUsage()
		}
	case "password":
		if args[1] == "generate" {
			handlePasswordGenerate(*orgAlias)
		} else {
			printUsage()
		}
	default:
		printUsage()
	}
}

func handleCreateExternal(src, target string) {
	if src == "" || target == "" {
		fmt.Println("Error: -s and -t flags are required")
		os.Exit(1)
	}

	fileName := "package.xml"

	fmt.Printf(">>> Retrieving from %s...\n", src)
	if err := os.WriteFile(fileName, []byte(xmlRetrieve), 0644); err != nil {
		die("Write error", err)
	}
	runCmd("sf", "project", "retrieve", "start", "--manifest", fileName, "--target-org", src)

	fmt.Printf("\n>>> Deploying to %s...\n", target)
	if err := os.WriteFile(fileName, []byte(xmlDeploy), 0644); err != nil {
		die("Update error", err)
	}
	runCmd("sf", "project", "deploy", "start", "--manifest", fileName, "--target-org", target)
}

func handlePasswordGenerate(alias string) {
	if alias == "" {
		fmt.Println("Error: -o flag is required")
		os.Exit(1)
	}

	cmd := exec.Command("sf", "org", "generate", "password", "--target-org", alias)
	output, err := cmd.CombinedOutput()
	strOutput := string(output)
	fmt.Print(strOutput)

	if err != nil {
		os.Exit(1)
	}

	re := regexp.MustCompile(`password "(.*)" for user (.*)\.`)
	matches := re.FindStringSubmatch(strOutput)

	if len(matches) > 2 {
		result := fmt.Sprintf("Login: %s\nPassword: %s", matches[2], matches[1])
		if err := copyToClipboard(result); err == nil {
			fmt.Println("\n✅ Credentials copied to clipboard!")
		}
	}
}

func copyToClipboard(text string) error {
	cmd := exec.Command("pbcopy")
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}
	in.Close()
	return cmd.Wait()
}

func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		die("Command failed", err)
	}
}

func die(msg string, err error) {
	fmt.Printf("FATAL: %s: %v\n", msg, err)
	os.Exit(1)
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  sfh -s <src> -t <target> create external")
	fmt.Println("  sfh -o <alias> password generate")
}
