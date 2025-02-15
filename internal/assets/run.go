package assets

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
)

//Download and compile WikipediaMarkupCleaner into a jar
//go:generate rm -fr wikiclean
//go:generate git clone https://github.com/negapedia/wikiclean.git
//go:generate docker run --rm --user 1000:1000 -v $PWD/wikiclean:/usr/src/wikiclean -w /usr/src/wikiclean maven mvn clean compile assembly:single
//go:generate mv wikiclean/target/wikiclean-1.2-SNAPSHOT-jar-with-dependencies.jar ../textnormalizer/WikipediaMarkupCleaner.jar
//go:generate rm -fr wikiclean

//Download badwords
//go:generate rm -fr badwords
//go:generate git clone https://github.com/negapedia/badwords.git
//go:generate mv badwords/data/ ../badwords/data/
//go:generate rm -fr badwords

//Regenerate bindata
//go:generate go-bindata -pkg $GOPACKAGE -nocompress -prefix \.\./ ../badwords/... ../destemmer/... ../textnormalizer/... ../topwordspageextractor/...

//Delete everything created or downloaded
//#go:generate rm -fr ../badwords/data ../textnormalizer/WikipediaMarkupCleaner.jar

//Run executes the asset program on the given directory with the given context and data
func Run(ctx context.Context, program, workdir string, args map[string]string) (err error) {
	/*var tmpDir string
	if tmpDir, err = ioutil.TempDir(workdir, "."+program); err != nil {
		return errors.Wrapf(err, "Unable to create a temporary directory in %v", workdir)
	}
	defer os.RemoveAll(tmpDir)*/

	/*if err = RestoreAssets(tmpDir, program); err != nil {
		return errors.Wrapf(err, "Unable to restore asset %s", program)
	}*/

	commandArgs := []string{"-C", "/go/src/github.com/negapedia/wikitfidf/internal/"+program,"runandselfclean"}
	for key, value := range args {
		commandArgs = append(commandArgs, fmt.Sprintf("%v=%v", key, value))
	}

	/*if program == "textnormalizer"{
		commandArgs[1] +=program
	} else if program == "destemmer" {
		call = "make -C /go/src/github.com/negapedia/wikitfidf/internal/destemmer"
	} else if program == "topwordspageextractor" {
		call = "make -C /go/src/github.com/negapedia/wikitfidf/internal/textnormalizer"
	}*/

	cmd := exec.CommandContext(ctx, "make", commandArgs...)
	fmt.Print("To run: ")
	fmt.Println(cmd.String())
	//cmd := exec.CommandContext(ctx, "ls > /data/ls.txt & ls / > /data/ls_.txt ")

	var cmdStderr bytes.Buffer
	cmd.Stderr = &cmdStderr
	/*cmd.Dir, err = filepath.Abs(filepath.Join(tmpDir, program))
	if err != nil {
		return errors.Wrapf(err, "Unable to convert to absolute path %s", tmpDir)
	}*/

	if err = cmd.Run(); err != nil {
		return errors.Wrap(err, "Call to external command failed, with the following error stream:\n"+cmdStderr.String())
	}

	return
}
