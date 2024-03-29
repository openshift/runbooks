package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/grafana/regexp"
	"k8s.io/klog/v2"
)

var (
	runbookRegex = regexp.MustCompile(`.*\.md`)

	editedByRegex           = regexp.MustCompile(`<!-- Edited by .*-->`)
	downstreamCommentsRegex = regexp.MustCompile(`<!--DS: (.*)-->`)
	multipleNewLinesRegex   = regexp.MustCompile(`\n\n+`)
	asciiDocLinksRegex      = regexp.MustCompile(`link:(https://[^[\]]+)\[([^[\]]+)\]`)
)

type runbook struct {
	name string

	lastLocalUpdate     time.Time
	upstreamLastUpdated time.Time
}

func listRunbooksThatNeedUpdate(downstreamRepo *git.Repository, upstreamRepo *git.Repository) ([]runbook, []runbook) {
	localRunbooks, err := findRunbooksLastUpdateDates(downstreamRepo, downstreamRunbooksDir)
	if err != nil {
		klog.Fatal(fmt.Errorf("failed to find runbooks last update dates: %w", err))
	}

	upstreamRunbooks, err := findRunbooksLastUpdateDates(upstreamRepo, upstreamRunbooksDir)
	if err != nil {
		klog.Fatal(fmt.Errorf("failed to find runbooks last update dates: %w", err))
	}

	return checkWhichRunbooksNeedUpdate(localRunbooks, upstreamRunbooks)
}

func findRunbooksLastUpdateDates(repo *git.Repository, dir string) (map[string]time.Time, error) {
	runbooksTree, err := getDirCurrentTree(repo, dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get current runbooks tree: %w", err)
	}

	runbooks := make(map[string]time.Time)
	for _, entry := range runbooksTree.Entries {
		if !runbookRegex.MatchString(entry.Name) {
			continue
		}

		lastCommitDate, err := findLastCommitDate(repo, dir, entry.Name, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to find last commit date: %w", err)
		}
		runbooks[entry.Name] = lastCommitDate
	}

	return runbooks, nil
}

func checkWhichRunbooksNeedUpdate(localRunbooks, upstreamRunbooks map[string]time.Time) ([]runbook, []runbook) {
	var runbooksToUpdate []runbook
	var runbooksToDelete []runbook

	for name, lastUpstreamUpdate := range upstreamRunbooks {
		lastLocalUpdate, ok := localRunbooks[name]

		if !ok {
			// minimal time to be considered as last update
			lastLocalUpdate = time.UnixMilli(0)
		}

		if lastLocalUpdate.Before(lastUpstreamUpdate) {
			runbooksToUpdate = append(runbooksToUpdate, runbook{
				name:                name,
				lastLocalUpdate:     lastLocalUpdate,
				upstreamLastUpdated: lastUpstreamUpdate,
			})
		}
	}

	for name, lastLocalUpdate := range localRunbooks {
		_, ok := upstreamRunbooks[name]
		if !ok {
			runbooksToDelete = append(runbooksToDelete, runbook{
				name:                name,
				lastLocalUpdate:     lastLocalUpdate,
				upstreamLastUpdated: time.Time{},
			})
		}
	}

	return runbooksToUpdate, runbooksToDelete
}

func copyRunbook(name string) error {
	from := path.Join(upstreamCloneDir, upstreamRunbooksDir, name)
	to := path.Join(downstreamCloneDir, downstreamRunbooksDir, name)

	file, err := os.ReadFile(from)
	if err != nil {
		return fmt.Errorf("failed to read runbook %s: %w", name, err)
	}

	content := string(file)

	// Replace all 'kubectl' with 'oc'
	content = strings.ReplaceAll(content, "kubectl", "oc")

	// Remove <!-- Edited by <name>, <date> --> comments
	content = editedByRegex.ReplaceAllString(content, "")

	// Remove all US comments
	content = removeTextBetweenTags(content, "<!--USstart-->", "<!--USend-->")

	// Uncomment DS comment - <!--DS: <content>-->
	content = downstreamCommentsRegex.ReplaceAllString(content, "$1")

	// Replace 'KubeVirt' with 'OpenShift Virtualization' when not in code blocks
	content = replaceOutsideCodeBlocks(content, "KubeVirt", "OpenShift Virtualization")

	// Replace AsciiDoc links with Markdown links
	content = asciiDocLinksRegex.ReplaceAllString(content, "[$2]($1)")

	// Remove multiple (2+) new lines
	content = multipleNewLinesRegex.ReplaceAllString(content, "\n\n")

	return createAndWriteFile(to, content)
}

func removeTextBetweenTags(content, startTag, endTag string) string {
	var result strings.Builder
	startIndex := 0

	for {
		// Find start and end index of the tags
		startIndex = strings.Index(content, startTag)
		endIndex := strings.Index(content, endTag)

		// If both tags exist, remove the text between them
		if startIndex != -1 && endIndex != -1 {
			result.WriteString(content[:startIndex])
			content = content[endIndex+len(endTag):]
		} else {
			result.WriteString(content)
			break
		}
	}

	return result.String()
}

func replaceOutsideCodeBlocks(text string, old string, new string) string {
	codeBlockPattern := "```[^`]*```"
	inlineCodePattern := "`[^`]*`"
	reCodeBlock := regexp.MustCompile(codeBlockPattern)
	reInlineCode := regexp.MustCompile(inlineCodePattern)

	newText := reCodeBlock.ReplaceAllStringFunc(text, func(match string) string {
		return match
	})

	newText = reInlineCode.ReplaceAllStringFunc(newText, func(match string) string {
		return match
	})

	newText = strings.ReplaceAll(newText, old, new)

	return newText
}

func createAndWriteFile(path, content string) error {
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		_, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", path, err)
		}
	}

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}
