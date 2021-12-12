package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"github.com/pkg/errors"
	//"flag"
	//"github.com/peterbourgon/ff/v3"
	libgit "github.com/libgit2/git2go/v31"
)

// Check bladeadoe wdewd wd
type Check struct {
	fn    func(*libgit.Diff, ...string) (bool, error)
	name  string
	value int
}

var checkList = []Check{
	{checkEmptyFile1, "emptyFile1Added", 1},
	{checkEmptyFile2, "emptyFile2Added", 1},
	{checkEmptyFile3, "emptyFile3Added", 1},
	{checkReadmeWithName, "readmeWithNameAdded", 2},
	{checkWrongCalculusAdded, "wrongCalculusAdded", 2},
	{checkCalculusFixed, "calculusFixed", 3},
}

func getCommits(repo *libgit.Repository) ([]*libgit.Commit, error) {
	commits := []*libgit.Commit{}

	headRef, err := repo.Head()
	if err != nil {
		if libgit.IsErrorCode(err, libgit.ErrorCodeUnbornBranch) {
			return commits, nil
		}
		return commits, err
	}
	headAnnotatedCommit, err := repo.AnnotatedCommitFromRef(headRef)
	if err != nil {
		return commits, err
	}
	headCommit, err := repo.LookupCommit(headAnnotatedCommit.Id())
	if err != nil {
		return commits, err
	}

	commits = append(commits, headCommit)
	commit := getParent(headCommit)
	for commit != nil {
		commits = append(commits, commit)
		commit = getParent(commit)
	}
	return commits, nil
}

func getParent(commit *libgit.Commit) *libgit.Commit {
	if commit == nil {
		return nil
	}
	if commit.ParentCount() == 0 {
		return nil
	}
	return commit.Parent(0)
}

func discardSignedCommits(commits []*libgit.Commit) []*libgit.Commit {
	res := []*libgit.Commit{}
	discardedcount := 0
	for _, c := range commits {
		committer := c.Committer()
		if committer != nil && committer.Name == "GitHub" {
			discardedcount += 1
		} else {
			res = append(res, c)
		}
	}
	if discardedcount > 0 {
		fmt.Printf("Discarding %d commits (commited made via GitHub GUI)\n", discardedcount)
	}
	return res
}

func discardTooRecentCommits(commits []*libgit.Commit, limit time.Time) []*libgit.Commit {
	finalList := []*libgit.Commit{}
	commitsDiscarded := []*libgit.Commit{}
	for _, commit := range commits {
		signature := commit.Author()
		commitDate := signature.When
		if !commitDate.Before(limit) {
			commitsDiscarded = append(commitsDiscarded, commit)
		} else {
			finalList = append(finalList, commit)
		}
	}
	if len(commitsDiscarded) > 0 {
		fmt.Printf("Discarding %d commits (commited after deadline)\n", len(commitsDiscarded))
	}
	return finalList
}

type GitTutorialGrade struct {
	Grade int
	OutOf int
	Comment string
}

func gradeTutorialFromGithubInfo(ghusername string, ghrepo string, firstName string, lastName string) (*GitTutorialGrade, error) {
	tutorialURL := fmt.Sprintf("https://github.com/%s/%s", ghusername, ghrepo)
	return gradeTutorialFromUrl(tutorialURL, firstName, lastName)
}

func credentialsCallback(url string, username string, allowedTypes libgit.CredType) (*libgit.Cred, error) {
	//return libgit.NewCredentialUsername("git")
	//return libgit.NewCredentialSSHKeyFromAgent("nicolas")
	return libgit.NewCredentialDefault()
	//res, cred := libgit.NewCredDefault()
	//return libgit.ErrorCode(res), &cred
	//ret, cred := git.NewCredSshKey("git", "/Users/odewahn/.ssh/id_rsa.pub", "/Users/odewahn/.ssh/id_rsa", "")
	//fmt.Print("Enter your username: ")
	//var user string
	//fmt.Scanln(&user)

	//fmt.Print("Enter password: ")
	//var password string
	//fmt.Scanln(&password)

	//ret, cred := git.NewCredUserpassPlaintext(user, password)
	//return , &cred
}

func certificateCheckCallback(cert *libgit.Certificate, valid bool, hostname string) libgit.ErrorCode {
	return 0
}

func gradeTutorialFromUrl(tutorialURL string, firstName string, lastName string) (*GitTutorialGrade, error) {
	if tutorialURL == "" {
		return &GitTutorialGrade{
			Grade: 0,
			OutOf: 0,
			Comment: "no url provided",
		}, nil
	}

	checkValidationTable := map[string]string{}
	for _, check := range checkList {
		checkValidationTable[check.name] = "failed"
	}

	cloneTo := "./tmp"
	err := os.RemoveAll(cloneTo)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(cloneTo)

	_, err = cloneRepo(tutorialURL)
	if err != nil {
		return nil, errors.Wrapf(err, "could not clone %s", tutorialURL)
	}

	repo, err := libgit.OpenRepository(cloneTo)
	if err != nil {
		return nil, err
	}
	commits, err := getCommits(repo)
	if err != nil {
		return nil, err
	}
	submissionDeadline, err := time.Parse(time.UnixDate, "Sun Dec 10 20:00:00 UTC 2021")
	if err != nil {
		return nil, err
	}
	commits = discardTooRecentCommits(commits, submissionDeadline)
	commits = discardSignedCommits(commits)
	fmt.Printf("Analyzing %d commits...\n", len(commits))

	commitCount := len(commits)
	var parent *libgit.Commit = nil
	for i := commitCount - 1; i >= 0; i-- {
		commit := commits[i]

		diff := getDiff(repo, commit, parent)
		checksValidated, err := runChecklist(diff, firstName, lastName)
		if err != nil {
			return nil, errors.Errorf("failed to run checklist: %+v", err)
		}
		for _, checkName := range checksValidated {
			checkValidationTable[checkName] = "succeeded"
		}

		parent = commit
	}

	if firstName == "" || lastName == "" {
		checkValidationTable["readmeWithNameAdded"] = "skipped"
	}

	grade, outOf, details := computeGrade(checkValidationTable)
	return &GitTutorialGrade{
		Grade: grade,
		OutOf: outOf,
		Comment: details,
	}, nil
}

func computeGrade(checkTable map[string]string) (int, int, string) {
	total, grade := 0, 0
	checksPassed, checksFailed, checksSkipped := []string{}, []string{}, []string{}
	for _, check := range checkList {
		status := checkTable[check.name]
		if status == "skipped" {
			checksSkipped = append(checksSkipped, check.name)
			continue
		}
		total += check.value
		if status == "succeeded" {
			checksPassed = append(checksPassed, check.name)
			grade += check.value
		} else {
			checksFailed = append(checksFailed, check.name)
		}
	}

	return grade, total, fmt.Sprintf("Checks passed: %+v, checks failed: %+v, checks skipped: %+v",
		checksPassed, checksFailed, checksSkipped)
}

func getDiff(repo *libgit.Repository, commit *libgit.Commit, parent *libgit.Commit) *libgit.Diff {
	diffOptions, err := libgit.DefaultDiffOptions()
	if err != nil {
		panic(err)
	}

	tree, err := commit.Tree()
	if err != nil {
		panic(err)
	}
	var parentTree *libgit.Tree = nil
	if parent != nil {
		parentTree, err = parent.Tree()
		if err != nil {
			panic(err)
		}
	}

	diff, err := repo.DiffTreeToTree(parentTree, tree, &diffOptions)
	if err != nil {
		panic(err)
	}
	return diff
}

func runChecklist(diff *libgit.Diff, firstName, lastName string) ([]string, error) {
	validatedChecks := []string{}
	for _, check := range checkList {
		var res bool
		var err error
		if check.name == "readmeWithNameAdded" {
			if firstName != "" && lastName != "" {
				res, err = check.fn(diff, firstName, lastName)
			}
		} else {
			res, err = check.fn(diff)
		}
		if err != nil {
			return validatedChecks, err
		}
		if res {
			validatedChecks = append(validatedChecks, check.name)
		}
	}
	return validatedChecks, nil
}

func checkEmptyFile1(diff *libgit.Diff, extraStr ...string) (bool, error) {
	return addsFile(diff, "empty_file_1", true), nil
}

func checkEmptyFile2(diff *libgit.Diff, extraStr ...string) (bool, error) {
	return addsFile(diff, "empty_file_2", true), nil
}

func checkEmptyFile3(diff *libgit.Diff, extraStr ...string) (bool, error) {
	return addsFile(diff, "empty_file_3", true), nil
}

func checkReadmeWithName(diff *libgit.Diff, extraStr ...string) (bool, error) {
	if !addsFile(diff, "README.md", false) {
		return false, nil
	}
	patches, err := grabPatchesThatAdd(diff, "README.md")
	if err != nil {
		return false, err
	}
	if len(extraStr) < 2 {
		return false, fmt.Errorf("Cannot perform readme check without firstname and lastname")
	}
	for _, patch := range patches {
		firstName, lastName := extraStr[0], extraStr[1]
		if Contains(patch, firstName, false) || Contains(patch, lastName, false) {
			return true, nil
		}
	}
	return false, nil
}

func checkWrongCalculusAdded(diff *libgit.Diff, extraStr ...string) (bool, error) {
	if !addsFile(diff, "modify_me.txt", false) {
		return false, nil
	}
	patches, err := grabPatchesThatAdd(diff, "modify_me.txt")
	if err != nil {
		return false, err
	}
	for _, patch := range patches {
		patch = strings.NewReplacer(" ", "").Replace(patch)
		if strings.Contains(patch, "2+2=5") {
			return true, nil
		}
	}
	return false, nil
}

func checkCalculusFixed(diff *libgit.Diff, extraStr ...string) (bool, error) {
	if !editsFile(diff, "modify_me.txt") {
		return false, nil
	}
	patches, err := grabPatchesThatEdit(diff, "modify_me.txt")
	if err != nil {
		return false, err
	}
	for _, patch := range patches {
		patch = strings.NewReplacer(" ", "").Replace(patch)
		scanner := bufio.NewScanner(strings.NewReader(patch))

		removesWrongCalculus, addsRightCalculus := false, false

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "-") && strings.Contains(line, "2+2=5") {
				removesWrongCalculus = true
			}
			if strings.HasPrefix(line, "+") && strings.Contains(line, "2+2=4") {
				addsRightCalculus = true
			}
		}
		if scanner.Err() != nil {
			continue
		}

		if (!addsRightCalculus && removesWrongCalculus) || (addsRightCalculus && !removesWrongCalculus) {
			fmt.Printf("Warn: checkCalculusFixed is partially done\n")
		}
		if addsRightCalculus && removesWrongCalculus {
			return true, nil
		}
	}

	return false, nil
}

func editsFile(diff *libgit.Diff, filepath string) bool {
	return len(findDeltaIndicesThatEdit(diff, filepath)) > 0
}

func grabPatchesThatEdit(diff *libgit.Diff, filepath string) ([]string, error) {
	deltaIndices := findDeltaIndicesThatEdit(diff, filepath)
	if len(deltaIndices) == 0 {
		return []string{}, fmt.Errorf("Attempting to get the delta of a diff "+
			"that is supposed to edit %s when that diff does not "+
			"actually edits that file (%+v)", filepath, diff)
	}
	patchesStr := []string{}
	for _, deltaIndex := range deltaIndices {
		patch, err := diff.Patch(deltaIndex)
		if err != nil {
			panic(err)
		}
		patchStr, err := patch.String()
		if err != nil {
			panic(err)
		}
		patchesStr = append(patchesStr, patchStr)
	}
	return patchesStr, nil
}

func addsFile(diff *libgit.Diff, filepath string, onlyThisOne bool) bool {
	deltaCount, err := diff.NumDeltas()
	if err != nil {
		panic(err)
	}
	if onlyThisOne && deltaCount != 1 {
		return false
	}
	return len(findDeltaIndicesThatAdd(diff, filepath)) > 0
}

func findDeltaIndicesRelatedTo(diff *libgit.Diff, filepath string) []int {
	indices := []int{}

	deltaCount, err := diff.NumDeltas()
	if err != nil {
		panic(err)
	}
	for i := 0; i < deltaCount; i++ {
		delta, err := diff.Delta(i)
		if err != nil {
			panic(err)
		}
		if delta.OldFile.Path == filepath && delta.NewFile.Path == filepath {
			indices = append(indices, i)
		}
	}
	return indices
}

func findDeltaIndicesThatAdd(diff *libgit.Diff, filepath string) []int {
	indices := findDeltaIndicesRelatedTo(diff, filepath)
	if len(indices) == 0 {
		return indices
	}
	goodIndices := []int{}

	for _, i := range indices {
		delta, err := diff.Delta(i)
		if err != nil {
			panic(err)
		}
		if delta.OldFile.Oid.IsZero() && !delta.NewFile.Oid.IsZero() {
			goodIndices = append(goodIndices, i)
		}
	}
	return goodIndices
}

func findDeltaIndicesThatEdit(diff *libgit.Diff, filepath string) []int {
	indices := findDeltaIndicesRelatedTo(diff, filepath)
	if len(indices) == 0 {
		return indices
	}

	goodIndices := []int{}

	for _, i := range indices {
		delta, err := diff.Delta(i)
		if err != nil {
			panic(err)
		}
		if !delta.OldFile.Oid.IsZero() && !delta.NewFile.Oid.IsZero() {
			goodIndices = append(goodIndices, i)
		}
	}
	return goodIndices
}

func grabPatchesThatAdd(diff *libgit.Diff, filepath string) ([]string, error) {
	deltaIndices := findDeltaIndicesThatAdd(diff, filepath)
	if len(deltaIndices) == 0 {
		return []string{}, fmt.Errorf("Attempting to get the delta of a diff "+
			"that is supposed to add %s when that diff does not "+
			"actually adds that file (%+v)", filepath, diff)
	}
	patchesStr := []string{}
	for _, deltaIndex := range deltaIndices {
		patch, err := diff.Patch(deltaIndex)
		if err != nil {
			panic(err)
		}
		patchStr, err := patch.String()
		if err != nil {
			panic(err)
		}
		patchesStr = append(patchesStr, patchStr)
	}
	return patchesStr, nil
}

// Contains checks that content contains substring. Strict allows more
// flexibility in checks as this function can be used to perform
// checks on first/last names, where we want to not be strict on
// "complex" letters (e.g. 'é', '-'...)
func Contains(content, substring string, strict bool) bool {
	if !strict {
		content, substring = strings.ToLower(content), strings.ToLower(substring)
		accentsReplacer := strings.NewReplacer("é", "e", "è", "e", "ë", "e", "ï", "i", "ô", " ")
		punctuationReplacer := strings.NewReplacer("-", " ", "_", " ", ".", " ", ",", " ", ":", " ")

		content, substring = accentsReplacer.Replace(content), accentsReplacer.Replace(substring)
		content, substring = punctuationReplacer.Replace(content), punctuationReplacer.Replace(substring)
	}

	return strings.Contains(content, substring)
}

//func main() {
//	fs := flag.NewFlagSet("osef", flag.ExitOnError)
//	var (
//		ghUsername  = fs.String("ghUsername", "", "Github usernamefor that student")
//		ghRepo  = fs.String("ghRepo", "", "Github repository path")
//		lastName = fs.String("lastName", "", "Student lastname, needed for the README check")
//		firstName = fs.String("firstName", "", "Student firstname, needed for the README check")
//	)
//	ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())
//
//	if *ghRepo == "" && *ghUsername == "" {
//		panic("No github username or repo path given")
//	}
//	if *lastName == "" {
//		fmt.Printf("No lastname given. Will not perform README check\n")
//	}
//	if *firstName == "" {
//		fmt.Printf("No firstname given. Will not perform README check\n")
//	}
//
//	xd, err := gradeTutorialFromGithubInfo(*ghUsername, *ghRepo, *firstName, *lastName)
//	fmt.Printf("Grade %d/%d, %s\n", grade, total, comment)
//}
