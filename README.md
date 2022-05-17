An introductory course to UNIX.


There are a few PDFs (preview at https://unix101-nicoche.koyeb.app):
- A tutorial with basics of UNIX: `tutorial-1`
- A tutorial for git CLI: `tutorial-git`
- 3 projects of varying difficulty:
  - `c-environment`, easy
  - `shell-scripting`, medium
  - `myps`, difficult
You may also find `subject-1`. It is an early version of shell-scripting.

To re-generate a PDF, say tutorial-1: `pdflatex tutorial-1.tex`

There is an automated tool to correct the git exercise of `tutorial-git` in `tutorial-git-moulinette/`. The code is pretty much garbage, but it works.

There is also an automated tool for `subject-1`. It is interactive, in the sense that the containing directory can be given as a whole to students; they can run the tests themselves and iteratively improve their solutions.
