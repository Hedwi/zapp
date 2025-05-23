package security

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Identity struct {
	ID            int
	Fingerprint   string
	Description   string
	Type          string
	DeveloperName string
	DeveloperID   string
}

func (i Identity) SecureString() string {
	// 개발자 이름 마스킹
	nameParts := strings.Fields(i.DeveloperName)
	maskedName := make([]string, len(nameParts))
	for j, part := range nameParts {
		maskedName[j] = strings.Repeat("*", len(part))
	}
	securedName := strings.Join(maskedName, " ")

	// 개발자 ID 마스킹 (마지막 5자리만 표시)
	idLength := len(i.DeveloperID)
	securedID := fmt.Sprintf("%s%s", i.DeveloperID[:5], strings.Repeat("*", idLength-5))

	return fmt.Sprintf("%s: %s (%s)", i.Type, securedName, securedID)
}
func (i Identity) String() string {
	return fmt.Sprintf("%s: %s (%s)", i.Type, i.DeveloperName, i.DeveloperID)
}

func FindIdentity(ctx context.Context, keychain string) ([]Identity, error) {
	var cmd *exec.Cmd
	if keychain == "" {
		cmd = exec.CommandContext(ctx, "security", "find-identity", "-v")
	} else {
		cmd = exec.CommandContext(ctx, "security", "find-identity", "-v", "-k", keychain)
	}
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}
	return parseFindIdentityOutput(string(output))
}

var (
	// lineRegexp Regular expressions for parsing certificate information
	lineRegexp = regexp.MustCompile(`^\s*(\d+)\) ([A-F0-9]+) "([^"]+)"$`)
	// descRegexp Regular expression to separate description string into Type, DeveloperName, DeveloperID
	descRegexp = regexp.MustCompile(`^(.*?):\s(.*?)\s\((.*?)\)$`)
)

func parseFindIdentityOutput(output string) (identities []Identity, err error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		matches := lineRegexp.FindStringSubmatch(line)
		if matches != nil {
			id, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, err
			}
			desc := matches[3]
			descMatches := descRegexp.FindStringSubmatch(desc)
			if descMatches != nil {
				identities = append(identities, Identity{
					ID:            id,
					Fingerprint:   matches[2],
					Description:   desc,
					Type:          descMatches[1],
					DeveloperName: descMatches[2],
					DeveloperID:   descMatches[3],
				})
			} else {
				identities = append(identities, Identity{
					ID: id, Fingerprint: matches[2], Description: desc,
				})
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return identities, nil
}
