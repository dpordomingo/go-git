package main

const helpLongMsg = `
NAME:
   %_COMMAND_NAME_% - Join two commits

SYNOPSIS:
  usage: %_COMMAND_NAME_% <path> <baseCommitRev> <commitRev> [-m <msg>] [--ff-only]
                           [--no-ff] [--no-commit] [--allow-unrelated-histories]
     or: %_COMMAND_NAME_% --help

 params:
    <path>           Path to the git repository
    <baseCommitRev>  Git revision of the commit that will be the base of the merge
    <commitRev>      Git revision of the commit that will be merged over the base

DESCRIPTION:
    %_COMMAND_NAME_% Incorporates changes from the passed <commitRev> over <baseCommitRev> and moves the HEAD of the repo to the merge commit.

OPTIONS:
    If no options passed, the merge commit will be avoided if it could be considered as a fast-forward, and if needed, it will be used a default merge commit message.

-m <msg>
    [NOT IMPLEMENTED]
    Set the commit message to be used for the merge commit (in case one is created).
    If not provided, an automated message will be generated.

--ff-only
    Refuse to merge and exit with 1 unless the current HEAD is already up to date or the merge can be resolved as a fast-forward.

--no-ff
    [NOT IMPLEMENTED]
    Create a merge commit even when the merge resolves as a fast-forward.

--no-commit
    [NOT IMPLEMENTED]
    With --no-commit perform the merge but pretend the merge failed and do not autocommit, to give the user a chance to
    inspect and further tweak the merge result before committing.

--allow-unrelated-histories
    [NOT IMPLEMENTED]
    By default, git merge command refuses to merge histories that do not share a common ancestor. This option can be
    used to override this safety when merging histories of two projects that started their lives independently.
    
DISCUSSION:
    Assume the following history exists and the current branch is "master":

             o---o---o---o---C <topic
            /       /
    ---C---B---o---A---o---G <master <HEAD

    Then "%_COMMAND_NAME_% <path> master topic" will replay the changes made on the topic branch since it diverged from master (i.e., A) until its current commit (i.e., C) on top of master, and record the result in a new commit (i.e., M) along with the names of the two parent commits and a log message from the user describing the changes.

             o---o---o---o---C <topic
            /       /         \
    ---C---B---o---A---o---G---M <HEAD

    Once the merge has been performed successfully HEAD will be updated to the new merge commit.
`

