package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"time"

	"gopkg.in/yaml.v3"
)

type RevertType struct {
	Name   string `yaml:"name"`
	SHA256 string `yaml:"sha256"`
}

type FixType struct {
	Id      string `yaml:"id"`
	Summary string `yaml:"summary"`
	Ref     string `yaml:"ref"`
}

type FixesType struct {
	BUGFIX   []FixType `yaml:"BUGFIX"`
	CVE      []FixType `yaml:"CVE"`
	SECURITY []FixType `yaml:"SECURITY"`
}

type ReferencesType struct {
	Type string `yaml:"type"`
	Link string `yaml:"link"`
}

type RequiredActionsType struct {
	SystemReboot   string   `yaml:"systemReboot"`
	ProductRestart string   `yaml:"productRestart"`
	ServiceRestart []string `yaml:"serviceRestart"`
}

type HotfixType struct {
	Name               string              `yaml:"name"`
	SHA256             string              `yaml:"sha256"`
	Revert             RevertType          `yaml:"revert"`
	TicketId           string              `yaml:"ticketId"`
	Released           string              `yaml:"released"`
	CompatibleReleases []string            `yaml:"compatibleReleases"`
	Type               string              `yaml:"type"`
	CompatibeNode      string              `yaml:"compatibleNode"`
	Summary            string              `yaml:"summary"`
	ImpactedArea       []string            `yaml:"impactedArea"`
	Fixes              FixesType           `yaml:"fixes"`
	Severity           string              `yaml:"severity"`
	References         []ReferencesType    `yaml:"references"`
	RequiredActions    RequiredActionsType `yaml:"requiredActions"`
	Incompatible       []string            `yaml:"incompatible"`
}

type MetadataType struct {
	Version   float64 `yaml:"version"`
	Generated string  `yaml:"generated"`
}

type NiosHotfixManifest struct {
	Type     string       `yaml:"type"`
	Metadata MetadataType `yaml:"metadata"`
	Data     []HotfixType `yaml:"data"`
}

func validate_regex(Value string, RE string) error {
	match, err := regexp.MatchString(RE, Value)
	if err != nil {
		return errors.New("Mathcing Regular Expression Failure")
	}
	if !match {
		return errors.New("Value Format is not matching with its Regular Expression")
	}
	return nil
}

func validate_datetime(Value string) error {
	_, err := time.Parse(time.RFC3339, Value)
	if err != nil {
		return errors.New("Not a valid TimeStamp")
	}
	return nil
}

func validate_manifest_file(ManifestData *NiosHotfixManifest) error {
	var hotfix_map = make(map[string]bool)

	if ManifestData.Type != "NiosHotfixManifest" {
		return errors.New("Hotfix Manifest File Type should be \"NiosHotfixManifest\"")
	}
	err := validate_datetime(ManifestData.Metadata.Generated)
	if err != nil {
		return errors.New("Hotfix Manifest File Generation TimeStamp is not Valid")
	}

	for idx, Hotfix := range ManifestData.Data {
		fmt.Printf("\n====================================================== %04d ======================================================\n", idx)

		fmt.Println("name:", Hotfix.Name)
		err := validate_regex(Hotfix.Name, ".+\\.bin")
		if err != nil {
			return fmt.Errorf("Hotfix Name format is not correct")
		}
		_, exists := hotfix_map[Hotfix.Name]
		if exists && Hotfix.Name != "" {
			return errors.New("Duplicate Hotfix entry is not allowed.")
		} else {
			hotfix_map[Hotfix.Name] = true
		}

		fmt.Println("sha256:", Hotfix.SHA256)
		err = validate_regex(Hotfix.SHA256, "[a-z0-9]{64}")
		if err != nil {
			return fmt.Errorf("Hotfix SHA256 format is not correct, the length should be 64 characters consisting only of numbers and lower alphabets")
		}

		fmt.Println("revert:")
		fmt.Println("    name:", Hotfix.Revert.Name)
		err = validate_regex(Hotfix.Revert.Name, "(^$)|(.+\\.bin)")
		if err != nil {
			return fmt.Errorf("Revert Hotfix Name format is not correct")
		}
		_, exists = hotfix_map[Hotfix.Revert.Name]
		if exists && Hotfix.Revert.Name != "" {
			return errors.New("Duplicate Hotfix entry is not allowed.")
		} else {
			hotfix_map[Hotfix.Revert.Name] = true
		}

		fmt.Println("    sha256:", Hotfix.Revert.SHA256)
		err = validate_regex(Hotfix.Revert.SHA256, "(^$)|([a-z0-9]{64})")
		if err != nil {
			return fmt.Errorf("Revert Hotfix SHA256 format is not correct, the length should be 64 characters consisting only of numbers and lower alphabets")
		}

		fmt.Println("ticketId:", Hotfix.TicketId)
		err = validate_regex(Hotfix.TicketId, "NIOS-[0-9]+")
		if err != nil {
			return fmt.Errorf("Hotfix ticketId should be a NIOS ticket id in NIOS-9999... format")
		}

		fmt.Println("released:", Hotfix.Released)
		err = validate_datetime(Hotfix.Released)
		if err != nil {
			return fmt.Errorf("%s - %s", Hotfix.Released, err.Error())
		}
		err = validate_regex(Hotfix.Released, "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}Z$")
		if err != nil {
			return fmt.Errorf("Hotfix released timestamp should be in 'yyyy-MM-ddThh:mm:ss.sssZ' format.")
		}

		fmt.Println("compatibleReleases:")
		if len(Hotfix.CompatibleReleases) <= 0 {
			return errors.New("There should be atleast one compatible release")
		}
		for _, ver := range Hotfix.CompatibleReleases {
			fmt.Println("    -", ver)
			err = validate_regex(ver, "NIOS-\\d\\.\\d\\.\\d")
			if err != nil {
				return fmt.Errorf("\"%s\" is not a valid NIOS release", ver)
			}
		}

		fmt.Println("type:", Hotfix.Type)
		if Hotfix.Type != "Generic" && Hotfix.Type != "Consolidated" {
			return errors.New("Hotfix type possible values(case sensitive) are [\"Generic\", \"Consolidated\"]")
		}

		fmt.Println("compatibleNode:", Hotfix.CompatibeNode)
		if Hotfix.CompatibeNode != "ALL" && Hotfix.CompatibeNode != "MASTER" && Hotfix.CompatibeNode != "MEMBER" {
			return errors.New("Hotfix compatibleNode possible values(case sensitive) are [\"ALL\", \"MASTER\", \"MEMBER\"]")
		}

		fmt.Println("summary:", Hotfix.Summary)
		err = validate_regex(Hotfix.Summary, ".+")
		if err != nil {
			return errors.New("Hotfix summary cannot be empty")
		}

		fmt.Println("impactedArea:")
		if len(Hotfix.ImpactedArea) <= 0 {
			return errors.New("There should be atleast one impacted area")
		}
		for _, ia := range Hotfix.ImpactedArea {
			fmt.Println("    -", ia)
			err = validate_regex(ia, ".+")
			if err != nil {
				return fmt.Errorf("%s - is not a valid impacted area", ia)
			}
		}

		fmt.Println("fixes:")
		fixes := reflect.ValueOf(Hotfix.Fixes)
		typeOfS := fixes.Type()
		for i := 0; i < fixes.NumField(); i++ {
			fix_type := typeOfS.Field(i).Name
			fmt.Printf("    %s:\n", fix_type)
			fix := fixes.Field(i).Interface().([]FixType)
			for _, f := range fix {
				fmt.Println("      - id:", f.Id)
				fmt.Println("        summary:", f.Summary)
				fmt.Println("        ref:", f.Ref)
				re := ".*"
				err_msg := ""
				switch fix_type {
				case "BUGFIX":
					re = "NIOS-[0-9]+"
					err_msg = fmt.Sprintf("%s id should be a NIOS ticket id in NIOS-9999... format", fix_type)
					break
				case "CVE":
					re = "CVE-[0-9]{4}-[0-9]+"
					err_msg = fmt.Sprintf("%s id should be a CVE ticket id in CVE-YYYY-9999... format", fix_type)
					break
				case "SECURITY":
					re = ".*" // TODO: Check the format for Security ticket
					//err_msg = fmt.Sprintf("%s id should be a NIOS ticket id in NIOS-9999... format", fix_type)
					break
				}
				err = validate_regex(f.Id, re)
				if err != nil {
					return errors.New(err_msg)
				}
				err = validate_regex(f.Summary, ".+")
				if err != nil {
					return fmt.Errorf("%s summary should not be empty", fix_type)
				}
				_, err = url.ParseRequestURI(f.Ref)
				if err != nil {
					return fmt.Errorf("%s reference \"%s\" is not a valid URL", fix_type, f.Ref)
				}
			}
		}

		fmt.Println("severity:", Hotfix.Severity)
		switch Hotfix.Severity {
		case "MANDATORY":
		case "IMPORTANT":
		case "RECOMMENDED":
		case "OPTIONAL":
			break
		default:
			return errors.New("Hotfix Severity possible values(case sensitive) are [\"MANDATORY\", \"IMPORTANT\", \"RECOMMENDED\", \"OPTIONAL\"]")
		}

		fmt.Println("references:")
		for _, r := range Hotfix.References {
			fmt.Println("  - type:", r.Type)
			fmt.Println("    link:", r.Link)
			err = validate_regex(r.Type, ".+")
			if err != nil {
				return fmt.Errorf("Hotfix references type should not be empty")
			}
			_, err = url.ParseRequestURI(r.Link)
			if err != nil {
				return fmt.Errorf("Hotfix references link \"%s\" is not a valid URL", r.Link)
			}
		}

		fmt.Println("requiredAction:")
		action := 0
		fmt.Println("    systemReboot:", Hotfix.RequiredActions.SystemReboot)
		switch Hotfix.RequiredActions.SystemReboot {
		case "Yes":
			action += 1
		case "No":
			break
		default:
			return fmt.Errorf("Hotfix requiredAction.systemReboot possible values(case sensitive) are [\"Yes\", \"No\"]")
		}
		fmt.Println("    productRestart:", Hotfix.RequiredActions.ProductRestart)
		switch Hotfix.RequiredActions.ProductRestart {
		case "Yes":
			action += 1
		case "No":
			break
		default:
			return fmt.Errorf("Hotfix requiredAction.productRestart possible values(case sensitive) are [\"Yes\", \"No\"]")
		}
		fmt.Println("    serviceRestart:")
		for _, sr := range Hotfix.RequiredActions.ServiceRestart {
			fmt.Println("      -", sr)
		}
		if len(Hotfix.RequiredActions.ServiceRestart) > 0 {
			action += 1
		}
		if action > 1 {
			return errors.New("Hotfix requiredAction should have maximum one action enabled out of systemReboot, productRestart and serviceRestart(multiple service restart is acceptable)")
		}

		fmt.Println("incompatible:")
		for _, i := range Hotfix.Incompatible {
			fmt.Println("  -", i)
			err := validate_regex(i, ".+\\.bin")
			if err != nil {
				return fmt.Errorf("Hotfix incompatible (i.e. Hotfixes which is not compatible with this hotfix) format is not correct")
			}
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:")
		fmt.Println(os.Args[0], "<cloud_manifest_file_path>")
		os.Exit(1)
	}

	HotfixManifestFile := os.Args[1]

	HFile, err := ioutil.ReadFile(HotfixManifestFile)
	if err != nil {
		fmt.Println("Error: Opening the manifest file", HotfixManifestFile)
		fmt.Println(err)
		os.Exit(1)
	}

	var ManifestData NiosHotfixManifest

	err = yaml.Unmarshal(HFile, &ManifestData)
	if err != nil {
		fmt.Println("Error: Converting Yaml data into objects")
		fmt.Println(err)
		os.Exit(1)
	}

	err = validate_manifest_file(&ManifestData)
	fmt.Printf("==================================================================================================================\n")
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	} else {
		fmt.Println("Everything Looks GOOD!!!")
	}
	fmt.Printf("==================================================================================================================\n")
}
